package entity

import "github.com/google/uuid"

type MediaType string

const (
	MediaTypePhoto     MediaType = "photo"
	MediaTypePortfolio MediaType = "portfolio"
)

type CandidateMedia struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CandidateID uuid.UUID `gorm:"type:uuid;not null;index" json:"candidate_id"`

	Type     MediaType `gorm:"type:varchar(20);not null" json:"type"`
	FilePath string    `gorm:"type:text;not null" json:"file_path"`
	FileName string    `gorm:"type:text" json:"file_name"`
	MimeType string    `gorm:"type:varchar(100)" json:"mime_type"`
	FileSize int64     `json:"file_size"`

	Candidate *Candidate `gorm:"foreignKey:CandidateID;references:ID" json:"candidate"`
}

func (CandidateMedia) TableName() string {
	return "candidate_media"
}
