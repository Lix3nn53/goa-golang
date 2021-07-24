package repository

import (
	"goa-golang/app/model"
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
			want: &userRepository{
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

	columns := []string{"id", "email", "cif", "postal_code", "country"}
	userID := int(1)
	mockUser := &model.User{
		ID:         userID,
		Name:       "FirstName",
		Cif:        "email@gmail.com",
		Country:    "es",
		PostalCode: "es",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		userID,
		mockUser.Cif,
		mockUser.Name,
		mockUser.PostalCode,
		mockUser.Country,
	)

	mock.ExpectQuery("SELECT id, cif, name, postal_code, country FROM users WHERE id = $1").WithArgs(userID).WillReturnRows(rows)

	foundUser, err := userPGRepository.FindByID(mockUser.ID)

	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.ID, userID)
}

func TestUserRepository_FindByID_IncorrectID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "db")
	defer sqlxDB.Close()

	userPGRepository := NewUserRepository(&storage.DbStore{DB: sqlxDB})

	columns := []string{"id", "cif", "name", "postal_code", "country"}
	userID := int(1)
	mockUser := &model.User{
		ID:         userID,
		Name:       "FirstName",
		Cif:        "email@gmail.com",
		Country:    "es",
		PostalCode: "es",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		userID,
		mockUser.Cif,
		mockUser.Name,
		mockUser.PostalCode,
		mockUser.Country,
	)

	mock.ExpectQuery("SELECT id, cif, name, postal_code, country FROM users WHERE id = $1").WithArgs(2).WillReturnRows(rows)

	foundUser, err := userPGRepository.FindByID(mockUser.ID)

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

	userID := int(1)
	mockUser := model.CreateUser{
		Name:       "FirstName",
		Cif:        "LastName",
		Country:    "es",
		PostalCode: "es",
	}

	ep := mock.ExpectPrepare("INSERT INTO users (name, cif, postal_code, country) VALUES ($1, $2, $3, $4) RETURNING id").WillBeClosed()
	ep.ExpectQuery().WithArgs(mockUser.Name, mockUser.Cif, mockUser.PostalCode, mockUser.Country).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	foundUser, err := userPGRepository.Create(mockUser)
	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.ID, userID)
}
