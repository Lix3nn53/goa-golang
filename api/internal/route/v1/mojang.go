package v1

import (
	"goa-golang/app/controller/mojangController"

	"github.com/gin-gonic/gin"
)

func SetupMojangRoute(mojang *gin.RouterGroup, c mojangController.MojangControllerInterface) *gin.RouterGroup {
	mojang.POST("/profiles", c.Profiles)

	return mojang
}
