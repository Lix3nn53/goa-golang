package authService

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/userRepository"
	"goa-golang/internal/logger"
)

//UserServiceInterface define the user service interface methods
type AuthServiceInterface interface {
	GoogleOauth2(code string) (user *userModel.User, err error)
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
func (s *AuthService) GoogleOauth2(code string) (user *userModel.User, err error) {
	tokenUrl := "https://oauth2.googleapis.com/token"
	userProfileUrl := "https://www.googleapis.com/oauth2/v2/userinfo"

	client := &http.Client{Timeout: 10 * time.Second}

	// Access token request
	data := map[string]interface{}{
		"code":          code,
		"client_id":     os.Getenv("GOOGLE_CLIENT_ID"),
		"client_secret": os.Getenv("GOOGLE_CLIENT_SECRET"),
		"redirect_uri":  os.Getenv("GOOGLE_REDIRECT_URI"),
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
	// print response
	j, err := json.MarshalIndent(accessTokenResponseJson, "", "\t")
	if err != nil {
		return nil, err
	}
	s.logger.Infof(string(j))
	// check for error field in json
	if accessTokenResponseJson["error"] != nil {
		err = errors.New(accessTokenResponseJson["error"].(string))
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

	j, err = json.MarshalIndent(userInfoResponseJson, "", "\t")
	if err != nil {
		return nil, err
	}
	s.logger.Infof(string(j))

	// FIND USER IF EXISTS
	user, err = s.userRepo.FindByID(userId)

	if errors.Is(err, sql.ErrNoRows) {
		// REGISTER USER IF DOES NOT EXIST
		err := s.userRepo.CreateUUID(userId)
		if err != nil {
			if !strings.Contains(err.Error(), "Error 1062: Duplicate entry") { // UUID is in goa_player but not in goa_player_web so lets skip to CreateWebData
				s.logger.Error(err.Error())
				return nil, err
			}
		}

		name := userInfoResponseJson["name"].(string)
		email := userInfoResponseJson["email"].(string)

		userModel := userModel.CreateUser{
			Email:      email,
			McUsername: name,
			Credits:    0,
		}
		user, err = s.userRepo.CreateWebData(userId, userModel)
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
