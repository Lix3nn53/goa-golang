package userRepository

import (
	"goa-golang/app/model/userModel"
	"goa-golang/internal/storage"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestUserRepositoryInit(t *testing.T) {
	type args struct {
		db *storage.DbStore
	}
	tests := []struct {
		name string
		args args
		want UserRepositoryInterface
	}{
		{
			name: "success",
			args: args{
				db: nil,
			},
			want: &UserRepository{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "db")
	defer sqlxDB.Close()

	userPGRepository := NewUserRepository(&storage.DbStore{DB: sqlxDB})

	columns := []string{"uuid", "email", "credits"}
	userID := string("1")
	mockUser := &userModel.User{
		UUID:    userID,
		Email:   "FirstName",
		Credits: 6,
	}

	rows := sqlmock.NewRows(columns).AddRow(
		mockUser.UUID,
		mockUser.Email,
		mockUser.Credits,
	)

	mock.ExpectQuery("SELECT uuid, email, credits FROM goa_player_web WHERE id = ?").WithArgs(userID).WillReturnRows(rows)

	foundUser, err := userPGRepository.FindByID(mockUser.UUID, "uuid")

	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.UUID, userID)
}

func TestUserRepository_FindByID_IncorrectID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "db")
	defer sqlxDB.Close()

	userPGRepository := NewUserRepository(&storage.DbStore{DB: sqlxDB})

	columns := []string{"uuid", "email", "credits"}
	userID := string("1")
	mockUser := &userModel.User{
		UUID:    userID,
		Email:   "email@gmail.com",
		Credits: 12,
	}

	rows := sqlmock.NewRows(columns).AddRow(
		userID,
		mockUser.Email,
		mockUser.Credits,
	)

	mock.ExpectQuery("SELECT uuid, email, credits FROM goa_player_web WHERE id = ?").WithArgs(2).WillReturnRows(rows)

	foundUser, err := userPGRepository.FindByID(mockUser.UUID, "uuid")

	require.Error(t, err)
	require.Nil(t, foundUser)
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "db")
	defer sqlxDB.Close()

	userPGRepository := NewUserRepository(&storage.DbStore{DB: sqlxDB})

	userID := string("1")
	mockUser := userModel.CreateUserMicrosoft{
		UUID: userID,
	}

	ep := mock.ExpectPrepare("INSERT INTO goa_player_web (uuid) VALUES (?)").WillBeClosed()
	ep.ExpectQuery().WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(userID))

	foundUser, err := userPGRepository.CreateWithMicrosoft(mockUser)
	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.UUID, userID)
}
