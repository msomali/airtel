package airtel

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
