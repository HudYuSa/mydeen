package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Scheluded Status = "scheduled"
	Live      Status = "live"
	Finished  Status = "finished"
)

type QuestionLength int

const (
	Short    QuestionLength = 160
	Medium   QuestionLength = 240
	Long     QuestionLength = 360
	VeryLong QuestionLength = 480
)

type MaxQuestions int

const (
	LowCount  MaxQuestions = 1
	MidCount  MaxQuestions = 3
	HighCount MaxQuestions = 5
)

type Event struct {
	EventID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	AdminID           uuid.UUID      `gorm:"not null"`
	EventName         string         `gorm:"not null"`
	Status            Status         `gorm:"not null"`
	Moderation        bool           `gorm:"not null"`
	MaxQuestions      MaxQuestions   `gorm:"not null"`
	MaxQuestionLength QuestionLength `gorm:"not null"`
	EventCode         string         `gorm:"not null"`
	StartDate         time.Time      `gorm:"not null"`
	CreatedAt         time.Time      `gorm:"not null"`
	UpdatedAt         time.Time      `gorm:"not null"`
	Admin             Admin          `gorm:"foreignKey:AdminID;references:AdminID"`
}
