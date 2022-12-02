package hcl2json_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dtcookie/assert"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2json"
)

func TestDashboards(t *testing.T) {
	testConfig(t, "../test_data/dashboards/example_a.tf", "../test_data/dashboards/example_a.json", func(me, ma map[string]interface{}) {
		delete(me, "metadata")
		delete(me, "id")
		me["dashboardMetadata"].(map[string]interface{})["name"] = ma["dashboardMetadata"].(map[string]interface{})["name"]
		me["dashboardMetadata"].(map[string]interface{})["preset"] = ma["dashboardMetadata"].(map[string]interface{})["preset"]
	})
}

func TestAlerting(t *testing.T) {
	testConfig(t, "../test_data/alerting/example_a.tf", "../test_data/alerting/example_a.json", func(me, ma map[string]interface{}) {
		me["name"] = ma["name"]
		if mgmz, ok := me["managementZone"]; ok && len(mgmz.(string)) == 0 {
			delete(me, "managementZone")
		}
		if mgmz, ok := ma["managementZone"]; ok && len(mgmz.(string)) == 0 {
			delete(ma, "managementZone")
		}
	})
}

func TestCustomServices(t *testing.T) {
	testConfig(t, "../test_data/custom_services/example_a.tf", "../test_data/custom_services/example_a.json")
}

func TestManagementZones(t *testing.T) {
	testConfig(t, "../test_data/management_zones/example_a.tf", "../test_data/management_zones/example_a.json")
}

func TestMaintenance(t *testing.T) {
	testConfig(t, "../test_data/maintenance/example_a.tf", "../test_data/maintenance/example_a.json")
}

func TestRequestAttributes(t *testing.T) {
	testConfig(t, "../test_data/request_attributes/example_a.tf", "../test_data/request_attributes/example_a.json")
}
func TestAutoTags(t *testing.T) {
	testConfig(t, "../test_data/auto_tags/example_a.tf", "../test_data/auto_tags/example_a.json", func(me, ma map[string]interface{}) {
	})
}

func TestAWSCredentials(t *testing.T) {
	testConfig(t, "../test_data/credentials/aws/example_a.tf", "../test_data/credentials/aws/example_a.json", func(me, ma map[string]interface{}) {
		delete(ma["authenticationData"].(map[string]interface{})["keyBasedAuthentication"].(map[string]interface{}), "secretKey")
	})
}

func TestAzureCredentials(t *testing.T) {
	testConfig(t, "../test_data/credentials/azure/example_a.tf", "../test_data/credentials/azure/example_a.json", func(me, ma map[string]interface{}) {
		delete(ma, "key")
	})
}
func TestK8sCredentials(t *testing.T) {
	testConfig(t, "../test_data/credentials/k8s/example_a.tf", "../test_data/credentials/k8s/example_a.json", func(me, ma map[string]interface{}) {
		delete(me, "endpointStatus")
		delete(me, "endpointStatusInfo")
		delete(ma, "authToken")
	})
}

func TestPCFCredentials(t *testing.T) {
	testConfig(t, "../test_data/credentials/cloudfoundry/example_a.tf", "../test_data/credentials/cloudfoundry/example_a.json", func(me, ma map[string]interface{}) {
		delete(me, "endpointStatus")
		delete(me, "endpointStatusInfo")
		delete(ma, "password")
	})
}

func TestServiceAnomalies(t *testing.T) {
	testConfig(t, "../test_data/anomalies/services/example_a.tf", "../test_data/anomalies/services/example_a.json")
}

func TestApplicationAnomalies(t *testing.T) {
	testConfig(t, "../test_data/anomalies/applications/example_a.tf", "../test_data/anomalies/applications/example_a.json")
}

func TestHostAnomalies(t *testing.T) {
	testConfig(t, "../test_data/anomalies/hosts/example_a.tf", "../test_data/anomalies/hosts/example_a.json")
}

func TestDatabaseAnomalies(t *testing.T) {
	testConfig(t, "../test_data/anomalies/databases/example_a.tf", "../test_data/anomalies/databases/example_a.json")
}

func TestCustomAnomalies(t *testing.T) {
	testConfig(t, "../test_data/anomalies/metric_events/example_a.tf", "../test_data/anomalies/metric_events/example_a.json", func(me, ma map[string]interface{}) {
		delete(me, "eventType")
		delete(me, "disabledReason")
		delete(me, "warningReason")
	})
}

