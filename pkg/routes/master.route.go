package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/HudYuSa/mydeen/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type MasterRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type masterRoutes struct {
	MasterController controllers.MasterController
}

func NewMasterRoutes(masterController controllers.MasterController) MasterRoutes {
	return &masterRoutes{
		MasterController: masterController,
	}
}

func (mr *masterRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/master")

	router.POST("/signup", mr.MasterController.SignUp)
	router.POST("/signin", mr.MasterController.SignIn)
	router.POST("/otp", mr.MasterController.OtpCheck)
	router.POST("/reissue_verification_code", mr.MasterController.ReissueVerificationCode)
	router.GET("/refresh", mr.MasterController.RefreshAccessToken)
	router.GET("/verify", mr.MasterController.VerifyEmail)
	router.GET("logout", mr.MasterController.LogOut)

	// protected routes
	router.Use(middlewares.AuthenticateMaster())
	router.GET("/profile", mr.MasterController.Profile)
	router.GET("/generate_invitation", mr.MasterController.GenerateInvitationCode)
}
