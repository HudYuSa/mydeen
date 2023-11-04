package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/gin-gonic/gin"
)

type CommonRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type commonRoutes struct {
	CommonController controllers.CommonController
}

func NewCommonRoutes(commonController controllers.CommonController) CommonRoutes {
	return &commonRoutes{
		CommonController: commonController,
	}
}

func (cr *commonRoutes) SetupRoutes(rg *gin.RouterGroup) {
	rg.GET("/check_role", cr.CommonController.CheckRole)

}
