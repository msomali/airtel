package airtel

var (
	_ AccountService      = (*Client)(nil)
	_ CollectionService   = (*Client)(nil)
	_ Authenticator       = (*Client)(nil)
	_ DisbursementService = (*Client)(nil)
	_ TransactionService  = (*Client)(nil)
	_ KYCService          = (*Client)(nil)
)
