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


package customservice

import (
	"context"
	"encoding/json"
	"log"

	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var classFilterResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The full name of the file / the name to match the file name with",
			Required:    true,
		},
		"match": {
			Type:        schema.TypeString,
			Description: "Matcher applying to the class name (ENDS_WITH, EQUALS or STARTS_WITH). STARTS_WITH can only be used if there is at least one annotation defined. Default value is EQUALS",
			Optional:    true,
			Default:     "EQUALS",
			// ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"ENDS_WITH", "EQUALS", "STARTS_WITH"}, false)),
		},
	},
}

var fileFilterResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The full name of the file / the name to match the file name with",
			Required:    true,
		},
		"match": {
			Type:        schema.TypeString,
			Description: "Matcher applying to the file name (ENDS_WITH, EQUALS or STARTS_WITH). Default value is ENDS_WITH (if applicable)",
			Optional:    true,
			// ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"ENDS_WITH", "EQUALS", "STARTS_WITH"}, false)),
		},
	},
}

var detectionRuleResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of the detection rule",
			Computed:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Rule enabled/disabled",
			Required:    true,
		},
		"file": {
			Type:        schema.TypeList,
			Description: "The PHP file containing the class or methods to instrument. Required for PHP custom service. Not applicable to Java and .NET",
			Optional:    true,
			MaxItems:    1,
			Elem:        fileFilterResource,
		},
		"class": {
			Type:        schema.TypeList,
			Description: "The fully qualified class or interface to instrument (or a substring if matching to a string). Required for Java and .NET custom services. Not applicable to PHP",
			Optional:    true,
			MaxItems:    1,
			Elem:        classFilterResource,
		},
		"annotations": {
			Type:        schema.TypeList,
			Description: "Additional annotations filter of the rule. Only classes where all listed annotations are available in the class itself or any of its superclasses are instrumented. Not applicable to PHP",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"method": {
			Type:        schema.TypeList,
			Description: "methods to instrument",
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Description: "The ID of the method rule",
						Computed:    true,
					},
					"name": {
						Type:        schema.TypeString,
						Description: "The method to instrument",
						Required:    true,
					},
					"returns": {
						Type:        schema.TypeString,
						Description: "Fully qualified type the method returns",
						Required:    true,
					},
					"arguments": {
						Type:        schema.TypeList,
						Description: "Fully qualified types of argument the method expects",
						Optional:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	},
}

// Resource produces the resource schema for Custom Services
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the custom service, displayed in the UI",
				Required:    true,
			},
			"order": {
				Type:        schema.TypeString,
				Description: "The order string. Sorting custom services alphabetically by their order string determines their relative ordering. Typically this is managed by Dynatrace internally and will not be present in GET responses",
				Optional:    true,
			},
			"technology": {
				Type:        schema.TypeString,
				Description: "Matcher applying to the file name (ENDS_WITH, EQUALS or STARTS_WITH). Default value is ENDS_WITH (if applicable)",
				Required:    true,
				// ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"dotNet", "go", "java", "nodeJS", "php"}, false)),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Custom service enabled/disabled",
				Required:    true,
			},
			"queue_entry_point": {
				Type:        schema.TypeBool,
				Description: "The queue entry point flag. Set to `true` for custom messaging services",
				Required:    true,
			},
			"queue_entry_point_type": {
				Type:        schema.TypeString,
				Description: "The queue entry point type (IBM_MQ, JMS, KAFKA, MSMQ or RABBIT_MQ)",
				Optional:    true,
				// ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"IBM_MQ", "JMS", "KAFKA", "MSMQ", "RABBIT_MQ"}, false)),
			},
			"process_groups": {
				Type:        schema.TypeList,
				Description: "The list of process groups the custom service should belong to",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rule": {
				Type:        schema.TypeList,
				Description: "The list of rules defining the custom service",
				Optional:    true,
				Elem:        detectionRuleResource,
			},
		},
	}
}

// func validateDiagFunc(validateFunc func(interface{}, string) ([]string, []error)) schema.SchemaValidateDiagFunc {
// 	return func(i interface{}, path cty.Path) diag.Diagnostics {
// 		warnings, errs := validateFunc(i, fmt.Sprintf("%+v", path))
// 		var diags diag.Diagnostics
// 		for _, warning := range warnings {
// 			diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: warning})
// 		}
// 		for _, err := range errs {
// 			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: err.Error()})
// 		}
// 		return diags
// 	}
// }

