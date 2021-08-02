package billingModel

import (
	"net/url"
	"strconv"
)

// Customer encapsulates details about a Customer registered in Stripe.
//
// see https://stripe.com/docs/api#customer_object
type Customer struct {
	ID         string `json:"id"`
	Desc       string `json:"description,omitempty"`
	Email      string `json:"email,omitempty"`
	Created    int64  `json:"created"`
	Balance    int64  `json:"account_balance"`
	Delinquent bool   `json:"delinquent"`
}

// StripePayment handles communication with Stripe
type StripePayment struct{}

// CreateCustomer Create new customer to Paypal and return teh user key.
func (r *StripePayment) CreateCustomer(params CustomerParams) (string, error) {
	// https://stripe.com/docs/api/customers/create
	customer := Customer{}

	values := url.Values{}
	appendCustomerParamsToValues(&params, &values)
	// simulate a request to striè
	// err := query("POST", "/v1/customers", values, &customer)
	customer.ID = "cus_J9Od1Dxs2IGSBU"
	return customer.ID, nil
}

// appendCustomerParamsToValues Helper functions
func appendCustomerParamsToValues(c *CustomerParams, values *url.Values) {
	// add optional parameters, if specified
	if c.Email != "" {
		values.Add("email", c.Email)
	}
	if c.Desc != "" {
		values.Add("description", c.Desc)
	}
	if c.Coupon != "" {
		values.Add("coupon", c.Coupon)
	}
	if c.Plan != "" {
		values.Add("plan", c.Plan)
	}
	if c.TrialEnd != 0 {
		values.Add("trial_end", strconv.FormatInt(c.TrialEnd, 10))
	}
	if c.AccountBalance != 0 {
		values.Add("account_balance", strconv.FormatInt(c.AccountBalance, 10))
	}
	if c.Quantity != 0 {
		values.Add("quantity", strconv.FormatInt(c.Quantity, 10))
	}

	// add metadata, if specified
	for k, v := range c.Metadata {
		values.Add("metadata["+k+"]", v)
	}

	// add optional credit card details, if specified
	if c.Card != nil {
		appendCardParamsToValues(c.Card, values)
	} else if len(c.Token) != 0 {
		values.Add("card", c.Token)
	}
}

func appendCardParamsToValues(c *CardParams, values *url.Values) {
	values.Add("card[number]", c.Number)
	values.Add("card[exp_month]", strconv.Itoa(c.ExpMonth))
	values.Add("card[exp_year]", strconv.Itoa(c.ExpYear))
	if c.Name != "" {
		values.Add("card[name]", c.Name)
	}
	if c.CVC != "" {
		values.Add("card[cvc]", c.CVC)
	}
	if c.Address1 != "" {
		values.Add("card[address_line1]", c.Address1)
	}
	if c.Address2 != "" {
		values.Add("card[address_line2]", c.Address2)
	}
	if c.AddressZip != "" {
		values.Add("card[address_zip]", c.AddressZip)
	}
	if c.AddressState != "" {
		values.Add("card[address_state]", c.AddressState)
	}
	if c.AddressCountry != "" {
		values.Add("card[address_country]", c.AddressCountry)
	}
}
