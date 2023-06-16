resource "dynatrace_ibm_mq_filters" "#name#" {
    cics_mq_queue_id_includes = ["StringA1","StringB1","StringC1"]
    cics_mq_queue_id_excludes = ["StringA2","StringB2","StringC2"]
    ims_mq_queue_id_excludes = ["StringA4","StringB4","StringC4"]
    ims_cr_trn_id_includes = ["StringA5","StringB5","StringC5"]
}