resource "dynatrace_queue_sharing_groups" "#name#" {
    name = "#name#"
    queue_managers = ["QueueManager1","QueueManager2","QueueManager3"]
    shared_queues = ["SharedQueue1","SharedQueue2"]
}