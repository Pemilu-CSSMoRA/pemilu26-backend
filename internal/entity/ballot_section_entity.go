package entity

import "github.com/google/uuid"

type BallotSelection struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	BallotID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_ballot_candidate" json:"ballot_id"`
	CandidateID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_ballot_candidate" json:"candidate_id"`

	Ballot    *Ballot    `gorm:"foreignKey:BallotID;references:ID" json:"ballot"`
	Candidate *Candidate `gorm:"foreignKey:CandidateID;references:ID" json:"candidate"`

	Timestamp
}

func (BallotSelection) TableName() string {
	return "ballot_selections"
}
