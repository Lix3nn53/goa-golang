package billingModel

// PaypalPayment handles communication with PayPal
type PaypalPayment struct {
}

// CreateCustomer Create new customer to Paypal and return teh user key.
func (r *PaypalPayment) CreateCustomer(params CustomerParams) (string, error) {
	return "paypal_J9Od1Dxs2IGSBU", nil // fake user key
}
