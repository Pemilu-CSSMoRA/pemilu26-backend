package repository

import (
	"context"
	"errors"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ElectionRepository interface {
	GetAllElectionWithPagination(ctx context.Context, tx *gorm.DB, metaReq pagination.Meta) (dto.GetAllElectionRepositoryResponse, error)
	GetElectionById(ctx context.Context, tx *gorm.DB, electionID uuid.UUID) (entity.Election, error)
	CreateElection(ctx context.Context, tx *gorm.DB, election entity.Election) (entity.Election, error)
	UpdateElection(ctx context.Context, tx *gorm.DB, election entity.Election) (entity.Election, error)
}

type electionRepository struct {
	db *gorm.DB
}

func NewElectionRepository(db *gorm.DB) ElectionRepository {
	return &electionRepository{
		db: db,
	}
}

func (r *electionRepository) GetElectionById(ctx context.Context, tx *gorm.DB, electionID uuid.UUID) (entity.Election, error) {
	if tx == nil {
		tx = r.db
	}

	var election entity.Election
	if err := tx.WithContext(ctx).Where("id = ?", electionID).Take(&election).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Election{}, dto.ErrElectionNotFound
		}
		return entity.Election{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return election, nil
}

func (r *electionRepository) GetAllElectionWithPagination(ctx context.Context, tx *gorm.DB, metaReq pagination.Meta) (dto.GetAllElectionRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var elections []entity.Election

	tx = tx.Model(&entity.Election{})

	if err := WithFilters(tx, &metaReq, AddModels(entity.Election{}, "elections")).
		Find(&elections).Error; err != nil {
		return dto.GetAllElectionRepositoryResponse{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return dto.GetAllElectionRepositoryResponse{
		Elections: elections,
		Meta:      metaReq,
	}, nil
}

func (r *electionRepository) CreateElection(ctx context.Context, tx *gorm.DB, election entity.Election) (entity.Election, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&election).Error; err != nil {
		return entity.Election{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return election, nil
}

func (r *electionRepository) UpdateElection(ctx context.Context, tx *gorm.DB, election entity.Election) (entity.Election, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&election).Error; err != nil {
		return entity.Election{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return election, nil
}
