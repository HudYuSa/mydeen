package models

import (
	"time"

	"github.com/google/uuid"
)

type status string

const (
	Scheluded status = "scheduled"
	Live      status = "live"
	Finished  status = "finished"
)

type questionLength int

const (
	Short    questionLength = 160
	Medium   questionLength = 240
	Long     questionLength = 360
	VeryLong questionLength = 540
)

type maxQuestions int

const (
	LowCount  maxQuestions = 1
	MidCount  maxQuestions = 3
	HighCount maxQuestions = 5
)

type Event struct {
	EventID           uuid.UUID
	AdminID           uuid.UUID
	EventName         string
	Status            string
	Moderation        bool
	MaxQuestions      int
	MaxQuestionLength int
	EventCode         string
	StartDate         time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Admin             Admin `gorm:"foreignKey:AdminID"`
}
