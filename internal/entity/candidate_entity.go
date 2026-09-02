package entity

import "github.com/google/uuid"

type Candidate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ElectionID  uuid.UUID `json:"election_id"`
	UserID      uuid.UUID `json:"user_id"`
	CandidateNo int64     `json:"candidate_no"`
	Vision      string    `json:"vision"`

	Election *Election `gorm:"foreignKey:ElectionID" json:"election,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Timestamp
}

func (Candidate) TableName() string {
	return "candidates"
}
