package userService

import (
	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/userRepository"
	"goa-golang/internal/logger"
	"goa-golang/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		userRepository userRepository.UserRepositoryInterface
	}
	tests := []struct {
		name string
		args args
		want UserServiceInterface
	}{
		{
			name: "success",
			args: args{
				userRepository: nil,
			},
			want: &UserService{
				userRepo: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.userRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_Store(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userR := mock.NewMockUserRepositoryInterface(ctrl)
	userService := NewUserService(userR)

	reqValue := userModel.CreateUser{
		Email:      "a@a.com",
		McUsername: "a",
		Credits:    7,
	}

	t.Run("Store", func(t *testing.T) {
		t.Parallel()

		user := userModel.CreateUser{
			Email:      reqValue.Email,
			McUsername: reqValue.McUsername,
			Credits:    reqValue.Credits,
		}

		userID := string("1")
		userRes := &userModel.User{
			UUID:       userID,
			Email:      user.Email,
			McUsername: user.McUsername,
			Credits:    user.Credits,
		}
		var err error

		userR.EXPECT().Create(userID, user).Return(userRes, err)

		response, err := userService.Store(userID, reqValue)

		require.NoError(t, err)
		require.NotNil(t, response)

		logger := logger.NewAPILogger()
		logger.InitLogger()
		logger.Info(response)
	})
}
