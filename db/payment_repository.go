package db

// CreatePaymentHistory creates a new payment history record
func (r *Repository) CreatePaymentHistory(payment *PaymentHistory) error {
	result := r.db.Create(payment)
	return result.Error
}
