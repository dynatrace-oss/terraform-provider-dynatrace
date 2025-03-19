package bindings

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "accounts:policy-bindings"

func Service(credentials *rest.Credentials) settings.CRUDService[*bindings.PolicyBinding] {
	return &service{
		serviceClient: NewPolicyService(credentials),
	}
}

type BindingServiceClient struct {
	client rest.Client
}

func (me *service) Create(ctx context.Context, v *bindings.PolicyBinding) (*api.Stub, error) {
	return me.serviceClient.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	return me.serviceClient.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.serviceClient.Delete(ctx, id)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.serviceClient.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	return me.serviceClient.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func NewPolicyService(credentials *rest.Credentials) *BindingServiceClient {
	return &BindingServiceClient{client: rest.ClusterV2Client(credentials)}
}

type service struct {
	serviceClient *BindingServiceClient
}

func (me *BindingServiceClient) SchemaID() string {
	return SchemaID
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *BindingServiceClient) Create(ctx context.Context, v *bindings.PolicyBinding) (*api.Stub, error) {
	id := joinID(v)
	var err error
	if err = me.Update(ctx, id, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id, Name: "PolicyBindings-" + id}, nil
}

func (me *BindingServiceClient) Get(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	if err = me.client.Get(ctx, fmt.Sprintf("/iam/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), 200).Finish(&v); err != nil {
		return err
	}
	if levelType == "cluster" {
		v.Cluster = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.GroupID = groupID
	for idx, policyID := range v.PolicyIDs {
		v.PolicyIDs[idx] = fmt.Sprintf("%s#-#%s#-#%s", policyID, levelType, levelID)
	}
	return nil
}

func (me *BindingServiceClient) Update(ctx context.Context, id string, bindings *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	policyIDs := []string{}
	for _, policyID := range bindings.PolicyIDs {
		uuid, policyLevelType, policyLevelID, err := policies.SplitID(policyID, levelType, levelID)
		if policyLevelID != levelID || policyLevelType != levelType {
			return fmt.Errorf("The policy %s is defined for %s = %s. It cannot be used within the scope %s = %s", uuid, policyLevelType, policyLevelID, levelType, levelID)
		}
		if err == nil {
			policyIDs = append(policyIDs, uuid)
		} else {
			policyIDs = append(policyIDs, policyID)
		}
	}
	bindings.PolicyIDs = policyIDs

	if err = me.client.Put(ctx, fmt.Sprintf("/iam/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), bindings, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *BindingServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubs api.Stubs

	clusterInfoResponse := struct {
		ClusterUUID string `json:"clusterUuid"`
	}{}
	if err = me.client.Get(ctx, "/license/consumption/hour", 200).Finish(&clusterInfoResponse); err != nil {
		return stubs, err
	}
	bindingsResponse := struct {
		PolicyBindings []struct {
			PolicyUUID string   `json:"policyUuid"`
			Groups     []string `json:"groups"`
		} `json:"policyBindings"`
	}{}
	if err = me.client.Get(ctx, fmt.Sprintf("/iam/repo/cluster/%s/bindings", clusterInfoResponse.ClusterUUID), 200).Finish(&bindingsResponse); err != nil {
		return stubs, err
	}
	bindingsMap := map[string]bool{}
	for _, bindings := range bindingsResponse.PolicyBindings {
		for _, group := range bindings.Groups {
			joinedId := fmt.Sprintf("%s#-#%s#-#%s", group, "cluster", clusterInfoResponse.ClusterUUID)
			if _, exists := bindingsMap[joinedId]; !exists {
				stubs = append(stubs, &api.Stub{ID: joinedId, Name: joinedId})
				bindingsMap[joinedId] = true
			}
		}
	}

	environmentsResponse := struct {
		Environments []struct {
			ID string `json:"id"`
		} `json:"environments"`
	}{}
	if err = me.client.Get(ctx, "/environments?pageSize=1000", 200).Finish(&environmentsResponse); err != nil {
		return stubs, err
	}
	for _, environment := range environmentsResponse.Environments {
		bindingsResponse := struct {
			PolicyBindings []struct {
				PolicyUUID string   `json:"policyUuid"`
				Groups     []string `json:"groups"`
			} `json:"policyBindings"`
		}{}
		if err = me.client.Get(ctx, fmt.Sprintf("/iam/repo/environment/%s/bindings", environment.ID), 200).Finish(&bindingsResponse); err != nil {
			return stubs, err
		}
		for _, bindings := range bindingsResponse.PolicyBindings {
			for _, group := range bindings.Groups {
				joinedId := fmt.Sprintf("%s#-#%s#-#%s", group, "environment", environment.ID)
				if _, exists := bindingsMap[joinedId]; !exists {
					stubs = append(stubs, &api.Stub{ID: joinedId, Name: joinedId})
					bindingsMap[joinedId] = true
				}
			}
		}
	}

	return stubs, nil
}

func (me *BindingServiceClient) Delete(ctx context.Context, id string) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var binding bindings.PolicyBinding
	if err = me.Get(ctx, id, &binding); err != nil {
		return err
	}
	for _, policyID := range binding.PolicyIDs {
		if err = me.client.Delete(ctx, fmt.Sprintf("/iam/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyID, groupID), 204).Finish(); err != nil {
			return err
		}
	}
	return nil
}

func splitID(id string) (groupID string, levelType string, levelID string, err error) {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%s is not a valid ID for a policy", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func joinID(binding *bindings.PolicyBinding) string {
	levelType, levelID := getLevel(binding)
	return fmt.Sprintf("%s#-#%s#-#%s", binding.GroupID, levelType, levelID)
}

func getLevel(binding *bindings.PolicyBinding) (string, string) {
	if len(binding.Cluster) > 0 {
		return "cluster", binding.Cluster
	}
	if len(binding.Environment) > 0 {
		return "environment", binding.Environment
	}
	return "global", "global"
}
