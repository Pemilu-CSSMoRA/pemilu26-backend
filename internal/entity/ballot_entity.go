package entity

import "github.com/google/uuid"

type Ballot struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ElectionID uuid.UUID `gorm:"type:uuid:not null;uniqueIndex:idx_ballot_election" json:"election_id"`
	VoterID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;idx_ballot_election" json:"voter_id"`

	Election   *Election         `gorm:"foreignKey:ElectionID" json:"election"`
	Selections []BallotSelection `gorm:"foreignKey:BallotID" json:"selections"`
	Voter      *ElectionVote     `gorm:"foreignKey:VoterID" json:"voter"`

	Timestamp
}

func (Ballot) TableName() string {
	return "ballots"
}
