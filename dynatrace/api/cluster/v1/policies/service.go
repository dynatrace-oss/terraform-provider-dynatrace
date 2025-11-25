/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package policies

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	policies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/policies/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type PolicyServiceClient struct {
	client rest.Client
}

func NewPolicyService(credentials *rest.Credentials) *PolicyServiceClient {
	return &PolicyServiceClient{client: rest.ClusterV2Client(credentials)}
}

func Service(credentials *rest.Credentials) settings.CRUDService[*policies.Policy] {
	return &service{
		serviceClient: NewPolicyService(credentials),
	}
}

type service struct {
	serviceClient *PolicyServiceClient
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.serviceClient.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *policies.Policy) error {
	return me.serviceClient.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return "accounts:policies"
}

func (me *service) Create(ctx context.Context, v *policies.Policy) (*api.Stub, error) {
	return me.serviceClient.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *policies.Policy) error {
	return me.serviceClient.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.serviceClient.Delete(ctx, id)
}

func (me *PolicyServiceClient) SchemaID() string {
	return "accounts:policies"
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *PolicyServiceClient) Create(ctx context.Context, v *policies.Policy) (*api.Stub, error) {
	var err error
	levelType, levelID := getLevel(v)

	var pcr PolicyCreateResponse
	if err = me.client.Post(ctx, fmt.Sprintf("/iam/repo/%s/%s/policies", levelType, levelID), v, 201).Finish(&pcr); err != nil {
		return nil, err
	}
	return &api.Stub{ID: joinID(pcr.UUID, v), Name: v.Name}, nil
}

func (me *PolicyServiceClient) Get(ctx context.Context, id string, v *policies.Policy) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}

	if err = me.client.Get(ctx, fmt.Sprintf("/iam/repo/%s/%s/policies/%s", levelType, levelID, uuid), 200).Finish(&v); err != nil {
		return err
	}
	if levelType == "cluster" {
		v.Cluster = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	return nil
}

func (me *PolicyServiceClient) Update(ctx context.Context, id string, user *policies.Policy) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}

	if err = me.client.Put(ctx, fmt.Sprintf("/iam/repo/%s/%s/policies/%s", levelType, levelID, uuid), user, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *PolicyServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubs api.Stubs

	clusterInfoResponse := struct {
		ClusterUUID string `json:"clusterUuid"`
	}{}

	if err = me.client.Get(ctx, "/license/consumption/hour", 200).Finish(&clusterInfoResponse); err != nil {
		return stubs, err
	}
	policiesResponse := struct {
		Policies []struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"policies"`
	}{}
	if err = me.client.Get(ctx, fmt.Sprintf("/iam/repo/cluster/%s/policies", clusterInfoResponse.ClusterUUID), 200).Finish(&policiesResponse); err != nil {
		return stubs, err
	}
	for _, policy := range policiesResponse.Policies {
		stubs = append(stubs, &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "cluster", clusterInfoResponse.ClusterUUID), Name: policy.Name})
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
		policiesResponse := struct {
			Policies []struct {
				UUID string `json:"uuid"`
				Name string `json:"name"`
			} `json:"policies"`
		}{}
		if err = me.client.Get(ctx, fmt.Sprintf("/iam/repo/environment/%s/policies", environment.ID), 200).Finish(&policiesResponse); err != nil {
			return stubs, err
		}
		for _, policy := range policiesResponse.Policies {
			stubs = append(stubs, &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "environment", environment.ID), Name: policy.Name})
		}
	}

	return stubs, nil
}

func (me *PolicyServiceClient) Delete(ctx context.Context, id string) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}
	return me.client.Delete(ctx, fmt.Sprintf("/iam/repo/%s/%s/policies/%s", levelType, levelID, uuid), 204).Finish()
}

func SplitID(id string) (uuid string, levelType string, levelID string, err error) {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%s is not a valid ID for a policy", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func joinID(uuid string, policy *policies.Policy) string {
	levelType, levelID := getLevel(policy)
	return fmt.Sprintf("%s#-#%s#-#%s", uuid, levelType, levelID)
}

func getLevel(policy *policies.Policy) (string, string) {
	if len(policy.Cluster) > 0 {
		return "cluster", policy.Cluster
	}
	if len(policy.Environment) > 0 {
		return "environment", policy.Environment
	}
	return "global", "global"
}
