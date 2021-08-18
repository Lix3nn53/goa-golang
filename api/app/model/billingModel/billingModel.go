package billingModel

import (
	appError "goa-golang/app/error"
)

// Identify Service types
type Identify string

const (
	// AccountPaypal Payment types
	AccountStripe Identify = "stripe"
	// AccountPaypal Payment types
	AccountPaypal Identify = "paypal"
)

// CardParams encapsulates options for Creating or Updating Credit Cards.
type CardParams struct {
	Name           string
	Number         string
	ExpMonth       int
	ExpYear        int
	CVC            string
	Address1       string
	Address2       string
	AddressCountry string
	AddressState   string
	AddressZip     string
}

// CustomerParams encapsulates options for creating and updating Customers.
type CustomerParams struct {
	Email          string
	Desc           string
	Card           *CardParams
	Token          string
	Coupon         string
	Plan           string
	TrialEnd       int64
	AccountBalance int64
	Metadata       map[string]string
	Quantity       int64
}

// CreateCustomer define the request params to add a new payment method
type CreateCustomer struct {
	Identify       Identify `json:"identify"`
	CustomerParams `json:"customerParams"`
}

// Payment define the payment interface struct
type Payment struct {
	Identify Identify `json:"identify"`
	CustomerParams
	PaymentMethod PaymentInterface
}

//PaymentInterface define the payment interface methods
type PaymentInterface interface {
	CreateCustomer(params CustomerParams) (string, error)
}

// StripeAdapter adapts Stripe API
type StripeAdapter struct {
	Payment *StripePayment
}

// CreateCustomer Call to the payment provider to generate a new cutomer
func (b *StripeAdapter) CreateCustomer(params CustomerParams) (string, error) {
	return b.Payment.CreateCustomer(params)
}

// PayPalAdapter adapts PayPal API
type PayPalAdapter struct {
	Payment *PaypalPayment
}

// CreateCustomer Init a new customer from Paypal
func (p *PayPalAdapter) CreateCustomer(params CustomerParams) (string, error) {
	return p.Payment.CreateCustomer(params)
}

// GetPaymentAdapter get Adapter
func GetPaymentAdapter(identify Identify) (PaymentInterface, error) {
	switch identify {
	case AccountPaypal:
		payPalAdapter := PayPalAdapter{
			Payment: &PaypalPayment{},
		}
		return &payPalAdapter, nil
	case AccountStripe:
		stripeAdapter := StripeAdapter{
			Payment: &StripePayment{},
		}
		return &stripeAdapter, nil
	}
	return nil, appError.ErrInvalidPaymentMethod
}
