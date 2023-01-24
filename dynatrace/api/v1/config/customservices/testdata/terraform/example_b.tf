resource "dynatrace_custom_service" "#name#" {
  name       = "#name#"
  technology = "java"
  enabled    = true
  rule {
    enabled = true
    class {
      name  = "com.example.Prefix"
      match = "EQUALS"
    }
    method {
      name      = "methodA"
      arguments = ["java.lang.String", "java.lang.String"]
      returns   = "java.lang.String"
    }
    method {
      name    = "methodB"
      arguments = [ ]
      returns = "void"
    }
    annotations = ["com.example.ExampleAnnotation"]
  }
  queue_entry_point      = true
  queue_entry_point_type = "KAFKA"
}
