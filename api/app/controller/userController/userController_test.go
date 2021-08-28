package userController

import (
	appError "goa-golang/app/error"
	"goa-golang/app/model/userModel"
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

	reqValue := &userModel.CreateUserMicrosoft{
		UUID: "1",
	}

	t.Run("Correct", func(t *testing.T) {
		userRes := &userModel.User{
			UUID: reqValue.UUID,
		}

		userUC.EXPECT().FindByID(1, "uuid").Return(userRes, nil)

		router := gin.Default()
		router.GET("/api/users/:id", userController.Info)
		ts := httptest.NewServer(router)
		defer ts.Close()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/1", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Incorrect", func(t *testing.T) {
		userUC.EXPECT().FindByID(2, "uuid").Return(nil, appError.ErrNotFound)

		router := gin.Default()
		router.GET("/api/users/:id", userController.Info)
		ts := httptest.NewServer(router)
		defer ts.Close()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/2", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Incorrect_2", func(t *testing.T) {

		router := gin.Default()
		router.GET("/api/users/:id", userController.Info)
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
