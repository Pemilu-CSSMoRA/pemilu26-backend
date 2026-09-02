package entity

import "github.com/google/uuid"

type Ballot struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ElectionID  uuid.UUID `json:"election_id"`
	VoterID     uuid.UUID `json:"voter_id"`
	CandidateID uuid.UUID `json:"candidate_id"`

	Election  *Election      `gorm:"foreignKey:ElectionID" json:"election"`
	Voter     *ElectionVoter `gorm:"foreignKey:VoterID" json:"voter"`
	Candidate *Candidate     `gorm:"foreignKey:CandidateID" json:"candidate"`
	Timestamp
}

func (Ballot) TableName() string {
	return "ballots"
}
