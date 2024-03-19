package controllers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/HudYuSa/mydeen/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	OtpCheck(ctx *gin.Context)
	UpdateUsername(ctx *gin.Context)
	UpdateEmail(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	LogOut(ctx *gin.Context)
	Profile(ctx *gin.Context)
}

type adminController struct {
	DB *gorm.DB
}

func NewAdminController(db *gorm.DB) AdminController {
	return &adminController{
		DB: db,
	}
}

func (ac *adminController) SignUp(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	var payload dtos.AdminSignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(payload)

	// hash the user password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// find invitation in database using payload.invitationCode
	invitation := models.Invitation{}
	invitationResult := ac.DB.WithContext(dbTimeoutCtx).Where("code = ? AND used = ?", payload.InvitationCode, false).First(&invitation)
	// if there's an error when fetching from db
	if invitationResult.Error == gorm.ErrRecordNotFound {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Cannot found your invitation code")
		return
	} else if invitationResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, invitationResult.Error.Error())
		return
	}

	// check for expiry
	fmt.Println(invitation)
	fmt.Println(invitation.ExpireDate)
	if time.Now().UTC().After(invitation.ExpireDate) {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Invitation expired")
		return
	}

	// create admin entity
	now := time.Now()
	randomCode, err := utils.GenerateRandomNumCodeLength(12)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
	}

	fmt.Println(randomCode)

	newAdmin := models.Admin{
		InvitationID: invitation.InvitationID,
		Username:     payload.Username,
		Email:        strings.ToLower(payload.Email),
		Password:     hashedPassword,
		AdminCode:    randomCode,
		Enable2fa:    false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// save new admin to database
	adminResult := ac.DB.WithContext(dbTimeoutCtx).Create(&newAdmin)
	// handle any possible error
	if adminResult.Error != nil && strings.Contains(adminResult.Error.Error(), "duplicate key value violates unique") {
		dtos.RespondWithError(ctx, http.StatusConflict, "Admin with that email already exist")
		return
	} else if adminResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, adminResult.Error.Error())
		return
	}

	// update the otp as used in database
	invitationUpdateResult := ac.DB.WithContext(dbTimeoutCtx).Model(&models.Invitation{}).Where("invitation_id = ?", invitation.InvitationID).Update("used", true)
	if invitationUpdateResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, invitationUpdateResult.Error.Error())
		return
	}

	// send the response
	dtos.RespondWithJson(ctx, http.StatusCreated, gin.H{
		"Message": "Successfully created admin",
	})
}

func (ac *adminController) SignIn(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	// try to bind the request body to the payload struct
	var payload dtos.AdminSignInInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(payload, "admin sign in payload")
	// find the admin in database
	admin := models.Admin{}
	adminResult := ac.DB.WithContext(dbTimeoutCtx).Where("email = ?", strings.ToLower(payload.Email)).First(&admin)
	// if there's an error when fetching from db
	if adminResult.Error == gorm.ErrRecordNotFound {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Invalid Email or password")
		return
	} else if adminResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, adminResult.Error.Error())
		return
	}

	// verify admin password
	if err := utils.VerifyPassword(admin.Password, payload.Password); err != nil {
		dtos.RespondWithError(ctx, http.StatusForbidden, "invalid email or password")
		return
	}

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.GlobalConfig.AccessTokenExpiresIn, dtos.GenerateAdminResponse(&admin), config.GlobalConfig.AccessTokenPrivateKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := utils.CreateToken(config.GlobalConfig.RefreshTokenExpiresIn, dtos.GenerateAdminResponse(&admin), config.GlobalConfig.RefreshTokenPrivateKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Get the client's request host
	host := ctx.Request.Host

	// Extract the domain from the request host
	parts := strings.Split(host, ":")
	domain := parts[0]

	// set accesstoken and refresh token to client cookie
	// max age time 60 so it become minute
	// set samesite to none
	fmt.Println(domain)
	ctx.SetCookie("access_token", accessToken, config.GlobalConfig.AccessTokenMaxAge*60, "/", domain, true, true)
	ctx.SetCookie("refresh_token", refreshToken, config.GlobalConfig.RefreshTokenMaxAge*60, "/", domain, true, true)

	// send response
	dtos.RespondWithJson(ctx, http.StatusOK, dtos.GenerateAdminResponse(&admin))
}

