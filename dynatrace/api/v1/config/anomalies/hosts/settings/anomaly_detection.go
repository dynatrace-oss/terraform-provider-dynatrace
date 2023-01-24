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

package hosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/connection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/cpu"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/disks"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/disks/inodes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/disks/slow"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/disks/space"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/gc"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/java"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/java/oom"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/java/oot"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/memory"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/droppedpackets"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/errors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/retransmission"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/tcp"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/utilization"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AnomalyDetection Configuration of anomaly detection for hosts.
type AnomalyDetection struct {
	NetworkDroppedPacketsDetection     *droppedpackets.DetectionConfig `json:"networkDroppedPacketsDetection"`     // Configuration of high number of dropped packets detection.
	HighNetworkDetection               *utilization.DetectionConfig    `json:"highNetworkDetection"`               // Configuration of high network utilization detection.
	NetworkHighRetransmissionDetection *retransmission.DetectionConfig `json:"networkHighRetransmissionDetection"` // Configuration of high retransmission rate detection.
	NetworkTcpProblemsDetection        *tcp.DetectionConfig            `json:"networkTcpProblemsDetection"`        // Configuration of TCP connectivity problems detection.
	NetworkErrorsDetection             *errors.DetectionConfig         `json:"networkErrorsDetection"`             // Configuration of high number of network errors detection.
	HighMemoryDetection                *memory.DetectionConfig         `json:"highMemoryDetection"`                // Configuration of high memory usage detection.
	HighCPUSaturationDetection         *cpu.DetectionConfig            `json:"highCpuSaturationDetection"`         // Configuration of high CPU saturation detection
	OutOfMemoryDetection               *oom.DetectionConfig            `json:"outOfMemoryDetection"`               // Configuration of Java out of memory problems detection.
	OutOfThreadsDetection              *oot.DetectionConfig            `json:"outOfThreadsDetection"`              // Configuration of Java out of threads problems detection.
	HighGcActivityDetection            *gc.DetectionConfig             `json:"highGcActivityDetection"`            // Configuration of high Garbage Collector activity detection.
	ConnectionLostDetection            *connection.LostDetectionConfig `json:"connectionLostDetection"`            // Configuration of lost connection detection.
	DiskSlowWritesAndReadsDetection    *slow.DetectionConfig           `json:"diskSlowWritesAndReadsDetection"`    // Configuration of slow running disks detection.
	DiskLowSpaceDetection              *space.DetectionConfig          `json:"diskLowSpaceDetection"`              // Configuration of low disk space detection.
	DiskLowInodesDetection             *inodes.DetectionConfig         `json:"diskLowInodesDetection"`             // Configuration of low disk inodes number detection.
}

func (me *AnomalyDetection) Name() string {
	return "host_anomalies"
}

func (me *AnomalyDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"memory": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of high memory usage detection",
			Elem:        &schema.Resource{Schema: new(memory.DetectionConfig).Schema()},
		},
		"cpu": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of high CPU saturation detection",
			Elem:        &schema.Resource{Schema: new(cpu.DetectionConfig).Schema()},
		},
		"gc": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of high Garbage Collector activity detection",
			Elem:        &schema.Resource{Schema: new(gc.DetectionConfig).Schema()},
		},
		"connections": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of lost connection detection",
			Elem:        &schema.Resource{Schema: new(connection.LostDetectionConfig).Schema()},
		},
		"network": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of network related anomalies",
			Elem:        &schema.Resource{Schema: new(network.DetectionConfig).Schema()},
		},
		"disks": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of disk related anomalies",
			Elem:        &schema.Resource{Schema: new(disks.DetectionConfig).Schema()},
		},
		"java": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of Java related anomalies",
			Elem:        &schema.Resource{Schema: new(java.DetectionConfig).Schema()},
		},
	}
}

