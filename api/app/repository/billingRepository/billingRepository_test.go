package billingRepository

import (
	"goa-golang/app/model/billingModel"
	"goa-golang/internal/storage"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestBillingRepositoryInit(t *testing.T) {
	type args struct {
		db *storage.DbStore
	}
	tests := []struct {
		name string
		args args
		want BillingRepositoryInterface
	}{
		{
			name: "success",
			args: args{
				db: nil,
			},
			want: &billingRepository{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillingRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBillingRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "db")
	defer sqlxDB.Close()

	userPGRepository := NewBillingRepository(&storage.DbStore{DB: sqlxDB})

	userID := int(1) // payment model.Payment, PaymentUserKey string, userID int
	key := "cus_124"

	serviceID := int(1)

	ep := mock.ExpectPrepare("INSERT INTO billing (identify, key, user_id) VALUES ($1, $2, $3) RETURNING id").WillBeClosed()
	ep.ExpectQuery().WithArgs(billingModel.AccountStripe, key, userID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(serviceID))

	err = userPGRepository.CreateBillingService(billingModel.AccountStripe, key, userID)
	require.NoError(t, err)
	require.Equal(t, 1, userID)
}
