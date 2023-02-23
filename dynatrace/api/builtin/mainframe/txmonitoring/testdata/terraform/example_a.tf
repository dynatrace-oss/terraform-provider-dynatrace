# ID vu9U3hXa3q0AAAABAB5idWlsdGluOm1haW5mcmFtZS50eG1vbml0b3JpbmcABnRlbmFudAAGdGVuYW50ACQwYWYxNWEwOS05YWM0LTMyZGEtOTZjZi01Y2Q3NjI1Y2MxNja-71TeFdrerQ
resource "dynatrace_mainframe_transaction_monitoring" "mainframe_transaction_monitoring" {
  group_cics_regions                             = true
  group_ims_regions                              = false
  monitor_all_ctg_protocols                      = false
  monitor_all_incoming_web_requests              = false
  node_limit                                     = 500
  zos_cics_service_detection_uses_transaction_id = false
  zos_ims_service_detection_uses_transaction_id  = false
}