func TestDiskAnomalies(t *testing.T) {
	testConfig(t, "../test_data/anomalies/disk_events/example_a.tf", "../test_data/anomalies/disk_events/example_a.json")
}

func TestCalculatedServiceMetrics(t *testing.T) {
	testConfig(t, "../test_data/metrics/calculated/service/example_a.tf", "../test_data/metrics/calculated/service/example_a.json", func(me, ma map[string]interface{}) {
		ma["tsmMetricKey"] = me["tsmMetricKey"]
	})
}

func TestServiceNaming(t *testing.T) {
	testConfig(t, "../test_data/naming/services/example_a.tf", "../test_data/naming/services/example_a.json")
}

func TestHostNaming(t *testing.T) {
	testConfig(t, "../test_data/naming/hosts/example_a.tf", "../test_data/naming/hosts/example_a.json")
}

func TestProcessGroupNaming(t *testing.T) {
	testConfig(t, "../test_data/naming/processgroups/example_a.tf", "../test_data/naming/processgroups/example_a.json")
}

func TestSLO(t *testing.T) {
	t.Skip()
	testConfig(t, "../test_data/slo/example_a.tf", "../test_data/slo/example_a.json")
}

func TestSpanEntryPoints(t *testing.T) {
	testConfig(t, "../test_data/spans/entrypoints/example_a.tf", "../test_data/spans/entrypoints/example_a.json")
}

func TestSpanCaptureRules(t *testing.T) {
	testConfig(t, "../test_data/spans/capture/example_a.tf", "../test_data/spans/capture/example_a.json")
}

func TestSpanContextPropagation(t *testing.T) {
	testConfig(t, "../test_data/spans/ctxprop/example_a.tf", "../test_data/spans/ctxprop/example_a.json")
}

func TestResourceAttributes(t *testing.T) {
	testConfig(t, "../test_data/spans/resattr/example_a.tf", "../test_data/spans/resattr/example_a.json", func(me, ma map[string]interface{}) {
		if elems, ok := me["attributeKeys"]; ok {
			if elems != nil {
				for _, item := range elems.([]interface{}) {
					delete(item.(map[string]interface{}), "masking")
				}
			}
		}
	})
}

func TestSpanAttributes(t *testing.T) {
	t.Skip()
	testConfig(t, "../test_data/spans/resattr/example_a.tf", "../test_data/spans/resattr/example_a.json", func(me, ma map[string]interface{}) {
		if elems, ok := me["attributeKeys"]; ok {
			if elems != nil {
				for _, item := range elems.([]interface{}) {
					delete(item.(map[string]interface{}), "masking")
				}
			}
		}
	})
}

func testConfig(t *testing.T, tfFile string, jsonFile string, anon ...func(me, ma map[string]interface{})) {
	t.Helper()
	assert := assert.New(t)
	var err error
	var configs []*hcl2json.Record
	var data []byte
	if configs, err = hcl2json.HCL2Config(tfFile); err != nil {
		t.Error(err)
		return
	}
	for _, config := range configs {
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			t.Error(err)
			return
		}
		// t.Log(string(data))
		ma := map[string]interface{}{}
		json.Unmarshal(data, &ma)
		data, err = os.ReadFile(jsonFile)
		if err != nil {
			t.Error(err)
			return
		}

		me := map[string]interface{}{}
		json.Unmarshal(data, &me)
		if _, ok := me["schemaId"]; ok {
			me = me["value"].(map[string]interface{})
		}
		delete(me, "id")
		delete(me, "name")
		delete(me, "displayName")
		delete(me, "label")
		delete(me, "metadata")
		delete(me, "connectionStatus")

		delete(ma, "id")
		delete(ma, "name")
		delete(ma, "displayName")
		delete(ma, "label")
		delete(ma, "metadata")
		if len(anon) == 1 && anon[0] != nil {
			anon[0](me, ma)
		}
		assert.Equals(streamline(me), streamline(ma))
		// out, _ := json.Marshal(me)
		// t.Log("expected", string(out))
		// out, _ = json.Marshal(ma)
		// t.Log("actual", string(out))
	}
}
