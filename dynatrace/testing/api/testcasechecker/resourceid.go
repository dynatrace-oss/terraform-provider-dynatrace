//go:build integration

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package testcasechecker

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type TestCaseCheckResourceIDs struct {
	resourceIDs map[string]string
}

// CheckCreate persists the current set IDs for every resource
func (tr *TestCaseCheckResourceIDs) CheckCreate(s *terraform.State) error {
	tr.resourceIDs = make(map[string]string)
	for name, resourceState := range s.RootModule().Resources {
		tr.resourceIDs[name] = resourceState.Primary.ID
	}
	return nil
}

// CheckUpdate checks if every resource is present in the update and that the ID matches
// Therefore, it can be used to check if resources were recreated or deleted instead of updated
func (tr *TestCaseCheckResourceIDs) CheckUpdate(s *terraform.State) error {
	errs := make([]error, 0)
	newResourceIDs := make(map[string]string)
	for name, resourceState := range s.RootModule().Resources {
		newResourceIDs[name] = resourceState.Primary.ID
	}
	for name, id := range tr.resourceIDs {
		newID, ok := newResourceIDs[name]
		if !ok {
			errs = append(errs, fmt.Errorf("could not find resource %q in updated state", name))
		} else if id != newID {
			errs = append(errs, fmt.Errorf("the resource %q was created instead of updated (old ID: %s, new ID: %s)", name, id, newID))
		}
	}
	return errors.Join(errs...)
}
