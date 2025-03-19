/**
* @license
* Copyright 2020 Dynatrace LLC
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

package publicendpoints

import (
	"context"

	publicendpoints "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/publicendpoints/settings"
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
	return &ServiceClient{client: rest.ClusterV1Client(credentials)}
}

type AddressSettings struct {
	Address *string `json:"address,omitempty"`
}

type AdditionalAddressesSettings struct {
	AdditionalAddresses []string `json:"additionalAddresses,omitempty"`
}

// Create TODO: documentation
func (cs *ServiceClient) Create(ctx context.Context, config *publicendpoints.Settings) error {
	return cs.Update(ctx, config)
}

// Update TODO: documentation
func (cs *ServiceClient) Update(ctx context.Context, config *publicendpoints.Settings) error {
	if config.WebUiAddress != nil {
		webUiAddress := AddressSettings{config.WebUiAddress}
		if err := cs.client.Post(ctx, "/endpoint/webUiAddress", webUiAddress, 200).Finish(); err != nil {
			return err
		}
	}
	if len(config.AdditionalWebUiAddresses) > 0 {
		additionalWebUiAddresses := AdditionalAddressesSettings{config.AdditionalWebUiAddresses}
		if err := cs.client.Post(ctx, "/endpoint/additionalWebUiAddresses", additionalWebUiAddresses, 200).Finish(); err != nil {
			return err
		}
	}
	if config.BeaconForwarderAddress != nil {
		beaconForwarderAddress := AddressSettings{config.BeaconForwarderAddress}
		if err := cs.client.Post(ctx, "/endpoint/beaconForwarderAddress", beaconForwarderAddress, 200).Finish(); err != nil {
			return err
		}
	}
	if config.CDNAddress != nil {
		cdnAddress := AddressSettings{config.CDNAddress}
		if err := cs.client.Post(ctx, "/endpoint/cdnAddress", cdnAddress, 200).Finish(); err != nil {
			return err
		}
	}

	return nil
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete() error {
	return nil
}

// Get TODO: documentation
func (cs *ServiceClient) Get(ctx context.Context) (*publicendpoints.Settings, error) {
	var err error
	webUiAddress := AddressSettings{}
	additionalWebUiAddresses := AdditionalAddressesSettings{}
	beaconForwarderAddress := AddressSettings{}
	cdnAddress := AddressSettings{}

	if err = cs.client.Get(ctx, "/endpoint/webUiAddress", 200).Finish(&webUiAddress); err != nil {
		return nil, err
	}
	if err = cs.client.Get(ctx, "/endpoint/additionalWebUiAddresses", 200).Finish(&additionalWebUiAddresses); err != nil {
		return nil, err
	}
	if err = cs.client.Get(ctx, "/endpoint/beaconForwarderAddress", 200).Finish(&beaconForwarderAddress); err != nil {
		return nil, err
	}
	if err = cs.client.Get(ctx, "/endpoint/cdnAddress", 200).Finish(&cdnAddress); err != nil {
		return nil, err
	}

	config := publicendpoints.Settings{
		WebUiAddress:             webUiAddress.Address,
		AdditionalWebUiAddresses: additionalWebUiAddresses.AdditionalAddresses,
		BeaconForwarderAddress:   beaconForwarderAddress.Address,
		CDNAddress:               cdnAddress.Address,
	}

	return &config, nil
}

// Get TODO: documentation
func (cs *ServiceClient) GetWebUiAddress(ctx context.Context) (*string, error) {
	webUiAddress := AddressSettings{}
	if err := cs.client.Get(ctx, "/endpoint/webUiAddress", 200).Finish(&webUiAddress); err != nil {
		return nil, err
	}

	return webUiAddress.Address, nil
}

// Get TODO: documentation
func (cs *ServiceClient) GetAdditionalWebUiAddresses(ctx context.Context) ([]string, error) {
	additionalWebUiAddresses := AdditionalAddressesSettings{}
	if err := cs.client.Get(ctx, "/endpoint/additionalWebUiAddresses", 200).Finish(&additionalWebUiAddresses); err != nil {
		return nil, err
	}

	return additionalWebUiAddresses.AdditionalAddresses, nil
}

// Get TODO: documentation
func (cs *ServiceClient) GetBeaconForwarderAddress(ctx context.Context) (*string, error) {
	beaconForwarderAddress := AddressSettings{}
	if err := cs.client.Get(ctx, "/endpoint/beaconForwarderAddress", 200).Finish(&beaconForwarderAddress); err != nil {
		return nil, err
	}

	return beaconForwarderAddress.Address, nil
}

// Get TODO: documentation
func (cs *ServiceClient) GetCDNAddress(ctx context.Context) (*string, error) {
	cdnAddress := AddressSettings{}
	if err := cs.client.Get(ctx, "/endpoint/cdnAddress", 200).Finish(&cdnAddress); err != nil {
		return nil, err
	}

	return cdnAddress.Address, nil
}
