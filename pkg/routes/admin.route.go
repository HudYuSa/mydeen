package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/HudYuSa/mydeen/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type AdminRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type adminRoutes struct {
	AdminController controllers.AdminController
}

func NewAdminRoutes(adminController controllers.AdminController) AdminRoutes {
	return &adminRoutes{
		AdminController: adminController,
	}
}

func (ar *adminRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/admin")

	router.POST("/signup", ar.AdminController.SignUp)
	router.POST("/signin", ar.AdminController.SignIn)
	router.POST("/otp", ar.AdminController.OtpCheck)
	router.GET("/refresh", ar.AdminController.RefreshAccessToken)
	router.GET("/logout", ar.AdminController.LogOut)

	router.Use(middlewares.AuthenticateAdmin())
	router.PATCH("/edit/username", ar.AdminController.UpdateUsername)
	router.PATCH("/edit/email", ar.AdminController.UpdateEmail)
}
