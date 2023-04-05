resource "dynatrace_opentelemetry_metrics" "#name#" {
  additional_attributes_to_dimension_enabled = true
  meter_name_to_dimension_enabled            = true
  scope                                      = "environment"
  additional_attributes {
    additional_attribute {
      enabled       = true
      attribute_key = "terraform.test.add"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "service.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "host.id"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "host.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "faas.id"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "faas.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.cluster.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.node.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.namespace.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.deployment.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.statefulset.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.daemonset.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.job.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.cronjob.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.pod.uid"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "k8s.pod.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "process.executable.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "dt.metrics.source"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "dt.entity.host"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "dt.entity.process_group_instance"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "dt.entity.host_group"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "dt.kubernetes.workload.kind"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "dt.kubernetes.workload.name"
    }
    additional_attribute {
      enabled       = true
      attribute_key = "dt.kubernetes.cluster.id"
    }
  }
  to_drop_attributes {
    to_drop_attribute {
      enabled       = true
      attribute_key = "terraform.test.drop"
    }
  }
}