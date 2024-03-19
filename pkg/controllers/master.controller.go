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
	"github.com/HudYuSa/mydeen/pkg/services"
	"github.com/HudYuSa/mydeen/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MasterController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	OtpCheck(ctx *gin.Context)
	ReissueVerificationCode(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	LogOut(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	Profile(ctx *gin.Context)
	GenerateInvitationCode(ctx *gin.Context)
}

type masterController struct {
	DB *gorm.DB
}

func NewMasterController(db *gorm.DB) MasterController {
	return &masterController{DB: db}
}

func (mc *masterController) SignUp(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	var payload dtos.MasterSignUpInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// hash the user password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// check duplicate master account that haven't been verified
	master := models.Master{}
	masterGetResult := mc.DB.WithContext(dbTimeoutCtx).Where("email = ?", strings.ToLower(payload.Email)).Find(&master)
	// check for any possible error
	if masterGetResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, masterGetResult.Error.Error())
		return
	}

	if !master.Verified && !reflect.ValueOf(master.MasterID).IsZero() {
		// verify master password
		if err := utils.VerifyPassword(master.Password, payload.Password); err != nil {
			dtos.RespondWithError(ctx, http.StatusForbidden, "invalid email or password")
			return
		}

		dtos.RespondWithError(ctx, http.StatusForbidden, "Please verify your email address first")
		return
	}

	// create master entity
	now := time.Now()
	newMaster := models.Master{
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Verified:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save the data to database using gorm
	masterResult := mc.DB.WithContext(dbTimeoutCtx).Create(&newMaster)
	// check for any possible error
	if masterResult.Error != nil && strings.Contains(masterResult.Error.Error(), "duplicate key value violates unique") {
		dtos.RespondWithError(ctx, http.StatusConflict, "Master with that email already exist")
		return
	} else if masterResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, masterResult.Error.Error())
		return
	}

	// generate verification code
	code, err := utils.GenerateVerificationToken()
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, "error generating verification code")
		return
	}
	fmt.Println(code)

	// create verification code entity to store in the database
	verificationCode := models.VerificationCode{
		MasterID:   newMaster.MasterID,
		Code:       code,
		ExpireDate: time.Now().UTC().Add(24 * 7 * time.Hour),
		Used:       false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// store the verification code in database
	vcResult := mc.DB.WithContext(dbTimeoutCtx).Create(&verificationCode)
	// check for any possible error
	if vcResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, vcResult.Error.Error())
		return
	}

	// send email with a link that contains a verification link to verify their email address
	emailErr := services.SendVerificationCodeGomail(verificationCode.Code, []string{newMaster.Email})
	if emailErr != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, emailErr.Error())
		return
	}

	// send the response
	dtos.RespondWithJson(ctx, http.StatusCreated, gin.H{
		"Message": "Successfully created master, check your email to verify your account",
	})
}

func (mc *masterController) SignIn(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	var payload dtos.MasterSignInInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(payload)
	// find the master in database
	master := models.Master{}
	result := mc.DB.WithContext(dbTimeoutCtx).Where("email = ?", strings.ToLower(payload.Email)).First(&master)
	// if there's an error when fetching from db
	if result.Error == gorm.ErrRecordNotFound {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Invalid Email or password")
		return
	} else if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, result.Error.Error())
		return
	}

	fmt.Println("ketemu master tanpa error")
	// verify master password
	if err := utils.VerifyPassword(master.Password, payload.Password); err != nil {
		dtos.RespondWithError(ctx, http.StatusForbidden, "invalid email or password")
		return
	}

	// if the master is not verified sent forbidden status
	if !master.Verified {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Please verify your email address first")
		return
	}

	// generate new otp code
	now := time.Now()
	otp := models.MasterOtp{
		Code:       utils.GenerateRandomCode(),
		MasterID:   master.MasterID,
		ExpireDate: time.Now().UTC().Add(10 * time.Minute),
		Used:       false,
		CreatedAt:  now,
	}

	// save the otp code to database
	otpResult := mc.DB.WithContext(dbTimeoutCtx).Create(&otp)
	if otpResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, otpResult.Error.Error())
		return
	}

	// send email with the otp code
	emailErr := services.SendOtpCode(otp.Code, []string{master.Email})
	if emailErr != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, emailErr.Error())
		return
	}

	dtos.RespondWithJson(ctx, http.StatusAccepted, gin.H{
		"otp_id": otp.MasterOtpID,
	})
}

