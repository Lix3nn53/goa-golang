package authService

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/userRepository"
	"goa-golang/internal/logger"
)

//UserServiceInterface define the user service interface methods
type AuthServiceInterface interface {
	GoogleOauth2() (user *userModel.User, err error)
}

// billingService handles communication with the user repository
type AuthService struct {
	userRepo userRepository.UserRepositoryInterface
	logger   logger.Logger
}

// NewUserService implements the user service interface.
func NewAuthService(userRepo userRepository.UserRepositoryInterface, logger logger.Logger) AuthServiceInterface {
	return &AuthService{
		userRepo,
		logger,
	}
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) GoogleOauth2() (user *userModel.User, err error) {
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

	json_data, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	accessTokenRequest, _ := http.NewRequest(http.MethodPost, tokenUrl, bytes.NewReader(json_data))
	accessTokenRequest.Header.Set("Content-Type", "application/json")
	accessTokenRequest.Header.Set("Accept", "application/json")

	// Access token response
	accessTokenResponse, err := client.Do((accessTokenRequest))
	if err != nil {
		return nil, err
	}
	defer accessTokenResponse.Body.Close()

	var accessTokenResponseJson map[string]interface{}
	err = json.NewDecoder(accessTokenResponse.Body).Decode(&accessTokenResponseJson)
	if err != nil {
		return nil, err
	}
	accessToken := accessTokenResponseJson["access_token"].(string)

	// User info request
	userInfoRequest, _ := http.NewRequest(http.MethodGet, userProfileUrl, nil)
	userInfoRequest.Header.Set("Authorization", "Bearer "+accessToken)

	// Access token response
	userInfoResponse, err := client.Do((userInfoRequest))
	if err != nil {
		return nil, err
	}
	defer userInfoResponse.Body.Close()

	var userInfoResponseJson map[string]interface{}
	err = json.NewDecoder(userInfoResponse.Body).Decode(&userInfoResponseJson)
	if err != nil {
		return nil, err
	}
	userId := userInfoResponseJson["id"].(string)

	// FIND USER IF EXISTS
	user, err = s.userRepo.FindByID(userId)

	if err == sql.ErrNoRows {
		// REGISTER USER IF DOES NOT EXIST
		name := userInfoResponseJson["name"].(string)

		userModel := userModel.CreateUser{
			Email:      "test@test.com",
			McUsername: name,
			Credits:    0,
		}

		user, err = s.userRepo.Create(userId, userModel)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
	} else if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}
