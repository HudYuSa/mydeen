package controllers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/HudYuSa/mydeen/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommonController interface {
	CheckRole(ctx *gin.Context)
}

type commonController struct {
	DB *gorm.DB
}

func NewCommonController(db *gorm.DB) CommonController {
	return &commonController{
		DB: db,
	}
}

func (cc *commonController) CheckRole(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	// sent user data everytime
	user := ctx.MustGet("user").(dtos.User)

	refreshToken := utils.GetToken(ctx, "refresh_token", "x-refresh-token")
	// if there's no token from header or cookie and
	if reflect.ValueOf(refreshToken).IsZero() {
		// dtos.RespondWithError(ctx, http.StatusUnauthorized, "you're not allowed to access this endpoint")
		dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
			"user":    dtos.GenerateUserResponse(&user),
			"account": false,
			"error":   true,
			"message": "you're not allowed to access this endpoint",
		})
		return
	}

	// validate the token
	sub, err := utils.ValidateToken(refreshToken, config.GlobalConfig.RefreshTokenPublicKey)
	if err != nil {
		// dtos.RespondWithError(ctx, http.StatusUnauthorized, err.Error())
		dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
			"user":    dtos.GenerateUserResponse(&user),
			"account": false,
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	userType, ok := sub["type"]
	fmt.Println(userType, ok)
	if ok {
		// check for 2 different user type
		// fmt.Println(userType)
		if userType == "master" {
			// find the master
			master := models.Master{}

			result := cc.DB.WithContext(dbTimeoutCtx).Where("master_id = ?", sub["master_id"]).First(&master)
			if result.Error != nil {
				// dtos.RespondWithError(ctx, http.StatusNotFound, "the user belonging to this token doesn't exist anymore")
				dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
					"user":    dtos.GenerateUserResponse(&user),
					"account": false,
					"error":   true,
					"message": result.Error.Error(),
				})
				return
			}

			dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
				"user":    dtos.GenerateUserResponse(&user),
				"account": dtos.GenerateMasterResponse(&master),
			})

		} else if userType == "admin" {
			// find the admin
			admin := models.Admin{}

			result := cc.DB.WithContext(dbTimeoutCtx).Where("admin_id = ?", sub["admin_id"]).First(&admin)
			if result.Error != nil {
				// dtos.RespondWithError(ctx, http.StatusNotFound, "the user belonging to this token doesn't exist anymore")
				dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
					"user":    dtos.GenerateUserResponse(&user),
					"account": false,
					"error":   true,
					"message": result.Error.Error(),
				})
				return
			}

			dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
				"user":    dtos.GenerateUserResponse(&user),
				"account": dtos.GenerateAdminResponse(&admin),
			})
		} else {
			dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
				"user":    dtos.GenerateUserResponse(&user),
				"account": false,
				"error":   true,
				"message": "you're not allowed to access this endpoint",
			})
		}
	} else {
		dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
			"user":    dtos.GenerateUserResponse(&user),
			"account": false,
			"error":   true,
			"message": "you're not allowed to access this endpoint",
		})
	}

}
