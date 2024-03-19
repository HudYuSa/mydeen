package middlewares

import (
	"net/http"
	"reflect"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/internal/connection"
	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/HudYuSa/mydeen/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthenticateMaster() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := utils.GetToken(ctx, "access_token", "Authorization")

		// if there's no token from header or cookie
		if reflect.ValueOf(accessToken).IsZero() {
			dtos.RespondWithError(ctx, http.StatusUnauthorized, "you're not allowed to access this endpoint")
			return
		}

		// validate the token and get the user from the sub/subject
		account, err := utils.ValidateToken(accessToken, config.GlobalConfig.AccessTokenPublicKey)
		if err != nil {
			dtos.RespondWithError(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		masterId, masterOk := account["master_id"]
		if masterOk {
			// find the master that belongs to the token
			var master models.Master
			masterResult := connection.DB.First(&master, "master_id = ?", masterId)
			if masterResult.Error != nil {
				dtos.RespondWithError(ctx, http.StatusInternalServerError, "you're not allowed to access this endpoint")
				return
			}

			ctx.Set("currentMaster", master)
			ctx.Next()
			return
		} else {
			dtos.RespondWithError(ctx, http.StatusUnauthorized, "you're not allowed to access this endpoint")
		}
	}
}
