package airtel

type (
	Service interface {
		Authenticator
		AccountService
		TransactionService
		CollectionService
		KYCService
		DisbursementService
	}
)
