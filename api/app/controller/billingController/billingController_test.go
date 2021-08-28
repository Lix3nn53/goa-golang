package billingController

import (
	"bytes"
	"encoding/json"
	appError "goa-golang/app/error"
	"goa-golang/app/model/billingModel"
	"goa-golang/app/model/userModel"
	"goa-golang/app/service/billingService"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"goa-golang/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBillingController_Store(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserServiceInterface(ctrl)
	billUC := mock.NewMockBillingServiceInterface(ctrl)

	// payUc := mock.NewMockPaymentAdapterCase(ctrl)

	apiLogger := logger.NewAPILogger()
	apiLogger.InitLogger()

	billingController := NewBillingController(billUC, userUC, apiLogger)

	t.Run("UserNotFound", func(t *testing.T) {
		userUC.EXPECT().FindByID(1, "uuid").Return(nil, appError.ErrNotFound)

		router := gin.Default()
		router.POST("/api/users/:id/paypal", billingController.AddCustomer)
		ts := httptest.NewServer(router)
		defer ts.Close()
		w := httptest.NewRecorder()

		var customer = billingModel.CreateCustomer{
			Identify: billingModel.AccountStripe,
			CustomerParams: billingModel.CustomerParams{
				Email: "test3@test.com",
				Desc:  "a 3rd test customer",
				Card: &billingModel.CardParams{
					Name:     "serRes.Name",
					Number:   "userRes.Cif",
					ExpYear:  time.Now().Year() + 1,
					ExpMonth: 1,
				},
			},
		}

		body, _ := json.Marshal(customer)
		req := httptest.NewRequest("POST", "/api/users/1/paypal", bytes.NewBuffer(body))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestBillingController_Store2(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserServiceInterface(ctrl)
	billUC := mock.NewMockBillingServiceInterface(ctrl)

	// payUc := mock.NewMockPaymentAdapterCase(ctrl)

	apiLogger := logger.NewAPILogger()
	apiLogger.InitLogger()

	billingController := NewBillingController(billUC, userUC, apiLogger)

	t.Run("InvalidIdentifyPayment", func(t *testing.T) {

		userExpected := userModel.User{
			UUID:    "1",
			Email:   "a",
			Credits: 13,
		}

		userRes := &userModel.User{
			UUID:    userExpected.UUID,
			Email:   userExpected.Email,
			Credits: userExpected.Credits,
		}

		userUC.EXPECT().FindByID(1, "uuid").Return(userRes, nil)

		const badIdentify billingModel.Identify = "strbadIdentifyipe"

		var customer = billingModel.CreateCustomer{
			Identify: badIdentify,
			CustomerParams: billingModel.CustomerParams{
				Email: "11@test.com",
				Desc:  "a 3rd test customer",
				Card: &billingModel.CardParams{
					// TODO fix uuids to actual numbers
					Number:   userRes.UUID,
					ExpYear:  time.Now().Year() + 1,
					ExpMonth: 1,
				},
			},
		}

		billUC.EXPECT().GetPaymentAdapter(customer).Return(nil, appError.ErrInvalidPaymentMethod)

		router := gin.Default()
		router.POST("/api/users/:id/paypal", billingController.AddCustomer)

		ts := httptest.NewServer(router)
		defer ts.Close()

		w := httptest.NewRecorder()
		body, _ := json.Marshal(customer)
		req := httptest.NewRequest("POST", "/api/users/1/paypal", bytes.NewBuffer(body))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}

func TestBillingController_Store3(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserServiceInterface(ctrl)
	billUC := mock.NewMockBillingServiceInterface(ctrl)

	// payUc := mock.NewMockPaymentAdapterCase(ctrl)

	apiLogger := logger.NewAPILogger()
	apiLogger.InitLogger()

	billingController := NewBillingController(billUC, userUC, apiLogger)

	t.Run("Correct", func(t *testing.T) {
		userRes := &userModel.User{
			UUID:    "1",
			Email:   "a",
			Credits: 14,
		}

		userUC.EXPECT().FindByID(1, "uuid").Return(userRes, nil)

		var customer = billingModel.CreateCustomer{
			Identify: billingModel.AccountStripe,
			CustomerParams: billingModel.CustomerParams{
				Email: "test3@test.com",
				Desc:  "a 3rd test customer",
				Card: &billingModel.CardParams{
					Number:   userRes.UUID,
					ExpYear:  time.Now().Year() + 1,
					ExpMonth: 1,
				},
			},
		}

		p := &billingModel.Payment{
			Identify:       customer.Identify,
			CustomerParams: customer.CustomerParams,
			PaymentMethod:  &mock.FakeAdapter{},
		}

		billUC.EXPECT().GetPaymentAdapter(customer).Return(p, nil)
		billUC.EXPECT().AddBilling(*userRes, *p).Return(nil)

		router := gin.Default()
		router.POST("/api/users/:id/paypal", billingController.AddCustomer)
		ts := httptest.NewServer(router)
		defer ts.Close()
		w := httptest.NewRecorder()

		body, _ := json.Marshal(customer)
		req := httptest.NewRequest("POST", "/api/users/1/paypal", bytes.NewBuffer(body))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNewBillingController(t *testing.T) {
	type args struct {
		uservice userService.UserServiceInterface
		service  billingService.BillingServiceInterface
		logger   logger.Logger
	}
	tests := []struct {
		name string
		args args
		want BillingControllerInterface
	}{
		{
			name: "success",
			args: args{
				service:  nil,
				uservice: nil,
				logger:   nil,
			},
			want: &BillingController{
				service:  nil,
				uservice: nil,
				logger:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillingController(tt.args.service, tt.args.uservice, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User controller = %v, want %v", got, tt.want)
			}
		})
	}
}
