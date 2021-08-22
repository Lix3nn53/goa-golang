package authService

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"goa-golang/app/model/playerModel"
	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/playerRepository"
	"goa-golang/app/repository/userRepository"
	"goa-golang/internal/logger"

	"github.com/golang-jwt/jwt"
)

//UserServiceInterface define the user service interface methods
type AuthServiceInterface interface {
	TokenBuildAccess(uuid string) (tokenString string, err error)
	tokenBuildRefresh(uuid string) (tokenString string, err error)
	TokenValidate(tokenString string, secret string) (userUUID string, err error)
	TokenValidateRefresh(tokenString string) (userUUID string, err error)
	Logout(uuid string, refreshToken string) error
	GoogleOauth2(code string) (refreshToken string, accessToken string, err error)
}

// billingService handles communication with the user repository
type AuthService struct {
	playerRepo playerRepository.PlayerRepositoryInterface
	userRepo   userRepository.UserRepositoryInterface
	logger     logger.Logger
}

// NewUserService implements the user service interface.
func NewAuthService(playerRepo playerRepository.PlayerRepositoryInterface,
	userRepo userRepository.UserRepositoryInterface,
	logger logger.Logger) AuthServiceInterface {
	return &AuthService{
		playerRepo,
		userRepo,
		logger,
	}
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) TokenBuildAccess(uuid string) (tokenString string, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    uuid,
		ExpiresAt: time.Now().Add(time.Duration(1) * time.Minute).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	hmacSampleSecret := []byte(os.Getenv("ACCESS_SECRET"))
	tokenString, err = token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) tokenBuildRefresh(uuid string) (tokenString string, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    uuid,
		ExpiresAt: time.Now().AddDate(0, 1, 0).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	hmacSampleSecret := []byte(os.Getenv("REFRESH_SECRET"))
	tokenString, err = token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}

	err = s.userRepo.AddSession(uuid, tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) TokenValidate(tokenString string, secret string) (userUUID string, err error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		hmacSampleSecret := []byte(secret)
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		fmt.Println("TOKEN DEBUG", claims.Issuer, claims.ExpiresAt)

		return claims.Issuer, nil
	} else {
		return "", err
	}
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) TokenValidateRefresh(tokenString string) (userUUID string, err error) {
	userUUID, err = s.TokenValidate(tokenString, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		return "", err
	}

	sessions, err := s.userRepo.GetSessions(userUUID)
	if err != nil {
		return "", err
	}

	split := strings.Split(sessions, "/")

	contains := false
	for _, value := range split {
		if value == tokenString {
			contains = true
			break
		}
	}

	if !contains {
		return "", errors.New("session is not active")
	}

	return userUUID, nil
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) Logout(uuid string, refreshToken string) error {
	err := s.userRepo.RemoveSession(uuid, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) GoogleOauth2(code string) (refreshToken string, accessToken string, err error) {
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
		return "", "", err
	}

	accessTokenRequest, _ := http.NewRequest(http.MethodPost, tokenUrl, bytes.NewReader(json_data))
	accessTokenRequest.Header.Set("Content-Type", "application/json")
	accessTokenRequest.Header.Set("Accept", "application/json")

	// Access token response
	accessTokenResponse, err := client.Do((accessTokenRequest))
	if err != nil {
		return "", "", err
	}
	defer accessTokenResponse.Body.Close()

	var accessTokenResponseJson map[string]interface{}
	err = json.NewDecoder(accessTokenResponse.Body).Decode(&accessTokenResponseJson)
	if err != nil {
		return "", "", err
	}
	// print response
	j, err := json.MarshalIndent(accessTokenResponseJson, "", "\t")
	if err != nil {
		return "", "", err
	}
	s.logger.Infof(string(j))
	// check for error field in json
	if accessTokenResponseJson["error"] != nil {
		err = errors.New(accessTokenResponseJson["error"].(string))
		return "", "", err
	}
	accessTokenGoogle := accessTokenResponseJson["access_token"].(string)

	// User info request
	userInfoRequest, _ := http.NewRequest(http.MethodGet, userProfileUrl, nil)
	userInfoRequest.Header.Set("Authorization", "Bearer "+accessTokenGoogle)

	// Access token response
	userInfoResponse, err := client.Do((userInfoRequest))
	if err != nil {
		return "", "", err
	}
	defer userInfoResponse.Body.Close()

	var userInfoResponseJson map[string]interface{}
	err = json.NewDecoder(userInfoResponse.Body).Decode(&userInfoResponseJson)
	if err != nil {
		return "", "", err
	}
	userId := userInfoResponseJson["id"].(string)

	j, err = json.MarshalIndent(userInfoResponseJson, "", "\t")
	if err != nil {
		return "", "", err
	}
	s.logger.Infof(string(j))

	// FIND USER IF EXISTS
	user, err := s.userRepo.FindByID(userId)

	if errors.Is(err, sql.ErrNoRows) {
		// REGISTER USER IF DOES NOT EXIST
		createPlayer := playerModel.CreatePlayer{
			UUID: userId,
		}

		err := s.playerRepo.CreateUUID(createPlayer)
		if err != nil {
			if !strings.Contains(err.Error(), "Error 1062: Duplicate entry") { // UUID is in goa_player but not in goa_player_web so lets skip to CreateWebData
				s.logger.Error(err.Error())
				return "", "", err
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
			return "", "", err
		}
	} else if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	refreshToken, err = s.tokenBuildRefresh(user.UUID)
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	accessToken, err = s.TokenBuildAccess(user.UUID)
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	return refreshToken, accessToken, nil
}
