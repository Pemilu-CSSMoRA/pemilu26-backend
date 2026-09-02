package migrations

import (
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Election{},
		&entity.ElectionVoter{},
		&entity.Candidate{},
		&entity.CandidateMission{},
		&entity.CandidateMedia{},
		&entity.Ballot{},
	); err != nil {
		return err
	}

	return nil
}
