resource "dynatrace_custom_service" "#name#" {
        name = "#name#"
        technology = "java"
        enabled = true
        rule {
                enabled = true
                class {
                    name = "com.example.Prefix"
                    match = "EQUALS"
                }
                method {
                        name = "methodA"
                        arguments = [ "java.lang.String", "java.lang.String" ]
                        returns = "java.lang.String"
                }
                method {
                        name = "methodB"
                        returns = "void"
                }
                annotations = [ "com.example.ExampleAnnotation" ]
        }
        rule {
                enabled = true
                class {
                    name = "com.example.Suffix"
                    match = "ENDS_WITH"
                }
                method {
                        name = "methodC"
                        arguments = [ "java.lang.String", "java.lang.String" ]
                        returns = "java.lang.String"
                }
                method {
                        name = "methodD"
                        returns = "void"
                }
        }
        queue_entry_point = false
}