func (mc *masterController) OtpCheck(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	var payload dtos.MasterOtpInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// get the otp from database
	otp := models.MasterOtp{}
	otpResult := mc.DB.WithContext(dbTimeoutCtx).Where("code = ? AND used = ?", payload.Code, false).First(&otp)
	if otpResult.Error == gorm.ErrRecordNotFound {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Invalid otp code")
		return
	} else if otpResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, otpResult.Error.Error())
		return
	}

	// check for expiry
	fmt.Println(otp.ExpireDate)
	fmt.Println(time.Now().UTC())
	if time.Now().UTC().After(otp.ExpireDate) {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Otp expired")
		return
	}

	// find master of the otp
	master := models.Master{}
	masterResult := mc.DB.WithContext(dbTimeoutCtx).Where("master_id = ?", otp.MasterID).First(&master)
	if masterResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, masterResult.Error.Error())
		return
	}

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.GlobalConfig.AccessTokenExpiresIn, dtos.GenerateMasterResponse(&master), config.GlobalConfig.AccessTokenPrivateKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := utils.CreateToken(config.GlobalConfig.RefreshTokenExpiresIn, dtos.GenerateMasterResponse(&master), config.GlobalConfig.RefreshTokenPrivateKey)
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
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", accessToken, config.GlobalConfig.AccessTokenMaxAge*60, "/", domain, true, true)
	ctx.SetCookie("refresh_token", refreshToken, config.GlobalConfig.RefreshTokenMaxAge*60, "/", domain, true, true)

	// update the otp as used in database
	otpUpdateResult := mc.DB.WithContext(dbTimeoutCtx).Model(&models.MasterOtp{}).Where("master_otp_id = ?", otp.MasterOtpID).Update("used", true)
	if otpUpdateResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, otpUpdateResult.Error.Error())
		return
	}

	dtos.RespondWithJson(ctx, http.StatusOK, dtos.GenerateMasterResponse(&master))
}

func (mc *masterController) ReissueVerificationCode(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	var payload dtos.MasterSignUpInput

	// try to bind the request body tot he payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// find the master in database
	master := models.Master{}
	masterResult := mc.DB.WithContext(dbTimeoutCtx).Where("email = ?", payload.Email).First(&master)
	// check for error
	if masterResult.Error == gorm.ErrRecordNotFound {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Invalid Email or password")
		return
	} else if masterResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, masterResult.Error.Error())
		return
	}

	// verify master password
	if err := utils.VerifyPassword(master.Password, payload.Password); err != nil {
		dtos.RespondWithError(ctx, http.StatusForbidden, "invalid email or password")
		return
	}

	// generate new verification code
	newVerificationCode, err := utils.GenerateVerificationToken()
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, "error generating verification code")
		return
	}
	fmt.Println(newVerificationCode)

	// Update verification code of the master
	vcResult := mc.DB.WithContext(dbTimeoutCtx).Model(&models.VerificationCode{}).Where("master_id = ?", master.MasterID).Update("code", newVerificationCode)
	if vcResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, vcResult.Error.Error())
		return
	}

	// send a new email with a link that contains a verification link to verify their email address
	emailErr := services.SendVerificationCodeGomail(newVerificationCode, []string{master.Email})
	if emailErr != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, emailErr.Error())
		return
	}

	// send the response
	dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
		"Message": "Successfully update verification code, check your email to verify your account",
	})
}

