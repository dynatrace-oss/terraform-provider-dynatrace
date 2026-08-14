resource "dynatrace_custom_service_order" "this" {
  dotnet = [
    dynatrace_custom_service.dotnet[0].id,
    dynatrace_custom_service.dotnet[1].id,
  ]
  java = [
    dynatrace_custom_service.java[0].id,
    dynatrace_custom_service.java[1].id,
  ]
}

resource "dynatrace_custom_service" "java" {
  count      = 2
  name       = "java-#name#-${count.index}"
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
      name      = "methodB"
      arguments = []
      returns   = "void"
    }
    annotations = ["com.example.ExampleAnnotation"]
  }
  queue_entry_point = false
}

resource "dynatrace_custom_service" "dotnet" {
  count      = 2
  name       = "dotnet-#name#-${count.index}"
  technology = "dotNet"
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
      name      = "methodB"
      arguments = []
      returns   = "void"
    }
    annotations = ["com.example.ExampleAnnotation"]
  }
  queue_entry_point = false
}
