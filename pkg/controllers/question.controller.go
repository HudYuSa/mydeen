package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"gorm.io/gorm"
)

type QuestionController interface {
	// http
	GetEventQuestions(ctx *gin.Context)
	GetUserTotalQuestions(ctx *gin.Context)
	// websocket
	CreateQuestion(s *melody.Session, b []byte)
	DeleteQuestion(s *melody.Session, b []byte)
	EditQuestion(s *melody.Session, b []byte)
}

type questionController struct {
	DB     *gorm.DB
	Melody *melody.Melody
}

func NewQuestionController(db *gorm.DB, melody *melody.Melody) QuestionController {
	return &questionController{
		DB:     db,
		Melody: melody,
	}
}

// http
func (qc *questionController) GetEventQuestions(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	user := ctx.MustGet("user").(dtos.User)

	eventId := ctx.Param("event_id")

	questions := []models.Question{}
	questionsResult := qc.DB.WithContext(dbTimeoutCtx).Preload("Likes").Where("event_id", eventId).Find(&questions)
	if questionsResult.Error != nil {
		switch questionsResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given code")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, questionsResult.Error.Error())
		}
		return
	}

	questionsResponse := []dtos.QuestionResponse{}
	for _, question := range questions {
		questionsResponse = append(questionsResponse, *dtos.GenerateQuestionResponse(&question, user))
	}

	dtos.RespondWithJson(ctx, http.StatusOK, questionsResponse)
}

func (qc *questionController) GetUserTotalQuestions(ctx *gin.Context) {
	// get user identifier
	// get event id
	// get questions data/total using the event id and user id
	// return total questions
}

