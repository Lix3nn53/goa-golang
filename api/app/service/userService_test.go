package service

import (
	"goa-golang/app/model"
	"goa-golang/app/repository"
	"goa-golang/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		userRepository repository.UserRepositoryInterface
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
			want: &userService{
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
	userR := mock.NewMockUserPGRepository(ctrl)
	userService := NewUserService(userR)

	reqValue := model.CreateUser{
		Name:       "a",
		Cif:        "a@a.com",
		Country:    "a",
		PostalCode: "a",
	}

	t.Run("Store", func(t *testing.T) {
		t.Parallel()

		user := model.CreateUser{
			Name:       reqValue.Name,
			Cif:        reqValue.Cif,
			Country:    reqValue.Country,
			PostalCode: reqValue.PostalCode,
		}

		userID := int(1)
		userRes := &model.User{
			ID:         userID,
			Name:       user.Name,
			Cif:        user.Cif,
			Country:    user.Country,
			PostalCode: user.PostalCode,
		}
		var err error

		userR.EXPECT().Create(user).Return(userRes, err)

		response, err := userService.Store(reqValue)

		require.NoError(t, err)
		require.NotNil(t, response)
	})
}
