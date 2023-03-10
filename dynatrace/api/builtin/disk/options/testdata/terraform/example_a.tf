# ID vu9U3hXa3q0AAAABABRidWlsdGluOmRpc2sub3B0aW9ucwAESE9TVAAQRDg5NUZGMDdDQzcxMDEyRQAkNmUzNzVmNDQtZmY3YS0zZDBjLThmOGUtYmFiZWU5MWQ2ZTMwvu9U3hXa3q0
resource "dynatrace_disk_options" "#name#" {
  nfs_show_all = true
  scope        = "HOST-1234567890000000"
  exclusions {
    exclusion {
      filesystem = "ntfs"
      mountpoint = "C:\\"
      os         = "OS_TYPE_WINDOWS"
    }
  }
}
