package authRepository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	appError "goa-golang/app/error"
	"goa-golang/app/model/userModel"
	"goa-golang/internal/storage"
	"net/http"
	"os"
	"time"
)

// billingRepository handles communication with the user store
type authRepository struct {
	db *storage.DbStore
}

//UserRepositoryInterface define the user repository interface methods
type AuthRepositoryInterface interface {
	GoogleOauth2() (user *userModel.User, err error)
}

// NewUserRepository implements the user repository interface.
func NewAuthRepository(db *storage.DbStore) AuthRepositoryInterface {
	return &authRepository{
		db,
	}
}

// FindByID implements the method to find a user from the store
func (r *authRepository) GoogleOauth2() (user *userModel.User, err error) {
	tokenUrl := "https://oauth2.googleapis.com/token"
	userProfileUrl := "https://www.googleapis.com/oauth2/v2/userinfo"

	client := &http.Client{Timeout: 10 * time.Second}

	// Access token request
	data := map[string]interface{}{
		"code":          "http://auth.xboxlive.com",
		"client_id":     os.Getenv("GOOGLE_CLIENT_ID"),
		"client_secret": os.Getenv("GOOGLE_CLIENT_SECRET"),
		"redirect_uri":  "http://localhost:3000/login/oauth2/callback",
		"grant_type":    "authorization_code",
	}

	json_data, err2 := json.Marshal(data)
	if err2 != nil {
		return nil, err2
	}

	accessTokenRequest, _ := http.NewRequest(http.MethodPost, tokenUrl, bytes.NewReader(json_data))
	accessTokenRequest.Header.Set("Content-Type", "application/json")
	accessTokenRequest.Header.Set("Accept", "application/json")

	// Access token response
	accessTokenResponse, err2 := client.Do((accessTokenRequest))
	if err2 != nil {
		return nil, err2
	}
	defer accessTokenResponse.Body.Close()

	var accessTokenResponseJson map[string]interface{}
	err2 = json.NewDecoder(accessTokenResponse.Body).Decode(&accessTokenResponseJson)
	if err2 != nil {
		return nil, err2
	}
	accessToken := accessTokenResponseJson["access_token"].(string)

	// User info request
	userInfoRequest, _ := http.NewRequest(http.MethodGet, userProfileUrl, nil)
	userInfoRequest.Header.Set("Authorization", "Bearer "+accessToken)

	// Access token response
	userInfoResponse, err2 := client.Do((userInfoRequest))
	if err2 != nil {
		return nil, err2
	}
	defer userInfoResponse.Body.Close()

	var userInfoResponseJson map[string]interface{}
	err2 = json.NewDecoder(userInfoResponse.Body).Decode(&userInfoResponseJson)
	if err2 != nil {
		return nil, err2
	}
	userId := userInfoResponseJson["id"].(string)

	// TODO
	user = &userModel.User{}

	var query = "SELECT uuid, email, mc_username, credits FROM goa_player_web WHERE id = $1"
	row := r.db.QueryRow(query, userId)

	if err := row.Scan(&user.UUID, &user.Email, &user.McUsername, &user.Credits); err != nil {
		if err == sql.ErrNoRows {
			return nil, appError.ErrNotFound
		}

		return nil, err
	}

	return user, nil
}
