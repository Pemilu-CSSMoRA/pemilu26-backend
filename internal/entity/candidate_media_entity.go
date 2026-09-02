package entity

import "github.com/google/uuid"

type CandidateMedia struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CandidateID uuid.UUID `json:"candidate_id"`
	FileURL     string    `json:"file_url"`

	Candidate *Candidate `gorm:"foreignKey:CandidateID" json:"candidate"`

	Timestamp
}

func (CandidateMedia) TableName() string {
	return "candidate_media"
}
