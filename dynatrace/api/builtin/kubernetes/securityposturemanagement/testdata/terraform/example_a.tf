resource "dynatrace_kubernetes_spm" "#name#" {
    scope = "KUBERNETES_CLUSTER-1234567890000000"
    configuration_dataset_pipeline_enabled = true
}