package GMiddleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goa-golang/services"
)

type Auth struct {
	UserService services.IUserService
}

func (s Auth) AuthMiddleware(c *gin.Context) {
	name := c.Param("name")

	if name != "lix3nn" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
	}

	c.Next()
}
