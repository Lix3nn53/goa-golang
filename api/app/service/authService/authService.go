package authService

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	MinecraftOauth2(code string) (refreshToken string, accessToken string, err error)
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

func (s *AuthService) onSuccessfulOauth(uuid string, userModel userModel.CreateUser) (refreshToken string, accessToken string, err error) {
	// FIND USER IF EXISTS
	user, err := s.userRepo.FindByID(uuid)

	if errors.Is(err, sql.ErrNoRows) {
		// REGISTER USER IF DOES NOT EXIST
		createPlayer := playerModel.CreatePlayer{
			UUID: uuid,
		}

		err := s.playerRepo.CreateUUID(createPlayer)
		if err != nil {
			if !strings.Contains(err.Error(), "Error 1062: Duplicate entry") { // UUID is in goa_player but not in goa_player_web so lets skip to CreateWebData
				s.logger.Error(err.Error())
				return "", "", err
			}
		}

		user, err = s.userRepo.CreateWebData(uuid, userModel)
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
		"redirect_uri":  os.Getenv("OAUTH_REDIRECT_URI"),
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

	name := userInfoResponseJson["name"].(string)
	email := userInfoResponseJson["email"].(string)

	userModel := userModel.CreateUser{
		Email:      email,
		McUsername: name,
		Credits:    0,
	}
	refreshToken, accessToken, err = s.onSuccessfulOauth(userId, userModel)
	if err != nil {
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) MinecraftOauth2(code string) (refreshToken string, accessToken string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	accessTokenMicrosoft, err := s.microsoftToken(client, code)
	if err != nil {
		return "", "", err
	}

	xblToken, uhs, err := s.microsoftXBL(client, accessTokenMicrosoft)
	if err != nil {
		return "", "", err
	}

	xstsToken, err := s.microsoftXSTS(client, xblToken)
	if err != nil {
		return "", "", err
	}

	minecraftAccessToken, err := s.microsoftMinecraftAuth(client, uhs, xstsToken)
	if err != nil {
		return "", "", err
	}

	uuid, name, err := s.microsoftMinecraftProfile(client, minecraftAccessToken)
	if err != nil {
		return "", "", err
	}

	userModel := userModel.CreateUser{
		McUsername: name,
		Credits:    0,
	}
	refreshToken, accessToken, err = s.onSuccessfulOauth(uuid, userModel)
	if err != nil {
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

func (s *AuthService) microsoftToken(client *http.Client, code string) (accessTokenMicrosoft string, err error) {
	reqUrl := "https://login.microsoftonline.com/consumers/oauth2/v2.0/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", os.Getenv("MICROSOFT_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("MICROSOFT_CLIENT_SECRET"))
	data.Set("redirect_uri", os.Getenv("OAUTH_REDIRECT_URI"))

	request, _ := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(data.Encode()))
	request.Close = true
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// Access token response
	response, err := client.Do((request))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var reponseJson map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&reponseJson)
	if err != nil {
		return "", err
	}
	// print response
	j, err := json.MarshalIndent(reponseJson, "", "\t")
	if err != nil {
		return "", err
	}
	s.logger.Infof(string(j))
	// check for error field in json
	if reponseJson["error"] != nil {
		err = errors.New(reponseJson["error"].(string))
		return "", err
	}
	accessTokenMicrosoft = reponseJson["access_token"].(string)

	return accessTokenMicrosoft, nil
}

func (s *AuthService) microsoftXBL(client *http.Client, accessTokenMicrosoft string) (xblToken string, uhs string, err error) {
	s.logger.Infof("microsoftXBL START")
	url := "https://user.auth.xboxlive.com/user/authenticate"

	// Xbox live auth request
	data := map[string]interface{}{
		"RelyingParty": "http://auth.xboxlive.com",
		"TokenType":    "JWT",
		"Properties": map[string]interface{}{
			"AuthMethod": "RPS",
			"SiteName":   "user.auth.xboxlive.com",
			"RpsTicket":  "d=" + accessTokenMicrosoft,
		},
	}

	json_data, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(json_data))
	request.Close = true
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("x-xbl-contract-version", "1")

	// Xbox live auth response
	response, err := client.Do((request))
	if err != nil {
		s.logger.Errorf(err.Error())
		return "", "", err
	}
	defer response.Body.Close()

	s.logger.Infof(response.Status)

	var responseJson map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJson)
	if err != nil {
		s.logger.Errorf(err.Error())
		return "", "", err
	}

	j, err := json.MarshalIndent(responseJson, "", "\t")
	if err != nil {
		s.logger.Errorf(err.Error())
		return "", "", err
	}
	s.logger.Infof(string(j))
	xblToken = responseJson["Token"].(string)
	uhs = responseJson["DisplayClaims"].(map[string]interface{})["xui"].([]interface{})[0].(map[string]interface{})["uhs"].(string)

	return xblToken, uhs, nil
}

func (s *AuthService) microsoftXSTS(client *http.Client, xblToken string) (xstsToken string, err error) {
	url := "https://xsts.auth.xboxlive.com/xsts/authorize"

	// Xsts request
	data := map[string]interface{}{
		"Properties": map[string]interface{}{
			"SandboxId": "RETAIL",
			"UserTokens": []string{
				xblToken,
			},
		},
		"RelyingParty": "rp://api.minecraftservices.com/",
		"TokenType":    "JWT",
	}

	json_data, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(json_data))
	request.Close = true
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("x-xbl-contract-version", "1")

	// Xsts response
	response, err := client.Do((request))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseJson map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJson)
	if err != nil {
		return "", err
	}

	j, err := json.MarshalIndent(responseJson, "", "\t")
	if err != nil {
		return "", err
	}
	s.logger.Infof(string(j))
	xstsToken = responseJson["Token"].(string)

	return xstsToken, nil
}

