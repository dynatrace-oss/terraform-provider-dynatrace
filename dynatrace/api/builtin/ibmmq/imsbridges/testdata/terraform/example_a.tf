resource "dynatrace_ims_bridges" "#name#" {
    name = "#name#"
    queue_managers {
        queue_manager {
            name = "QueueManagerExample1"
        }
        queue_manager {
            name = "QueueManagerExample2"
            queue_manager_queue = ["Queue1","Queue2"]
        }
    }
}