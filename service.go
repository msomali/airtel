package airtel

var _ Service = (*Client)(nil)

type (
	Service interface {
		CollectionService
		DisbursementService
		Authenticator
		AccountService
		KYCService
		TransactionService
	}
)
