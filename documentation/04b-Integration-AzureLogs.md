# Dynatrace Integration Configuration: Azure Logs Terraform Example

Azure log forwarding allows you to stream Azure logs from Azure Event Hubs into Dynatrace logs via an Azure Function App instance, additional information can be found under [Azure Logs](https://www.dynatrace.com/support/help/setup-and-configuration/setup-on-cloud-platforms/microsoft-azure-services/azure-integrations/set-up-log-forwarder-azure).

The configuration below is an example of how to configure the Azure log forwarder with Terraform.

```
locals {
  deployment_name  = "deployment_name"
  target_url       = "target_url"
  target_api_token = "target_api_token"
}

# ######################################################################################
# CREATE RESOURCE GROUP
# ######################################################################################

resource "azurerm_resource_group" "resgrp" {
  name     = "example-resources"
  location = "East US"
}

# ######################################################################################
# SET UP AZURE EVENT HUBS INSTANCE
# ######################################################################################

resource "azurerm_eventhub_namespace" "evthubns" {
  name                = "acceptanceTestEventHubNamespace"
  location            = azurerm_resource_group.resgrp.location
  resource_group_name = azurerm_resource_group.resgrp.name
  sku                 = "Basic"
  capacity            = 2

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "evthub" {
  name                = "acceptanceTestEventHub"
  namespace_name      = azurerm_eventhub_namespace.evthubns.name
  resource_group_name = azurerm_resource_group.resgrp.name
  partition_count     = 2
  message_retention   = 2
}

resource "azurerm_eventhub_namespace_authorization_rule" "evthubnsauthrule" {
  name                = "example-auth-rule"
  namespace_name      = azurerm_eventhub_namespace.evthubns.name
  resource_group_name = azurerm_eventhub_namespace.evthubns.resource_group_name
  listen              = true
  send                = false
  manage              = false
}

# ######################################################################################
# STORAGE ACCOUNT
# ######################################################################################

resource "azurerm_storage_account" "storage" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.resgrp.name
  location                 = azurerm_resource_group.resgrp.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                       = "examplekeyvault"
  location                   = azurerm_resource_group.resgrp.location
  resource_group_name        = azurerm_resource_group.resgrp.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
  sku_name                   = "standard"
}

# ######################################################################################
# CONFIGURE MONITORING SETTINGS - ROUTE TO EVENT HUB INSTANCE
# ######################################################################################

resource "azurerm_monitor_diagnostic_setting" "example" {
  name                           = "example-diagnostic-setting"
  target_resource_id             = azurerm_storage_account.storacc.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.evthubnsauthrule.authorization_rule_id

  log {
    category = "AuditEvent"
    enabled  = true

    retention_policy {
      enabled = true
      days    = 30
    }
  }
}

# #################################################################################################################
# -----------------------------------------------------------------------------------------------------------------
#
#   EVERYTHING FROM HERE ON REPLICATES WHAT IS USUALLY DONE BY
#   https://github.com/dynatrace-oss/dynatrace-azure-log-forwarder/releases/latest/download/dynatrace-azure-logs.sh
#
# -----------------------------------------------------------------------------------------------------------------
# #################################################################################################################


# ######################################################################################
# DEPLOYMENT GROUP
# ######################################################################################

data "http" "template" {
  url = "https://github.com/dynatrace-oss/dynatrace-azure-log-forwarder/releases/download/release-0.1.3/dynatrace-azure-forwarder.json"
}

resource "azurerm_resource_group_template_deployment" "example" {
  name                = "example-deployment"
  resource_group_name = azurerm_resource_group.resgrp.name
  deployment_mode     = "Incremental" # Unsure here. Should this be "Complete"?
  template_content    = data.http.template.body
  parameters = jsondecode({
    "forwarderName" : "${local.deployment_name}",
    "targetUrl" : "${local.target_url}",
    "targetAPIToken" : "${local.target_api_token}",
    "eventHubConnectionString" : "${azurerm_eventhub_namespace_authorization_rule.evthubnsauthrule.primary_connection_string}",
    "eventHubName" : "${azurerm_eventhub.evthub.name}",
    "requireValidCertificate" : true,
    "selfMonitoringEnabled" : false,
    "deployActiveGateContainer" : false,
    "targetPaasToken" : "",
    "filterConfig" : "",
    "resourceTags" : "\"LogsForwarderDeployment\":\"${local.deployment_name}\""
  })
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-app-service-plan"
  location            = azurerm_resource_group.resgrp.location
  resource_group_name = azurerm_resource_group.resgrp.name
  sku {
    tier = "Standard"
    size = "S1"
  }
}

# wait some time to allow functionapp to warmup ?
resource "azurerm_web_app" "example" {
  name                = "${local.deployment_name}-function"
  location            = azurerm_resource_group.resgrp.location
  resource_group_name = azurerm_resource_group.resgrp.name
  app_service_plan_id = azurerm_app_service_plan.example.id

  site_config {
    dotnet_framework_version = "v4.0"
  }
}

resource "azurerm_storage_container" "container" {
  name                  = "webapp-container"
  storage_account_name  = azurerm_storage_account.storage.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "log_forwarder_blob" {
  name                   = "dynatrace-azure-log-forwarder.zip"
  storage_account_name   = azurerm_storage_account.storage.name
  storage_container_name = azurerm_storage_container.container.name
  type                   = "Block"
  source                 = "https://github.com/dynatrace-oss/dynatrace-azure-log-forwarder/releases/download/release-0.1.3/dynatrace-azure-log-forwarder.zip"
}

resource "azurerm_app_service" "example" {
  name                = "${local.deployment_name}-function"
  location            = azurerm_resource_group.resgrp.location
  resource_group_name = azurerm_resource_group.resgrp.name
  app_service_plan_id = azurerm_app_service_plan.example.id

  site_config {
    always_on                = true
    dotnet_framework_version = "v4.0"
  }

  storage_account {
    name           = azurerm_storage_account.storage.name
    type           = "AzureBlob"
    account_key    = azurerm_storage_account.storage.primary_access_key
    container_name = azurerm_storage_container.container.name
    path           = "site/wwwroot"
  }
}

```