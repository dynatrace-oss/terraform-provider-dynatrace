resource "dynatrace_queue_manager" "#name#" {
  name = "#name#"
  clusters = ["Cluster 1","Cluster 2"]
  alias_queues {
    alias_queue {
      alias_queue_name = "Alias Queue A"
      base_queue_name = "Base Queue A"
    }
    alias_queue {
      alias_queue_name = "Alias Queue B"
      base_queue_name = "Base Queue B"
      cluster_visibility = ["Cluster 1"]  
    }
    alias_queue {
      alias_queue_name = "Alias Queue C"
      base_queue_name = "Base Queue C"
      cluster_visibility = ["Cluster 1", "Cluster 2"]  
    }
  }
  remote_queues {
    remote_queue {
      local_queue_name = "Local Queue A"
      remote_queue_name = "Remote Queue A"
      remote_queue_manager = "Remote Queue Manager A"
    }
    remote_queue {
      local_queue_name = "Local Queue B"
      remote_queue_name = "Remote Queue B"
      remote_queue_manager = "Remote Queue Manager B"
      cluster_visibility = ["Cluster 1"]    
    }
    remote_queue {
      local_queue_name = "Local Queue C"
      remote_queue_name = "Remote Queue C"
      remote_queue_manager = "Remote Queue Manager C"
      cluster_visibility = ["Cluster 1","Cluster 2"]    
    }
  }
  cluster_queues {
    cluster_queue {
      local_queue_name = "Local Queue A"
    }
    cluster_queue {
      local_queue_name = "Local Queue B"
      cluster_visibility = ["Cluster 1"]    
    }
    cluster_queue {
      local_queue_name = "Local Queue C"
      cluster_visibility = ["Cluster 1", "Cluster 2"]    
    }
  }
}
