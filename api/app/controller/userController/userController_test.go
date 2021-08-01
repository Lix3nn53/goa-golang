package userController

import (
	error2 "goa-golang/app/error"
	"goa-golang/app/model"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"goa-golang/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMicroservice_Find(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserServiceInterface(ctrl)

	apiLogger := logger.NewAPILogger()
	apiLogger.InitLogger()

	userController := NewUserController(userUC, apiLogger)

	reqValue := &model.CreateUser{
		Name:       "FirstName",
		Cif:        "email@gmail.com",
		Country:    "es",
		PostalCode: "es",
	}

	t.Run("Correct", func(t *testing.T) {
		userRes := &model.User{
			ID:         1,
			Name:       reqValue.Name,
			Cif:        reqValue.Cif,
			Country:    reqValue.Country,
			PostalCode: reqValue.PostalCode,
		}

		userUC.EXPECT().FindByID(1).Return(userRes, nil)

		router := gin.Default()
		router.GET("/api/users/:id", userController.Find)
		ts := httptest.NewServer(router)
		defer ts.Close()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/1", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Incorrect", func(t *testing.T) {
		userUC.EXPECT().FindByID(2).Return(nil, error2.ErrNotFound)

		router := gin.Default()
		router.GET("/api/users/:id", userController.Find)
		ts := httptest.NewServer(router)
		defer ts.Close()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/2", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Incorrect_2", func(t *testing.T) {

		router := gin.Default()
		router.GET("/api/users/:id", userController.Find)
		ts := httptest.NewServer(router)
		defer ts.Close()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/pa", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestNewUserController(t *testing.T) {
	type args struct {
		service userService.UserServiceInterface
		logger  logger.Logger
	}
	tests := []struct {
		name string
		args args
		want UserControllerInterface
	}{
		{
			name: "success",
			args: args{
				service: nil,
				logger:  nil,
			},
			want: &UserController{
				service: nil,
				logger:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserController(tt.args.service, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User controller = %v, want %v", got, tt.want)
			}
		})
	}
}