func (s *AuthService) microsoftMinecraftAuth(client *http.Client, uhs string, xstsToken string) (minecraftAccessToken string, err error) {
	url := "https://api.minecraftservices.com/authentication/login_with_xbox"

	// Xsts request
	data := map[string]interface{}{
		"identityToken":       "XBL3.0 x=" + uhs + ";" + xstsToken,
		"ensureLegacyEnabled": true,
	}

	json_data, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(json_data))
	request.Close = true
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Xsts response
	response, err := client.Do((request))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseJson map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJson)
	if err != nil {
		return "", err
	}

	j, err := json.MarshalIndent(responseJson, "", "\t")
	if err != nil {
		return "", err
	}
	s.logger.Infof(string(j))
	minecraftAccessToken = responseJson["access_token"].(string)

	return minecraftAccessToken, nil
}

/* func (s *AuthService) microsoftMinecraftOwnership(client *http.Client, minecraftAccessToken string) (owner bool, err error) {
	url := "https://api.minecraftservices.com/entitlements/mcstore"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Close = true
	request.Header.Set("Authorization", "Bearer "+minecraftAccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Xsts response
	response, err := client.Do((request))
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	var responseJson map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJson)
	if err != nil {
		return false, err
	}

	j, err := json.MarshalIndent(responseJson, "", "\t")
	if err != nil {
		return false, err
	}
	s.logger.Infof(string(j))
	minecraftAccessToken = responseJson["access_token"].(string)

	return true, nil
} */

func (s *AuthService) microsoftMinecraftProfile(client *http.Client, minecraftAccessToken string) (uuid string, name string, err error) {
	url := "https://api.minecraftservices.com/minecraft/profile"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Close = true
	request.Header.Set("Authorization", "Bearer "+minecraftAccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Xsts response
	response, err := client.Do((request))
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	var responseJson map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJson)
	if err != nil {
		return "", "", err
	}

	j, err := json.MarshalIndent(responseJson, "", "\t")
	if err != nil {
		return "", "", err
	}
	s.logger.Infof(string(j))

	if _, ok := responseJson["error"]; ok {
		return "", "", errors.New("your xbox account does not own minecraft")
	}

	uuid = responseJson["id"].(string)
	name = responseJson["name"].(string)

	return uuid, name, nil
}
