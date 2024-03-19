package middlewares

import (
	"reflect"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/internal/connection"
	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/HudYuSa/mydeen/pkg/utils"
	"github.com/olahol/melody"
)

// this function will authenticate admin and response with false if there's an error and true if there isn't
func WSAuthenticateAdmin(s *melody.Session, group dtos.WebSocketGroup) bool {
	accessToken := utils.GetTokenWS(s, "access_token", "Authorization")

	// if there's no token from header or cookie
	if reflect.ValueOf(accessToken).IsZero() {
		s.Write(dtos.WebSocketRespondError(group, "You're not allowed to acces this endpoint"))
		return false
	}

	// validate the token and get the user from the sub/subject
	account, err := utils.ValidateToken(accessToken, config.GlobalConfig.AccessTokenPublicKey)
	if err != nil {
		s.Write(dtos.WebSocketRespondError(group, err.Error()))
		return false
	}

	adminId, adminOk := account["admin_id"]
	if adminOk {
		var admin models.Admin
		adminResult := connection.DB.First(&admin, "admin_id = ?", adminId)
		if adminResult.Error != nil {
			s.Write(dtos.WebSocketRespondError(group, "you're not allowed to access this endpoint"))
			return false
		}

		// if there's no error set the admin and then return false
		s.Set("currentAdmin", admin)
		return true
	} else {
		s.Write(dtos.WebSocketRespondError(group, "you're not allowed to access this endpoing"))
		return false
	}
}
