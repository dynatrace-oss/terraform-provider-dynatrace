package services

import "strings"

func DefaultMetrics(metric string) []*AzureMonitoredMetric {
	for _, service := range DefaultServices {
		if strings.ToLower(service.ServiceName) == strings.ToLower(metric) {
			return service.MonitoredMetrics
		}
	}
	return nil
}

var DefaultServices = []*Settings{
	{
		ServiceName: "cloud:azure:storage:storageaccounts:table",
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
		ServiceName: "cloud:azure:storage:storageaccounts",
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
		ServiceName: "cloud:azure:storage:storageaccounts:file",
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
		ServiceName: "cloud:azure:storage:storageaccounts:blob",
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
		ServiceName: "cloud:azure:storage:storageaccounts:queue",
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
		ServiceName:      "azure_service_bus_namespace",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_vm",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_redis_cache",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_load_balancer",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_sql",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_event_hub_namespace",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_web_app",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_iot_hub",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_application_gateway",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_cosmos_db",
		MonitoredMetrics: nil,
	},
	{
		ServiceName:      "azure_api_management_service",
		MonitoredMetrics: nil,
	},
}
