package entity

import "github.com/google/uuid"

type Candidate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ElectionID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_election_candidate;uniqueIndex:idx_no_election" json:"election_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_election_candidate" json:"user_id"`
	CandidateNo int64     `gorm:"uniqueIndex:idx_no_election" json:"candidate_no"`
	Vision      string    `json:"vision"`

	Election *Election          `gorm:"foreignKey:ElectionID" json:"election,omitempty"`
	User     *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Missions []CandidateMission `gorm:"foreignKey:CandidateID" json:"candidate_mission"`
	Media    []CandidateMedia   `gorm:"foreignKey:CandidateID;references:ID" json:"media"`

	Timestamp
}

func (Candidate) TableName() string {
	return "candidates"
}
