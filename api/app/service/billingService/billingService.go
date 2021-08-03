package billingService

import (
	appError "goa-golang/app/error"
	"goa-golang/app/model/billingModel"
	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/billingRepository"
)

//BillingServiceInterface define the user service interface methods
type BillingServiceInterface interface {
	AddBilling(user userModel.User, payment billingModel.Payment) error
	GetPaymentAdapter(customer billingModel.CreateCustomer) (*billingModel.Payment, error)
}

// BillingService handles communication with the user repository
type BillingService struct {
	paymentRepo billingRepository.BillingRepositoryInterface
}

// NewUserService implements the user service interface.
func NewBillingService(paymentRepo billingRepository.BillingRepositoryInterface) BillingServiceInterface {
	return &BillingService{
		paymentRepo,
	}
}

// FindByID implements the method to store a new a user model
func (s *BillingService) AddBilling(user userModel.User, payment billingModel.Payment) error {

	key, err := payment.PaymentMethod.CreateCustomer(payment.CustomerParams)
	if err != nil {
		return err
	}

	return s.paymentRepo.CreateBillingService(payment.Identify, key, user.ID)
}

// FindByID implements the method to store a new a user model
func (s *BillingService) GetPaymentAdapter(customer billingModel.CreateCustomer) (*billingModel.Payment, error) {
	p, err := billingModel.GetPaymentAdapter(customer.Identify)

	if err != nil {
		return nil, appError.InvalidPaymentMethod
	}

	return &billingModel.Payment{
		Identify:       customer.Identify,
		CustomerParams: customer.CustomerParams,
		PaymentMethod:  p,
	}, err
}
