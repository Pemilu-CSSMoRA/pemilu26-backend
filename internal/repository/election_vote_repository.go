package repository

import (
	"context"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ElectionVoteRepository interface {
	GetVotersByElectionID(ctx context.Context, tx *gorm.DB, electionID uuid.UUID, metaReq pagination.Meta) (dto.GetAllElectionVoteRepositoryResponse, error)
	GetByElectionIDAndUserIDs(ctx context.Context, tx *gorm.DB, electionID uuid.UUID, userIDs []uuid.UUID) ([]entity.ElectionVote, error)
	CreateElectionVotes(ctx context.Context, tx *gorm.DB, electionVotes []entity.ElectionVote) ([]entity.ElectionVote, error)
	DeleteByElectionIDAndUserIDs(ctx context.Context, tx *gorm.DB, electionID uuid.UUID, userIDs []uuid.UUID) error
	WithinTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type electionVoteRepository struct {
	db *gorm.DB
}

func NewElectionVoteRepository(db *gorm.DB) ElectionVoteRepository {
	return &electionVoteRepository{db: db}
}

func (r *electionVoteRepository) GetVotersByElectionID(ctx context.Context, tx *gorm.DB, electionID uuid.UUID, metaReq pagination.Meta) (dto.GetAllElectionVoteRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var electionVotes []entity.ElectionVote
	query := tx.WithContext(ctx).Model(&entity.ElectionVote{}).Where("election_id = ?", electionID)
	if err := WithFilters(query, &metaReq, AddModels(entity.ElectionVote{}, "election_voters")).
		Preload("ElectionVoter").Find(&electionVotes).Error; err != nil {
		return dto.GetAllElectionVoteRepositoryResponse{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return dto.GetAllElectionVoteRepositoryResponse{ElectionVotes: electionVotes, Meta: metaReq}, nil
}

func (r *electionVoteRepository) GetByElectionIDAndUserIDs(ctx context.Context, tx *gorm.DB, electionID uuid.UUID, userIDs []uuid.UUID) ([]entity.ElectionVote, error) {
	if tx == nil {
		tx = r.db
	}

	var electionVotes []entity.ElectionVote
	if err := tx.WithContext(ctx).Where("election_id = ? AND user_id IN ?", electionID, userIDs).Find(&electionVotes).Error; err != nil {
		return nil, MyError.Wrap(err, dto.ErrGeneral)
	}
	return electionVotes, nil
}

func (r *electionVoteRepository) CreateElectionVotes(ctx context.Context, tx *gorm.DB, electionVotes []entity.ElectionVote) ([]entity.ElectionVote, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&electionVotes).Error; err != nil {
		return nil, MyError.Wrap(err, dto.ErrGeneral)
	}
	return electionVotes, nil
}

func (r *electionVoteRepository) DeleteByElectionIDAndUserIDs(ctx context.Context, tx *gorm.DB, electionID uuid.UUID, userIDs []uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	result := tx.WithContext(ctx).Where("election_id = ? AND user_id IN ?", electionID, userIDs).Delete(&entity.ElectionVote{})
	if result.Error != nil {
		return MyError.Wrap(result.Error, dto.ErrGeneral)
	}
	if result.RowsAffected != int64(len(userIDs)) {
		return dto.ErrElectionVoteNotFound
	}
	return nil
}

func (r *electionVoteRepository) WithinTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