func (ac *adminController) OtpCheck(ctx *gin.Context) {
	dtos.RespondWithError(ctx, http.StatusNotFound, "api not implmented")
}

func (ac *adminController) UpdateUsername(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	var payload dtos.UpdateAdminUsername

	// try to bind the reauest body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(payload)

	UpdateUsernameResult := ac.DB.WithContext(dbTimeoutCtx).Model(&currentAdmin).Where("admin_id = ?", currentAdmin.AdminID).Update("username", payload.Username)
	if UpdateUsernameResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateUsernameResult.Error.Error())
		return
	}

	dtos.RespondWithJson(ctx, http.StatusOK, dtos.GenerateAdminResponse(&currentAdmin))
}

func (ac *adminController) UpdateEmail(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	var payload dtos.UpdateAdminEmail

	// try to bind the reauest body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(payload)

	UpdateUsernameResult := ac.DB.WithContext(dbTimeoutCtx).Model(&currentAdmin).Where("admin_id = ?", currentAdmin.AdminID).Update("email", payload.Email)
	if UpdateUsernameResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateUsernameResult.Error.Error())
		return
	}

	dtos.RespondWithJson(ctx, http.StatusOK, dtos.GenerateAdminResponse(&currentAdmin))
}

func (ac *adminController) UpdatePassword(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ac *adminController) RefreshAccessToken(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	// if there's no token from header or cookie
	refreshToken := utils.GetToken(ctx, "refresh_token", "x-refresh-token")
	if reflect.ValueOf(refreshToken).IsZero() {
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "you're not allowed to access this endpoint")
		return
	}

	// validate the token
	sub, err := utils.ValidateToken(refreshToken, config.GlobalConfig.RefreshTokenPublicKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	// find the user that has the refresh token
	var admin models.Admin
	result := ac.DB.WithContext(dbTimeoutCtx).Where("admin_id = ?", sub["admin_id"]).First(&admin)
	fmt.Println(sub["admin_id"])
	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "the user belonging to this token doesn't exist anymore")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, result.Error.Error())
		}
		return
	}

	// reissue new accesstoken
	accessToken, err := utils.CreateToken(config.GlobalConfig.AccessTokenExpiresIn, dtos.GenerateAdminResponse(&admin), config.GlobalConfig.AccessTokenPrivateKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Get the client's request host
	host := ctx.Request.Host

	// Extract the domain from the request host
	parts := strings.Split(host, ":")
	domain := parts[0]

	// set accesstoken and refresh token to client cookie
	// max age time 60 so it become minute
	// set samesite to none
	fmt.Println(domain)

	// set new accesstoken cookie
	ctx.SetCookie("access_token", accessToken, config.GlobalConfig.AccessTokenMaxAge*60, "/", domain, true, true)

	dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (ac *adminController) LogOut(ctx *gin.Context) {
	// Get the client's request host
	host := ctx.Request.Host

	// Extract the domain from the request host
	parts := strings.Split(host, ":")
	domain := parts[0]

	// set accesstoken and refresh token to client cookie
	// max age time 60 so it become minute
	// set samesite to none
	fmt.Println(domain)
	ctx.SetCookie("access_token", "", -1, "/", domain, true, true)
	ctx.SetCookie("refresh_token", "", -1, "/", domain, true, true)

	dtos.RespondWithJson(ctx, http.StatusOK, "successfully logout user")
}

func (ac *adminController) Profile(ctx *gin.Context) {
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)
	dtos.RespondWithJson(ctx, http.StatusOK, dtos.GenerateAdminResponse(&currentAdmin))
}
