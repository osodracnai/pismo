package domain

const (
	OperationTypePurchase            int64 = 1
	OperationTypeInstallmentPurchase int64 = 2
	OperationTypeWithdrawal          int64 = 3
	OperationTypePayment             int64 = 4
)

func IsValidOperationType(operationTypeID int64) bool {
	switch operationTypeID {
	case OperationTypePurchase, OperationTypeInstallmentPurchase, OperationTypeWithdrawal, OperationTypePayment:
		return true
	default:
		return false
	}
}

func NormalizeAmount(operationTypeID int64, amount float64) float64 {
	signedAmount := amount
	if signedAmount < 0 {
		signedAmount = -signedAmount
	}

	switch operationTypeID {
	case OperationTypePayment:
		return signedAmount
	default:
		return -signedAmount
	}
}