func (me *AnomalyDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"memory":      me.HighMemoryDetection,
		"cpu":         me.HighCPUSaturationDetection,
		"gc":          me.HighGcActivityDetection,
		"connections": me.ConnectionLostDetection,
		"java": &java.DetectionConfig{
			OutOfMemoryDetection:  me.OutOfMemoryDetection,
			OutOfThreadsDetection: me.OutOfThreadsDetection,
		},
		"network": &network.DetectionConfig{
			NetworkDroppedPacketsDetection:     me.NetworkDroppedPacketsDetection,
			HighNetworkDetection:               me.HighNetworkDetection,
			NetworkHighRetransmissionDetection: me.NetworkHighRetransmissionDetection,
			NetworkTcpProblemsDetection:        me.NetworkTcpProblemsDetection,
			NetworkErrorsDetection:             me.NetworkErrorsDetection,
		},
		"disks": &disks.DetectionConfig{
			Speed:  me.DiskSlowWritesAndReadsDetection,
			Space:  me.DiskLowSpaceDetection,
			Inodes: me.DiskLowInodesDetection,
		},
	})
}

func (me *AnomalyDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	me.NetworkDroppedPacketsDetection = &droppedpackets.DetectionConfig{Enabled: false}
	me.HighNetworkDetection = &utilization.DetectionConfig{Enabled: false}
	me.NetworkHighRetransmissionDetection = &retransmission.DetectionConfig{Enabled: false}
	me.NetworkTcpProblemsDetection = &tcp.DetectionConfig{Enabled: false}
	me.NetworkErrorsDetection = &errors.DetectionConfig{Enabled: false}
	me.DiskSlowWritesAndReadsDetection = &slow.DetectionConfig{Enabled: false}
	me.DiskLowSpaceDetection = &space.DetectionConfig{Enabled: false}
	me.DiskLowInodesDetection = &inodes.DetectionConfig{Enabled: false}
	me.HighMemoryDetection = &memory.DetectionConfig{Enabled: false}
	me.HighCPUSaturationDetection = &cpu.DetectionConfig{Enabled: false}
	me.OutOfMemoryDetection = &oom.DetectionConfig{Enabled: false}
	me.OutOfThreadsDetection = &oot.DetectionConfig{Enabled: false}
	me.HighGcActivityDetection = &gc.DetectionConfig{Enabled: false}
	me.ConnectionLostDetection = &connection.LostDetectionConfig{Enabled: false}

	if _, ok := decoder.GetOk("network.#"); ok {
		cfg := new(network.DetectionConfig)

		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "network", 0)); err != nil {
			return err
		}
		me.NetworkDroppedPacketsDetection = cfg.NetworkDroppedPacketsDetection
		me.HighNetworkDetection = cfg.HighNetworkDetection
		me.NetworkHighRetransmissionDetection = cfg.NetworkHighRetransmissionDetection
		me.NetworkTcpProblemsDetection = cfg.NetworkTcpProblemsDetection
		me.NetworkErrorsDetection = cfg.NetworkErrorsDetection
	}
	if _, ok := decoder.GetOk("disks.#"); ok {
		cfg := new(disks.DetectionConfig)

		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "disks", 0)); err != nil {
			return err
		}
		me.DiskSlowWritesAndReadsDetection = cfg.Speed
		me.DiskLowSpaceDetection = cfg.Space
		me.DiskLowInodesDetection = cfg.Inodes
	}
	if _, ok := decoder.GetOk("memory.#"); ok {
		me.HighMemoryDetection = new(memory.DetectionConfig)
		if err := me.HighMemoryDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "memory", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("cpu.#"); ok {
		me.HighCPUSaturationDetection = new(cpu.DetectionConfig)
		if err := me.HighCPUSaturationDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "cpu", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("gc.#"); ok {
		me.HighGcActivityDetection = new(gc.DetectionConfig)
		if err := me.HighGcActivityDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "gc", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("connections.#"); ok {
		me.ConnectionLostDetection = new(connection.LostDetectionConfig)
		if err := me.ConnectionLostDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "connections", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("java.#"); ok {
		cfg := new(java.DetectionConfig)

		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "java", 0)); err != nil {
			return err
		}
		me.OutOfMemoryDetection = cfg.OutOfMemoryDetection
		me.OutOfThreadsDetection = cfg.OutOfThreadsDetection
	}
	return nil
}
