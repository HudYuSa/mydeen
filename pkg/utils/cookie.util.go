package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func SetNewUserCookie(ctx *gin.Context, user *dtos.User) {

	userJson := dtos.EncodeJson(user)

	// Get the client's request host
	host := ctx.Request.Host

	// Extract the domain from the request host
	parts := strings.Split(host, ":")
	domain := parts[0]

	// set accesstoken and refresh token to client cookie
	// max age time 60 so it become minute
	// set samesite to none
	fmt.Println(domain)
	ctx.SetCookie("user", string(userJson), int(200*365*24*time.Hour), "/", domain, true, true)
}