func resourceDataToCustomService(data *schema.ResourceData) *customservices.CustomService {
	var customService customservices.CustomService
	customService.Enabled = getBool(data, "enabled")
	customService.Name = getString(data, "name")
	customService.Order = getString(data, "order")
	customService.QueueEntryPoint = getBool(data, "queue_entry_point")
	processGroups := data.Get("process_groups")
	if processGroups != nil {
		customService.ProcessGroups = []string{}
		for _, processGroup := range processGroups.([]interface{}) {
			customService.ProcessGroups = append(customService.ProcessGroups, processGroup.(string))
		}
	}
	customService.QueueEntryPointType = customservices.QueueEntryPointType(data.Get("queue_entry_point_type").(string))
	customService.Rules = extractDetectionRules(data.Get("rule"))
	return &customService
}

func marshal(v interface{}) string {
	result := ""
	if data, err := json.MarshalIndent(v, "", "  "); err != nil {
		result = err.Error()
	} else {
		result = string(data)
	}
	return result
}

// Create expects configuration data for a Custom Service within the given ResourceData
// and sends it via POST to the Dynatrace Server in order to create that resource
func Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("customservice.Create(...)")
	}

	conf := m.(*config.ProviderConfiguration)

	var err error

	customService := resourceDataToCustomService(d)

	technology := customservices.Technology(d.Get("technology").(string))

	rest.Verbose = config.HTTPVerbose

	customServices := customservices.NewService(conf.DTenvURL, conf.APIToken)
	var stub *api.EntityShortRepresentation
	if stub, err = customServices.Create(customService, technology); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(stub.ID)
	return Read(ctx, d, m)
}

// Read queries for the contents of a Custom Service with the given ID (within ResourceData)
func Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("customservice.Read(...)")
	}
	var diags diag.Diagnostics

	var err error

	rest.Verbose = config.HTTPVerbose

	technology := customservices.Technology(d.Get("technology").(string))

	conf := m.(*config.ProviderConfiguration)
	customServices := customservices.NewService(conf.DTenvURL, conf.APIToken)
	var customService *customservices.CustomService
	if customService, err = customServices.Get(d.Id(), technology, true); err != nil {
		return diag.FromErr(err)
	}

	d.Set("enabled", customService.Enabled)
	d.Set("name", customService.Name)
	d.Set("order", customService.Order)
	d.Set("queue_entry_point", customService.QueueEntryPoint)
	d.Set("queue_entry_point_type", customService.QueueEntryPointType)
	d.Set("process_groups", customService.ProcessGroups)

	rules := make([]interface{}, 0)
	for _, detectionRule := range customService.Rules {
		rule := make(map[string]interface{}, 0)
		rule["id"] = detectionRule.ID
		rule["enabled"] = detectionRule.Enabled

		rule["annotations"] = detectionRule.Annotations

		if detectionRule.FileName != "" || detectionRule.FileNameMatcher != "" {
			rule["file"] = []interface{}{
				map[string]interface{}{
					"name":  detectionRule.FileName,
					"match": string(detectionRule.FileNameMatcher),
				},
			}
		}
		if detectionRule.ClassName != "" || detectionRule.ClassNameMatcher != "" {
			rule["class"] = []interface{}{
				map[string]interface{}{
					"name":  detectionRule.ClassName,
					"match": string(detectionRule.ClassNameMatcher),
				},
			}
		}

		methodRules := make([]interface{}, 0)
		for _, method := range detectionRule.MethodRules {
			methodRule := map[string]interface{}{
				"id":        method.ID,
				"name":      method.MethodName,
				"returns":   method.ReturnType,
				"arguments": []interface{}{},
			}
			methodArguments := []interface{}{}
			for _, arg := range method.ArgumentTypes {
				methodArguments = append(methodArguments, arg)
			}
			methodRule["arguments"] = methodArguments
			methodRules = append(methodRules, methodRule)
		}
		rule["method"] = methodRules
		rules = append(rules, rule)

	}
	d.Set("rule", rules)

	return diags
}

// Update expects configuration data for a Custom Service within the given ResourceData
// and sends a PUT request with the modified contents to the Dynatrace Server
func Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("customservice.Update(...)")
	}
	conf := m.(*config.ProviderConfiguration)

	var err error

	customService := resourceDataToCustomService(d)
	customService.ID = d.Id()
	technology := customservices.Technology(d.Get("technology").(string))

	rest.Verbose = config.HTTPVerbose

	customServices := customservices.NewService(conf.DTenvURL, conf.APIToken)
	if err = customServices.Update(customService, technology); err != nil {
		return diag.FromErr(err)
	}

	return Read(ctx, d, m)
}

// Delete removes the Custom Service with the given ID (within ResourceData) from the Dynatrace Server
func Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("customservice.Delete(...)")
	}
	conf := m.(*config.ProviderConfiguration)

	technology := customservices.Technology(d.Get("technology").(string))

	rest.Verbose = config.HTTPVerbose

	customServices := customservices.NewService(conf.DTenvURL, conf.APIToken)
	if err := customServices.Delete(d.Id(), technology); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
