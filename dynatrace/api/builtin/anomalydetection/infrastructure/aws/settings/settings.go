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

package aws

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Ec2CandidateHighCpuDetection     *Ec2CandidateHighCpuDetectionConfig     `json:"ec2CandidateHighCpuDetection"`
	ElbHighConnectionErrorsDetection *ElbHighConnectionErrorsDetectionConfig `json:"elbHighConnectionErrorsDetection"`
	LambdaHighErrorRateDetection     *LambdaHighErrorRateDetectionConfig     `json:"lambdaHighErrorRateDetection"`
	RdsHighCpuDetection              *RdsHighCpuDetectionConfig              `json:"rdsHighCpuDetection"`
	RdsHighMemoryDetection           *RdsHighMemoryDetectionConfig           `json:"rdsHighMemoryDetection"`
	RdsHighWriteReadLatencyDetection *RdsHighWriteReadLatencyDetectionConfig `json:"rdsHighWriteReadLatencyDetection"`
	RdsLowStorageDetection           *RdsLowStorageDetectionConfig           `json:"rdsLowStorageDetection"`
	RdsRestartsSequenceDetection     *RdsRestartsSequenceDetectionConfig     `json:"rdsRestartsSequenceDetection"`
}

func (me *Settings) Name() string {
	return "aws_anomalies"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ec_2_candidate_high_cpu_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Ec2CandidateHighCpuDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"elb_high_connection_errors_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ElbHighConnectionErrorsDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"lambda_high_error_rate_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(LambdaHighErrorRateDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rds_high_cpu_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(RdsHighCpuDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rds_high_memory_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(RdsHighMemoryDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rds_high_write_read_latency_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(RdsHighWriteReadLatencyDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rds_low_storage_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(RdsLowStorageDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rds_restarts_sequence_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(RdsRestartsSequenceDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"ec_2_candidate_high_cpu_detection":     me.Ec2CandidateHighCpuDetection,
		"elb_high_connection_errors_detection":  me.ElbHighConnectionErrorsDetection,
		"lambda_high_error_rate_detection":      me.LambdaHighErrorRateDetection,
		"rds_high_cpu_detection":                me.RdsHighCpuDetection,
		"rds_high_memory_detection":             me.RdsHighMemoryDetection,
		"rds_high_write_read_latency_detection": me.RdsHighWriteReadLatencyDetection,
		"rds_low_storage_detection":             me.RdsLowStorageDetection,
		"rds_restarts_sequence_detection":       me.RdsRestartsSequenceDetection,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"ec_2_candidate_high_cpu_detection":     &me.Ec2CandidateHighCpuDetection,
		"elb_high_connection_errors_detection":  &me.ElbHighConnectionErrorsDetection,
		"lambda_high_error_rate_detection":      &me.LambdaHighErrorRateDetection,
		"rds_high_cpu_detection":                &me.RdsHighCpuDetection,
		"rds_high_memory_detection":             &me.RdsHighMemoryDetection,
		"rds_high_write_read_latency_detection": &me.RdsHighWriteReadLatencyDetection,
		"rds_low_storage_detection":             &me.RdsLowStorageDetection,
		"rds_restarts_sequence_detection":       &me.RdsRestartsSequenceDetection,
	})
}
