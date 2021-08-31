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

type MyCustomClaims struct {
	IssuerField string `json:"iss_field,omitempty"`
	jwt.StandardClaims
}

//UserServiceInterface define the user service interface methods
type AuthServiceInterface interface {
	TokenBuildAccess(id string, idField string) (tokenString string, err error)
	tokenBuildRefresh(id string, idField string) (tokenString string, err error)
	TokenValidate(tokenString string, secret string) (id string, idField string, err error)
	TokenValidateRefresh(tokenString string) (id string, idField string, err error)
	Logout(id string, idField string, refreshToken string) error
	GoogleOauth2(code string) (refreshToken string, accessToken string, err error)
	TwitchOauth2(code string) (refreshToken string, accessToken string, err error)
	DiscordOauth2(code string) (refreshToken string, accessToken string, err error)
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
func (s *AuthService) TokenBuildAccess(id string, idField string) (tokenString string, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		idField,
		jwt.StandardClaims{
			Issuer:    id,
			ExpiresAt: time.Now().Add(time.Duration(1) * time.Minute).Unix(),
		},
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
func (s *AuthService) tokenBuildRefresh(id string, idField string) (tokenString string, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		idField,
		jwt.StandardClaims{
			Issuer:    id,
			ExpiresAt: time.Now().AddDate(0, 1, 0).Unix(),
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	hmacSampleSecret := []byte(os.Getenv("REFRESH_SECRET"))
	tokenString, err = token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}

	err = s.userRepo.AddSession(id, idField, tokenString)
	if err != nil {
		s.logger.Error(err.Error())
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) TokenValidate(tokenString string, secret string) (id string, idField string, err error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		hmacSampleSecret := []byte(secret)
		return hmacSampleSecret, nil
	})
	if err != nil {
		return "", "", err
	}

	s.logger.Info("TESTTTTT")

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {

		s.logger.Info("claims: %s", claims)

		return claims.Issuer, claims.IssuerField, nil
	} else {
		s.logger.Info("claims: %s", claims)
		return "", "", errors.New("token is not valid")
	}
}

func (s *AuthService) TokenValidateRefresh(tokenString string) (id string, idField string, err error) {
	id, idField, err = s.TokenValidate(tokenString, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		return "", "", err
	}

	sessions := s.userRepo.GetSessions(id, idField)
	if !sessions.Valid {
		return "", "", errors.New("sessions is not valid")
	}

	split := strings.Split(sessions.String, "/")

	contains := false
	for _, value := range split {
		if value == tokenString {
			contains = true
			break
		}
	}

	if !contains {
		return "", "", errors.New("session is not active")
	}

	return id, idField, nil
}

func (s *AuthService) Logout(id string, field string, refreshToken string) error {
	err := s.userRepo.RemoveSession(id, field, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

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
	googleId := userInfoResponseJson["id"].(string)

	j, err = json.MarshalIndent(userInfoResponseJson, "", "\t")
	if err != nil {
		return "", "", err
	}
	s.logger.Infof(string(j))

	email := userInfoResponseJson["email"].(string)

	userModel := userModel.CreateUserGoogle{
		GoogleId: googleId,
		Email:    email,
	}

	// FIND USER IF EXISTS
	user, err := s.userRepo.FindByID(googleId, "google_id")

	if errors.Is(err, sql.ErrNoRows) {
		// REGISTER USER IF DOES NOT EXIST
		user, err = s.userRepo.CreateWithGoogle(userModel)
		if err != nil {
			s.logger.Error(err.Error())
			return "", "", err
		}
	} else if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	s.logger.Info("USER: ", user)

	refreshToken, err = s.tokenBuildRefresh(user.GoogleId, "google_id")
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	accessToken, err = s.TokenBuildAccess(user.GoogleId, "google_id")
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

func (s *AuthService) DiscordOauth2(code string) (refreshToken string, accessToken string, err error) {
	tokenUrl := "https://discord.com/api/v8/oauth2/token"
	userProfileUrl := "https://discord.com/api/v8/oauth2/@me"

	client := &http.Client{Timeout: 10 * time.Second}

	// Access token request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", os.Getenv("DISCORD_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
	data.Set("redirect_uri", os.Getenv("OAUTH_REDIRECT_URI"))

	accessTokenRequest, _ := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(data.Encode()))
	accessTokenRequest.Close = true
	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	accessTokenDiscord := accessTokenResponseJson["access_token"].(string)

	// User info request
	userInfoRequest, _ := http.NewRequest(http.MethodGet, userProfileUrl, nil)
	userInfoRequest.Header.Set("Authorization", "Bearer "+accessTokenDiscord)

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

	j, err = json.MarshalIndent(userInfoResponseJson, "", "\t")
	if err != nil {
		return "", "", err
	}
	s.logger.Infof(string(j))

	discordUser := userInfoResponseJson["user"].(map[string]interface{})

	discordId := discordUser["id"].(string)

	userModel := userModel.CreateUserDiscord{
		DiscordId: discordId,
	}

	// FIND USER IF EXISTS
	user, err := s.userRepo.FindByID(discordId, "discord_id")

	if errors.Is(err, sql.ErrNoRows) {
		// REGISTER USER IF DOES NOT EXIST
		user, err = s.userRepo.CreateWithDiscord(userModel)
		if err != nil {
			s.logger.Error(err.Error())
			return "", "", err
		}
	} else if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	s.logger.Info("USER: ", user)

	refreshToken, err = s.tokenBuildRefresh(user.DiscordId, "discord_id")
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	accessToken, err = s.TokenBuildAccess(user.DiscordId, "discord_id")
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

func (s *AuthService) TwitchOauth2(code string) (refreshToken string, accessToken string, err error) {
	tokenUrl := "https://id.twitch.tv/oauth2/token"
	userProfileUrl := "https://api.twitch.tv/helix/users"

	client := &http.Client{Timeout: 10 * time.Second}

	// Access token request
	// Access token request
	data := map[string]interface{}{
		"code":          code,
		"client_id":     os.Getenv("TWITCH_CLIENT_ID"),
		"client_secret": os.Getenv("TWITCH_CLIENT_SECRET"),
		"redirect_uri":  os.Getenv("OAUTH_REDIRECT_URI"),
		"grant_type":    "authorization_code",
	}

	json_data, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	s.logger.Infof(string(json_data))

	accessTokenRequest, _ := http.NewRequest(http.MethodPost, tokenUrl, bytes.NewReader(json_data))
	accessTokenRequest.Close = true
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
	accessTokenTwitch := accessTokenResponseJson["access_token"].(string)

	// User info request
	userInfoRequest, _ := http.NewRequest(http.MethodGet, userProfileUrl, nil)
	userInfoRequest.Header.Set("Authorization", "Bearer "+accessTokenTwitch)
	userInfoRequest.Header.Set("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))

	// Access token response
	userInfoResponse, err := client.Do((userInfoRequest))
	if err != nil {
		return "", "", err
	}
	defer userInfoResponse.Body.Close()

	var userInfoResponseJsonArray map[string]interface{}
	err = json.NewDecoder(userInfoResponse.Body).Decode(&userInfoResponseJsonArray)
	if err != nil {
		return "", "", err
	}

	j, err = json.MarshalIndent(userInfoResponseJsonArray, "", "\t")
	if err != nil {
		return "", "", err
	}
	s.logger.Infof(string(j))

	twitchUsers := userInfoResponseJsonArray["data"].([]map[string]inte rface{})

	twitchUser := twitchUsers[0]

	twitchId := twitchUser["id"].(string)
	email := twitchUser["email"].(string)

	userModel := userModel.CreateUserTwitch{
		TwitchId: twitchId,
		Email:    email,
	}

	// FIND USER IF EXISTS
	user, err := s.userRepo.FindByID(twitchId, "twitch_id")

	if errors.Is(err, sql.ErrNoRows) {
		// REGISTER USER IF DOES NOT EXIST
		user, err = s.userRepo.CreateWithTwitch(userModel)
		if err != nil {
			s.logger.Error(err.Error())
			return "", "", err
		}
	} else if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	s.logger.Info("USER: ", user)

	refreshToken, err = s.tokenBuildRefresh(user.TwitchId, "twitch_id")
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	accessToken, err = s.TokenBuildAccess(user.TwitchId, "twitch_id")
	if err != nil {
		s.logger.Error(err.Error())
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

	uuid, _, err := s.microsoftMinecraftProfile(client, minecraftAccessToken)
	if err != nil {
		return "", "", err
	}

	userModel := userModel.CreateUserMicrosoft{
		UUID: uuid,
	}

	// FIND USER IF EXISTS
	user, err := s.userRepo.FindByID(uuid, "uuid")

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

		user, err = s.userRepo.CreateWithMicrosoft(userModel)
		if err != nil {
			s.logger.Error(err.Error())
			return "", "", err
		}
	} else if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	refreshToken, err = s.tokenBuildRefresh(user.UUID, "uuid")
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	accessToken, err = s.TokenBuildAccess(user.UUID, "uuid")
	if err != nil {
		s.logger.Error(err.Error())
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
		return "", "", errors.New("your microsoft account does not own minecraft")
	}

	uuid = responseJson["id"].(string)
	name = responseJson["name"].(string)

	return uuid, name, nil
}
