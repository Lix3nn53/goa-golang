package authController

import (
	"errors"
	appError "goa-golang/app/error"
	"goa-golang/app/service/authService"
	"goa-golang/internal/logger"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//UserControllerInterface define the user controller interface methods
type AuthControllerInterface interface {
	GoogleOauth2(c *gin.Context)
	TwitchOauth2(c *gin.Context)
	DiscordOauth2(c *gin.Context)
	MinecraftOauth2(c *gin.Context)
	RefreshAccessToken(c *gin.Context)
	Logout(c *gin.Context)
	AuthMiddleware() gin.HandlerFunc
}

// UserController handles communication with the user service
type AuthController struct {
	service authService.AuthServiceInterface
	logger  logger.Logger
}

// NewUserController implements the user controller interface.
func NewAuthController(service authService.AuthServiceInterface, logger logger.Logger) AuthControllerInterface {
	return &AuthController{
		service,
		logger,
	}
}

type AuthResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type RefreshAccessResponse struct {
	AccessToken string `json:"access_token"`
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) GoogleOauth2(c *gin.Context) {
	code := c.Query("code")

	refreshToken, accessToken, err := uc.service.GoogleOauth2(code)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusBadRequest, err)
		return
	}

	response := AuthResponse{RefreshToken: refreshToken, AccessToken: accessToken}
	c.JSON(http.StatusOK, response)
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) TwitchOauth2(c *gin.Context) {
	code := c.Query("code")

	refreshToken, accessToken, err := uc.service.TwitchOauth2(code)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusBadRequest, err)
		return
	}

	response := AuthResponse{RefreshToken: refreshToken, AccessToken: accessToken}
	c.JSON(http.StatusOK, response)
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) DiscordOauth2(c *gin.Context) {
	code := c.Query("code")

	refreshToken, accessToken, err := uc.service.DiscordOauth2(code)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusBadRequest, err)
		return
	}

	response := AuthResponse{RefreshToken: refreshToken, AccessToken: accessToken}
	c.JSON(http.StatusOK, response)
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) MinecraftOauth2(c *gin.Context) {
	code := c.Query("code")

	refreshToken, accessToken, err := uc.service.MinecraftOauth2(code)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusBadRequest, err)
		return
	}

	response := AuthResponse{RefreshToken: refreshToken, AccessToken: accessToken}
	c.JSON(http.StatusOK, response)
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) RefreshAccessToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		appError.Respond(c, http.StatusBadRequest, errors.New("no authorization header provided"))
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		appError.Respond(c, http.StatusBadRequest, errors.New("could not find bearer token in authorization header"))
		return
	}

	id, idField, err := uc.service.TokenValidateRefresh(token)
	if err != nil {
		appError.Respond(c, http.StatusBadRequest, err)
		return
	}

	accessToken, err := uc.service.TokenBuildAccess(id, idField)
	if err != nil {
		appError.Respond(c, http.StatusInternalServerError, err)
		return
	}

	response := RefreshAccessResponse{AccessToken: accessToken}
	c.JSON(http.StatusOK, response)
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) Logout(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		appError.Respond(c, http.StatusBadRequest, errors.New("no authorization header provided"))
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		appError.Respond(c, http.StatusBadRequest, errors.New("could not find bearer token in authorization header"))
		return
	}

	id, idField, err := uc.service.TokenValidateRefresh(token)
	if err != nil {
		appError.Respond(c, http.StatusUnauthorized, err)
		return
	}

	err = uc.service.Logout(id, idField, token)
	if err != nil {
		appError.Respond(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (uc *AuthController) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			appError.Respond(c, http.StatusBadRequest, errors.New("no authorization header provided"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			appError.Respond(c, http.StatusBadRequest, errors.New("could not find bearer token in authorization header"))
			c.Abort()
			return
		}

		id, idField, err := uc.service.TokenValidate(token, os.Getenv("ACCESS_SECRET"))
		if err != nil {
			appError.Respond(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		c.Set("tokenID", id)
		c.Set("tokenIDField", idField)
		c.Next()
		// after request
	}
}
