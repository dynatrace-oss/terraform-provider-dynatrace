resource "dynatrace_transaction_start_filters" "transaction_start_filters" {
  cics_terminal_transaction_ids = [ "DTAX", "ATAX" ]
  cics_transaction_ids          = [ "TIPU" ]
  ims_transaction_ids           = [ "FAKE" ]
  ims_terminal_transaction_ids = [ "DTAX", "ATAX" ]
}