// websocket
func (qc *questionController) CreateQuestion(s *melody.Session, b []byte) {
	// dbtimeoutctx for websocket
	dbTimeoutCtx, cancel := context.WithTimeout(s.Request.Context(), time.Duration(config.GlobalConfig.DatabaseTimeout)*time.Millisecond)
	defer cancel()

	user := s.Request.Context().Value("user").(dtos.User)
	log.Println(user)

	var payload dtos.CreateQuestionInput

	if err := json.Unmarshal(b, &payload); err != nil {
		s.Write(dtos.WebSocketRespondError(dtos.Question, err.Error()))
		return
	}

	// create new question instance
	log.Println("payload: ", payload)
	now := time.Now().UTC()

	eventId, err := uuid.Parse(payload.EventID)
	if err != nil {
		s.Write(dtos.WebSocketRespondError(dtos.Question, "no event with the given id"))
	}

	newQuestion := models.Question{
		EventID:   eventId,
		UserID:    user.ID,
		Username:  payload.Username,
		Content:   payload.Content,
		Starred:   false,
		Approved:  false,
		Answered:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	log.Println("question instance: ", newQuestion)
	// save new question to the database
	questionResult := qc.DB.WithContext(dbTimeoutCtx).Create(&newQuestion)
	if questionResult.Error != nil && strings.Contains(questionResult.Error.Error(), "duplicate key value violates unique") {
		log.Println(questionResult.Error.Error())
		s.Write(dtos.WebSocketRespondError(dtos.Question, questionResult.Error.Error()))
		return
	} else if questionResult.Error != nil {
		log.Println(questionResult.Error.Error())
		s.Write(dtos.WebSocketRespondError(dtos.Question, questionResult.Error.Error()))
		return
	}

	// respond back to the client websocket
	log.Println(newQuestion)
	qc.Melody.Broadcast(dtos.WebSocketRespondJson(dtos.Question, dtos.CreateQuestionType, dtos.GenerateQuestionResponse(&newQuestion, user)))
}

func (qc *questionController) DeleteQuestion(s *melody.Session, b []byte) {
	// dbtimeoutctx for websocket
	dbTimeoutCtx, cancel := context.WithTimeout(s.Request.Context(), time.Duration(config.GlobalConfig.DatabaseTimeout)*time.Millisecond)
	defer cancel()

	// get current user
	user := s.Request.Context().Value("user").(dtos.User)
	log.Println(user.ID, "user")

	// kasi authentication untuk user
	// cek id user yang punya pertanyaan sebelum hapus pertanyaannya

	var payload dtos.DeleteQuestionInput

	if err := json.Unmarshal(b, &payload); err != nil {
		s.Write(dtos.WebSocketRespondError(dtos.Question, err.Error()))
		return
	}

	log.Println("payload: ", payload)

	// start a transaction
	tx := qc.DB.Begin()

	// find the question
	question := models.Question{}
	questionResult := tx.WithContext(dbTimeoutCtx).Where("question_id = ?", payload.QuestionID).First(&question)
	if questionResult.Error != nil {
		tx.Rollback()
		switch questionResult.Error.Error() {
		case "record not found":
			s.Write(dtos.WebSocketRespondError(dtos.Question, "there is no question with the given id"))
		default:
			s.Write(dtos.WebSocketRespondError(dtos.Question, questionResult.Error.Error()))
		}
		return
	}

	// check if user is the admin that created the question
	if question.UserID != user.ID {
		tx.Rollback()
		s.Write(dtos.WebSocketRespondError(dtos.Question, "You're not allowed to access this endpoint"))
		return
	}

	deleteQuestionResult := tx.WithContext(dbTimeoutCtx).Delete(&models.Question{}, "question_id = ?", question.QuestionID)
	if deleteQuestionResult.Error != nil {
		log.Println(deleteQuestionResult.Error.Error())
		s.Write(dtos.WebSocketRespondError(dtos.Question, deleteQuestionResult.Error.Error()))
		return
	}

	// commit the transaction
	tx.Commit()

	// kirim question id nya biar nanti di frontend semua active connection bisa delete question itu dari storenya
	qc.Melody.Broadcast(dtos.WebSocketRespondJson(dtos.Question, dtos.DeleteQuestionType, map[string]any{
		"question_id": payload.QuestionID,
	}))
}

func (qc *questionController) EditQuestion(s *melody.Session, b []byte) {
	// dbtimeoutctx for websocket
	dbTimeoutCtx, cancel := context.WithTimeout(s.Request.Context(), time.Duration(config.GlobalConfig.DatabaseTimeout)*time.Millisecond)
	defer cancel()

	// get current user
	user := s.Request.Context().Value("user").(dtos.User)
	log.Println(user)

	var payload dtos.EditQuestionInput

	if err := json.Unmarshal(b, &payload); err != nil {
		s.Write(dtos.WebSocketRespondError(dtos.Question, err.Error()))
		return
	}

	// create new question instance
	log.Println("payload: ", payload)

	// start a transaction
	tx := qc.DB.Begin()

	// find the question
	question := models.Question{}
	questionResult := tx.WithContext(dbTimeoutCtx).Where("question_id = ?", payload.QuestionID).First(&question)
	if questionResult.Error != nil {
		tx.Rollback()
		switch questionResult.Error.Error() {
		case "record not found":
			s.Write(dtos.WebSocketRespondError(dtos.Question, "there is no question with the given id"))
		default:
			s.Write(dtos.WebSocketRespondError(dtos.Question, questionResult.Error.Error()))
		}
		return
	}

	// check if user is the admin that created the question
	if question.UserID != user.ID {
		tx.Rollback()
		s.Write(dtos.WebSocketRespondError(dtos.Question, "You're not allowed to access this endpoint"))
		return
	}

	// update question data
	question.Content = payload.Content

	UpdateQuestionResult := tx.WithContext(dbTimeoutCtx).Where("question_id = ?", payload.QuestionID).Save(&question)
	if UpdateQuestionResult.Error != nil {
		log.Println(UpdateQuestionResult.Error.Error())
		s.Write(dtos.WebSocketRespondError(dtos.Question, UpdateQuestionResult.Error.Error()))
		return
	}

	qc.Melody.Broadcast(dtos.WebSocketRespondJson(dtos.Question, dtos.EditQuestionType, map[string]string{
		"question_id": payload.QuestionID,
		"content":     payload.Content,
	}))
}
