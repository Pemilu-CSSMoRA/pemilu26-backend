package entity

import "github.com/google/uuid"

type CandidateMission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CandidateID uuid.UUID `json:"candidate_id"`
	Position    int64     `json:"position"`
	Content     string    `json:"content"`

	Candidate *Candidate `gorm:"foreignKey:CandidateID"`
	Timestamp
}

func (CandidateMission) TableName() string {
	return "candidate_missions"
}
