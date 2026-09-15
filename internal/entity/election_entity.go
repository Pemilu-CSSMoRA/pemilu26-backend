package entity

import "github.com/google/uuid"

type ElectionStatus string

const (
	ElectionCreated ElectionStatus = "created"
	ElectionStarted ElectionStatus = "start"
	ElectionEnded   ElectionStatus = "end"
)

type Election struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"type:text;not null" json:"name"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Year        string         `gorm:"type:text;not null" json:"year"`
	Status      ElectionStatus `gorm:"type:text;not null" json:"status"`
	Timestamp
}

func (Election) TableName() string {
	return "elections"
}
