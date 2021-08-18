package billingRepository

import (
	"goa-golang/app/model/billingModel"
	"goa-golang/internal/storage"
)

// billingRepository handles communication with the user store
type billingRepository struct {
	db *storage.DbStore
}

//BillingRepositoryInterface define the user repository interface methods
type BillingRepositoryInterface interface {
	CreateBillingService(identity billingModel.Identify, key string, userID string) error
}

// NewBillingRepository implements the billing repository interface.
func NewBillingRepository(db *storage.DbStore) BillingRepositoryInterface {
	return &billingRepository{
		db,
	}
}

// CreateBillingService Create implements the method to persist a Payment user
func (r *billingRepository) CreateBillingService(identify billingModel.Identify, PaymentUserKey string, userID string) error {
	createUserQuery := `INSERT INTO billing (identify, key, user_id) 
		VALUES (?, ?, ?)
		RETURNING id`

	stmt, err := r.db.Prepare(createUserQuery)
	if err != nil {
		return nil
	}
	defer stmt.Close()

	var paymentID int
	err = stmt.QueryRow(identify, PaymentUserKey, userID).Scan(&paymentID)
	return err
}
