package entity

import "github.com/google/uuid"

type ElectionVoter struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ElectionID      uuid.UUID `json:"election_id"`
	ElectionVoterID uuid.UUID `json:"voter_id"`

	Election      *Election `gorm:"foreignKey:ElectionID" json:"election"`
	ElectionVoter *User     `gorm:"foreignKey:ElectionVoterID" json:"voter"`
	Timestamp
}

func (ElectionVoter) TableName() string {
	return "election_voters"
}
