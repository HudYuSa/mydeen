package controllers

import (
	"gorm.io/gorm"
)

type WebsocketController interface {
}

type websocketController struct {
	DB *gorm.DB
}

func NewWebsocketController(db *gorm.DB) WebsocketController {
	return &websocketController{
		DB: db,
	}
}
