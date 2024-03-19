package middlewares

import (
	"log"
	"net/http"

	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/HudYuSa/mydeen/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get user cookie
		userCookie, err := utils.GetCookie(ctx, "user")
		if err != nil {
			// if there's no cookie user found
			// then set a new user cookie for the client
			if err == http.ErrNoCookie {
				user := dtos.User{
					ID: uuid.New(),
				}

				utils.SetNewUserCookie(ctx, &user)
			} else {
				dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
				return
			}
		}

		// set user inside context
		var user dtos.User
		if userCookie == "" {
			user = dtos.User{
				ID: uuid.New(),
			}
			utils.SetNewUserCookie(ctx, &user)

			ctx.Set("user", user)
		} else {
			decodeErr := dtos.DecodeJson([]byte(userCookie), &user)
			if decodeErr != nil {
				log.Println(decodeErr.Error())
			}
			ctx.Set("user", user)
		}
	}
}
