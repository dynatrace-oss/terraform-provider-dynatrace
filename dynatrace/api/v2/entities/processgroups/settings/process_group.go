package processgroups

import entity "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"

type ProcessGroup struct {
	DisplayName string `json:"displayName"`
	Properties  struct {
		AWSNameTag      string `json:"awsNameTag"`
		AzureHostName   string `json:"azureHostName"`
		AzureSiteName   string `json:"azureSiteName"`
		BoshName        string `json:"boshName"`
		ConditionalName string `json:"conditionalName"`
		CustomizedName  string `json:"customizedName"`
		DetectedName    string `json:"detectedName"`
		GCPZone         string `json:"gcpZone"`
		ListenPorts     []int  `json:"listenPorts"`
		ManagementZones []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"managementZones"`
		Metadata []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"metadata"`
		SoftwareTechnologies []struct {
			Type    string `json:"type"`
			Edition string `json:"edition"`
			Version string `json:"version"`
		} `json:"softwareTechnologies"`
		OneAgentCustomHostName string      `json:"oneAgentCustomHostName"`
		Tags                   entity.Tags `json:"tags"`
		// CustomPgMetadata map[string]any `json:"customPgMetadata"`
	} `json:"properties"`
}

/*


{
  "type": "PROCESS_GROUP",
  "displayName": "Process Group",
  "dimensionKey": "dt.entity.process_group",
  "entityLimitExceeded": false,
  "tags": "List",
  "managementZones": "List",
  "fromRelationships": [
    {
      "id": "isNetworkClientOfProcessGroup",
      "toTypes": [
        "AWS_APPLICATION_LOAD_BALANCER",
        "AWS_LAMBDA_FUNCTION",
        "AWS_NETWORK_LOAD_BALANCER",
        "AZURE_API_MANAGEMENT_SERVICE",
        "AZURE_APPLICATION_GATEWAY",
        "AZURE_COSMOS_DB",
        "AZURE_EVENT_HUB_NAMESPACE",
        "AZURE_FUNCTION_APP",
        "AZURE_IOT_HUB",
        "AZURE_LOAD_BALANCER",
        "AZURE_REDIS_CACHE",
        "AZURE_SERVICE_BUS_NAMESPACE",
        "AZURE_SQL_SERVER",
        "AZURE_STORAGE_ACCOUNT",
        "CUSTOM_DEVICE_GROUP",
        "ELASTIC_LOAD_BALANCER",
        "PROCESS_GROUP",
        "RELATIONAL_DATABASE_SERVICE"
      ]
    },
    {
      "id": "isPgOfCa",
      "toTypes": [
        "CLOUD_APPLICATION"
      ]
    },
    {
      "id": "isPgOfCai",
      "toTypes": [
        "CLOUD_APPLICATION_INSTANCE"
      ]
    },
    {
      "id": "isPgOfCg",
      "toTypes": [
        "CONTAINER_GROUP"
      ]
    },
    {
      "id": "runsOn",
      "toTypes": [
        "HOST"
      ]
    },
    {
      "id": "runsOnResource",
      "toTypes": [
        "CUSTOM_DEVICE"
      ]
    }
  ],
  "toRelationships": [
    {
      "id": "isDockerContainerOfPg",
      "fromTypes": [
        "DOCKER_CONTAINER_GROUP_INSTANCE"
      ]
    },
    {
      "id": "isHostGroupOf",
      "fromTypes": [
        "HOST_GROUP"
      ]
    },
    {
      "id": "isInstanceOf",
      "fromTypes": [
        "PROCESS_GROUP_INSTANCE"
      ]
    },
    {
      "id": "isNamespaceOfPg",
      "fromTypes": [
        "CLOUD_APPLICATION_NAMESPACE"
      ]
    },
    {
      "id": "isNetworkClientOfProcessGroup",
      "fromTypes": [
        "AWS_APPLICATION_LOAD_BALANCER",
        "AWS_LAMBDA_FUNCTION",
        "AWS_NETWORK_LOAD_BALANCER",
        "AZURE_API_MANAGEMENT_SERVICE",
        "AZURE_APPLICATION_GATEWAY",
        "AZURE_COSMOS_DB",
        "AZURE_EVENT_HUB_NAMESPACE",
        "AZURE_FUNCTION_APP",
        "AZURE_IOT_HUB",
        "AZURE_LOAD_BALANCER",
        "AZURE_REDIS_CACHE",
        "AZURE_SERVICE_BUS_NAMESPACE",
        "AZURE_SQL_SERVER",
        "AZURE_STORAGE_ACCOUNT",
        "CUSTOM_DEVICE_GROUP",
        "ELASTIC_LOAD_BALANCER",
        "PROCESS_GROUP",
        "RELATIONAL_DATABASE_SERVICE"
      ]
    },
    {
      "id": "isPgAppOf",
      "fromTypes": [
        "AZURE_FUNCTION_APP",
        "AZURE_WEB_APP"
      ]
    },
    {
      "id": "runsOn",
      "fromTypes": [
        "SERVICE"
      ]
    }
  ]
}
*/
