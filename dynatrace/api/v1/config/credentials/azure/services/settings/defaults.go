package services

import "strings"

func DefaultMetrics(metric string) []*AzureMonitoredMetric {
	for _, service := range DefaultServices {
		if strings.ToLower(service.Name) == strings.ToLower(metric) {
			return service.MonitoredMetrics
		}
	}
	return nil
}

var DefaultServices = []*Settings{
	{
		Name: "cloud:azure:storage:storageaccounts:table",
		MonitoredMetrics: []*AzureMonitoredMetric{
			{
				Name:       "TableCapacity",
				Dimensions: []string{},
			},
			{
				Name:       "Transactions",
				Dimensions: []string{},
			},
		},
	},
	{
		Name: "cloud:azure:storage:storageaccounts",
		MonitoredMetrics: []*AzureMonitoredMetric{
			{
				Name:       "Transactions",
				Dimensions: []string{},
			},
			{
				Name:       "UsedCapacity",
				Dimensions: []string{},
			},
		},
	},
	{
		Name: "cloud:azure:storage:storageaccounts:file",
		MonitoredMetrics: []*AzureMonitoredMetric{
			{
				Name:       "FileCapacity",
				Dimensions: []string{},
			},
			{
				Name:       "Transactions",
				Dimensions: []string{},
			},
		},
	},
	{
		Name: "cloud:azure:storage:storageaccounts:blob",
		MonitoredMetrics: []*AzureMonitoredMetric{
			{
				Name:       "BlobCapacity",
				Dimensions: []string{},
			},
			{
				Name:       "Transactions",
				Dimensions: []string{},
			},
		},
	},
	{
		Name: "cloud:azure:storage:storageaccounts:queue",
		MonitoredMetrics: []*AzureMonitoredMetric{
			{
				Name:       "QueueCapacity",
				Dimensions: []string{},
			},
			{
				Name:       "Transactions",
				Dimensions: []string{},
			},
		},
	},
	{
		Name:             "azure_service_bus_namespace",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_vm",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_redis_cache",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_load_balancer",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_sql",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_event_hub_namespace",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_web_app",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_iot_hub",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_application_gateway",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_cosmos_db",
		MonitoredMetrics: nil,
	},
	{
		Name:             "azure_api_management_service",
		MonitoredMetrics: nil,
	},
}
