package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"

	"goa-golang/config"
	"goa-golang/models"
	"goa-golang/requests"
	"goa-golang/services"
	"goa-golang/viewModels"
)

type AuthController struct {
	AuthService services.IAuthService
}

// Login godoc
// @Summary
// @Description
// @Tags Auth
// @Accept  json
// @Accept  multipart/form-data
// @Accept  application/x-www-form-urlencoded
// @Produce json
// @Param email body string true "<code>required</code>  <code>min:4</code> <code>max:50</code> <code>must be email</code>" minlength(4) maxlength(50)
// @Param password body string true "<code>required</code>  <code>min:8</code> <code>max:50</code>" minlength(8) maxlength(50)
// @Param platform body string true "<code>required</code>  <code>In('panel', 'web', 'mobile')/code>"
// @Success 200 {object} viewModels.HTTPSuccessResponse{data=viewModels.Login}
// @Failure 422 {object} viewModels.HTTPErrorResponse{}
// @Failure 400 {object} viewModels.Message{}
// @Failure 500 {object} viewModels.Message{}
// @Router /v1/login [post]
func (a AuthController) Login(c *gin.Context) (err error) {

	// Request Bind And Validation
	var request requests.LoginRequest

	if err := c.ShouldBindWith(&request, binding.Query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var user models.User
	user, err = a.AuthService.GetUserByEmail(request.Body.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusUnprocessableEntity, viewModels.ValidationResponse(map[string]string{
				"email": "email or password is incorrect",
			}))
		} else {
			return echo.ErrInternalServerError
		}
	}

	var verify bool
	verify, err = a.AuthService.Check(request.Body.Email, request.Body.Password)
	if !verify {
		return c.JSON(http.StatusUnprocessableEntity, viewModels.ValidationResponse(map[string]string{
			"email": "email or password is incorrect",
		}))
	}

	accessTokenExp := time.Now().Add(time.Hour * 720).Unix()

	claims := &config.JwtCustomClaims{
		AuthID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var accessToken string
	accessToken, err = token.SignedString([]byte(config.Conf.SecretKey))
	if err != nil {
		return
	}

	// Response
	return c.JSON(http.StatusOK, viewModels.SuccessResponse(viewModels.Login{
		AccessToken:    accessToken,
		AccessTokenExp: accessTokenExp,
		User:           user,
	}))
}
