package entity

import "github.com/google/uuid"

type CandidateMission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CandidateID uuid.UUID `gorm:"type:uuid;not null" json:"candidate_id"`
	Position    int64     `gorm:"not null" json:"position"`
	Content     string    `gorm:"type: text;not null" json:"content"`

	Candidate *Candidate `gorm:"foreignKey:CandidateID" json:"candidate"`
	Timestamp
}

func (CandidateMission) TableName() string {
	return "candidate_missions"
}
