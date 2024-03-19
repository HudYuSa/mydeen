package controllers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/olahol/melody"
	"gorm.io/gorm"
)

type LikeController interface {
	// http
	// websocket
	ToggleLike(s *melody.Session, b []byte)
}

type likeController struct {
	DB     *gorm.DB
	Melody *melody.Melody
}

func NewLikeController(db *gorm.DB, melody *melody.Melody) LikeController {
	return &likeController{
		DB:     db,
		Melody: melody,
	}
}

// http
// websocket
func (lc *likeController) ToggleLike(s *melody.Session, b []byte) {
	// dbtimeoutctx for websocket
	dbTimeoutCtx, cancel := context.WithTimeout(s.Request.Context(), time.Duration(config.GlobalConfig.DatabaseTimeout)*time.Millisecond)
	defer cancel()

	// get user
	user := s.Request.Context().Value("user").(dtos.User)
	log.Println(user)

	// get like
	var payload dtos.ToggleLikeInput

	if err := json.Unmarshal(b, &payload); err != nil {
		s.Write(dtos.WebSocketRespondError(dtos.Like, err.Error()))
		return
	}

	like := models.Like{}
	// check for like in database
	checkLikeResult := lc.DB.WithContext(dbTimeoutCtx).Where("question_id = ? AND user_id = ?", payload.QuestionID, user.ID).First(&like)

	if checkLikeResult.Error != nil {
		// there's no like then create a new like in the database
		if checkLikeResult.Error == gorm.ErrRecordNotFound {
			like.QuestionID = payload.QuestionID
			like.UserID = user.ID

			// save new like to database
			likeResult := lc.DB.WithContext(dbTimeoutCtx).Create(&like)
			if likeResult.Error != nil {
				log.Println(likeResult.Error.Error())
				s.Write(dtos.WebSocketRespondError(dtos.Like, likeResult.Error.Error()))
				return
			}

			// respond with new like
			lc.Melody.Broadcast(dtos.WebSocketRespondJson(dtos.Like, dtos.ToggleLikeType, dtos.GenerateLikeResponse(&like, true)))
			return
		} else {
			log.Println(checkLikeResult.Error.Error())
			s.Write(dtos.WebSocketRespondError(dtos.Like, checkLikeResult.Error.Error()))
			return
		}
	}

	// if there's like then delete it from the database
	likeResult := lc.DB.WithContext(dbTimeoutCtx).Where("question_id = ? AND user_id = ?", payload.QuestionID, user.ID).Delete(&models.Like{})
	if likeResult.Error != nil {
		log.Println(checkLikeResult.Error.Error())
		s.Write(dtos.WebSocketRespondError(dtos.Like, checkLikeResult.Error.Error()))
		return
	}
	// respond for deleting like
	lc.Melody.Broadcast(dtos.WebSocketRespondJson(dtos.Like, dtos.ToggleLikeType, dtos.GenerateLikeResponse(&like, false)))
}