func (mc *masterController) RefreshAccessToken(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	refreshToken := utils.GetToken(ctx, "refresh_token", "x-refresh-token")

	// if there's no token from header or cookie
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
	var master models.Master
	result := mc.DB.WithContext(dbTimeoutCtx).Where("master_id = ?", sub["master_id"]).First(&master)
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
	accessToken, err := utils.CreateToken(config.GlobalConfig.AccessTokenExpiresIn, dtos.GenerateMasterResponse(&master), config.GlobalConfig.AccessTokenPrivateKey)
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
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", accessToken, config.GlobalConfig.AccessTokenMaxAge*60, "/", "localhost", true, true)

	dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (mc *masterController) LogOut(ctx *gin.Context) {
	// Get the client's request host
	host := ctx.Request.Host

	// Extract the domain from the request host
	parts := strings.Split(host, ":")
	domain := parts[0]

	// set accesstoken and refresh token to client cookie
	// max age time 60 so it become minute
	// set samesite to none
	fmt.Println(domain)
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("access_token", "", -1, "/", domain, true, true)
	ctx.SetCookie("refresh_token", "", -1, "/", domain, true, true)

	dtos.RespondWithJson(ctx, http.StatusOK, "successfully logout user")
}

func (mc *masterController) VerifyEmail(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	code := ctx.Query("verification_code")

	fmt.Println(code)

	if code == "" {
		dtos.RespondWithError(ctx, http.StatusBadRequest, "Invalid verification code")
		return
	}

	// Verify the verification code
	verificationCode := models.VerificationCode{}
	vcResult := mc.DB.WithContext(dbTimeoutCtx).Where("code = ? AND used = ?", code, false).First(&verificationCode)
	if vcResult.Error == gorm.ErrRecordNotFound {
		dtos.RespondWithError(ctx, http.StatusNotFound, "Invalid verification code")
		return
	} else if vcResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, vcResult.Error.Error())
		return
	}

	// check for expiry
	if time.Now().UTC().After(verificationCode.ExpireDate) {
		dtos.RespondWithError(ctx, http.StatusForbidden, "Verification code expired")
		return
	}

	// update the master data on database
	masterUpdateResult := mc.DB.WithContext(dbTimeoutCtx).Model(&models.Master{}).Where("master_id = ?", verificationCode.MasterID).Update("verified", true)
	if masterUpdateResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, masterUpdateResult.Error.Error())
		return
	}

	vcUpdateResult := mc.DB.WithContext(dbTimeoutCtx).Model(&models.VerificationCode{}).Where("verification_code_id = ?", verificationCode.VerificationCodeID).Update("used", true)
	if vcUpdateResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, vcUpdateResult.Error.Error())
		return
	}

	dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
		"Message": "Email verified successfully",
	})
}

func (mc *masterController) Profile(ctx *gin.Context) {
	currentMaster := ctx.MustGet("currentMaster").(models.Master)
	dtos.RespondWithJson(ctx, http.StatusOK, dtos.GenerateMasterResponse(&currentMaster))
}

func (mc *masterController) GenerateInvitationCode(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	currentMaster := ctx.MustGet("currentMaster").(models.Master)
	fmt.Println(currentMaster)
	// create new invitation entity
	now := time.Now()
	newInvitation := models.Invitation{
		MasterID:   currentMaster.MasterID,
		Code:       utils.GenerateRandomCodeLength(12),
		ExpireDate: time.Now().UTC().Add(24 * 30 * time.Hour),
		Used:       false,
		CreatedAt:  now,
	}

	invitationResult := mc.DB.WithContext(dbTimeoutCtx).Create(&newInvitation)
	if invitationResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, invitationResult.Error.Error())
		return
	}

	dtos.RespondWithJson(ctx, http.StatusCreated, dtos.GenerateInvitationResponse(&newInvitation))

}
