package envs

import (
	"context"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

// NewService creates a new Service Client
// baseURL should look like this: "https://siz65484.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(credentials *rest.Credentials) *ServiceClient {
	return &ServiceClient{client: rest.ClusterV2Client(credentials)}
}

func evalRetry(rerr *rest.Error, environment *Environment) bool {
	if len(rerr.ConstraintViolations) > 0 {
		for _, violation := range rerr.ConstraintViolations {
			if violation.Message == "Cannot set Synthetic monitors quota as Synthetic monitors are not allowed for this license." {
				environment.Quotas.Synthetic = nil
				return true
			} else if violation.Message == "Cannot set DEM units quota as DEM units are not allowed for this license." {
				environment.Quotas.DEMUnits = nil
				return true
			} else if violation.Message == "Cannot set Log Monitoring retention as Log Monitoring is not configured for this cluster." {
				environment.Storage.LogMonitoringRetention = nil
				environment.Storage.LogMonitoringStorage = nil
				return true
			} else if violation.Message == "Cannot set Log Monitoring quota as Log Monitoring is not allowed for this license." {
				environment.Quotas.LogMonitoring = nil
				return true
			}
		}
	}
	return false
}

// Create TODO: documentation
func (cs *ServiceClient) Create(ctx context.Context, environment *Environment) (*api.Stub, error) {
	var err error

	if len(opt.String(environment.ID)) > 0 {
		return nil, errors.New("you MUST NOT provide an ID within the Dashboard payload upon creation")
	}
	var stub api.Stub
	retry := true
	for retry {
		if err = cs.client.Post(ctx, "/environments", environment, 201).Finish(&stub); err != nil {
			switch rerr := err.(type) {
			case *rest.Error:
				retry = evalRetry(rerr, environment)
			case rest.Error:
				retry = evalRetry(&rerr, environment)
			default:
				return nil, err
			}
		} else {
			retry = false
		}
	}
	return &stub, err
}

// Update TODO: documentation
func (cs *ServiceClient) Update(ctx context.Context, environment *Environment) error {
	var err error
	retry := true

	for retry {
		if err = cs.client.Put(ctx, fmt.Sprintf("/environments/%s", opt.String(environment.ID)), environment, 204).Finish(); err != nil {
			switch rerr := err.(type) {
			case *rest.Error:
				retry = evalRetry(rerr, environment)
			case rest.Error:
				retry = evalRetry(&rerr, environment)
			default:
				return err
			}
		} else {
			retry = false
		}
	}
	return err
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("empty ID provided for the Dashboard to delete")
	}
	env, err := cs.Get(ctx, id)
	if err != nil {
		return err
	}
	if env.State == States.Enabled {
		env.State = States.Disabled
		retry := true
		for retry {
			if err = cs.Update(ctx, env); err != nil {
				switch rerr := err.(type) {
				case *rest.Error:
					retry = evalRetry(rerr, env)
				case rest.Error:
					retry = evalRetry(&rerr, env)
				default:
					return err
				}
			} else {
				retry = false
			}
		}
	}
	return cs.client.Delete(ctx, fmt.Sprintf("/environments/%s", id), 204).Finish()
}

// Get TODO: documentation
func (cs *ServiceClient) Get(ctx context.Context, id string) (*Environment, error) {
	if len(id) == 0 {
		return nil, errors.New("empty ID provided for the environment to fetch")
	}

	var err error

	var environment Environment
	if err = cs.client.Get(ctx, fmt.Sprintf("/environments/%s?includeConsumptionInfo=true&includeStorageInfo=true", id), 200).Finish(&environment); err != nil {
		return nil, err
	}
	return &environment, nil
}

// ListAll TODO: documentation
func (cs *ServiceClient) ListAll(ctx context.Context) (*EnvironmentList, error) {
	var err error

	var environmentList EnvironmentList
	if err = cs.client.Get(ctx, "/environments", 200).Finish(&environmentList); err != nil {
		return nil, err
	}
	return &environmentList, nil
}

type EnvironmentList struct {
	Environments []*Environment `json:"environments"`
}
