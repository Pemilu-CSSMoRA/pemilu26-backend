package entity

import "github.com/google/uuid"

type VoterStatus string

const (
	Enrolled VoterStatus = "enrolled"
	Voted    VoterStatus = "voted"
)

type ElectionVote struct {
	ID         uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ElectionID uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_election_voter" json:"election_id"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_election_voter" json:"voter_id"`
	Status     VoterStatus `json:"status"`

	Election      *Election `gorm:"foreignKey:ElectionID" json:"election"`
	ElectionVoter *User     `gorm:"foreignKey:UserID" json:"voter"`
	Timestamp
}

func (ElectionVote) TableName() string {
	return "election_voters"
}
