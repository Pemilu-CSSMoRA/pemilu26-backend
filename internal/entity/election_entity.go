package entity

import "github.com/google/uuid"

type election_status string

const (
	ElectionCreated election_status = "created"
	ElectionStarted election_status = "start"
	ElectionEnded   election_status = "end"
)

type Election struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      election_status `json:"status" gorm:"default:created"`
	Timestamp
}

func (Election) TableName() string {
	return "elections"
}
