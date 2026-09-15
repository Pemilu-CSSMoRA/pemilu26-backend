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

type CandidateRepository interface {
	GetAllCandidateWithPagination(ctx context.Context, tx *gorm.DB, metaReq pagination.Meta) (dto.GetAllCandidateRepositoryResponse, error)
	GetCandidateById(ctx context.Context, tx *gorm.DB, candidateID uuid.UUID) (entity.Candidate, error)
	CreateCandidate(ctx context.Context, tx *gorm.DB, candidate entity.Candidate) (entity.Candidate, error)
	UpdateCandidate(ctx context.Context, tx *gorm.DB, candidate entity.Candidate) (entity.Candidate, error)
}

type candidateRepository struct {
	db *gorm.DB
}

func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &candidateRepository{
		db: db,
	}
}

func (r *candidateRepository) GetCandidateById(ctx context.Context, tx *gorm.DB, candidateID uuid.UUID) (entity.Candidate, error) {
	if tx == nil {
		tx = r.db
	}

	var candidate entity.Candidate
	if err := tx.WithContext(ctx).Where("id = ?", candidateID).Take(&candidate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Candidate{}, dto.ErrCandidateNotFound
		}
		return entity.Candidate{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return candidate, nil
}

func (r *candidateRepository) GetAllCandidateWithPagination(ctx context.Context, tx *gorm.DB, metaReq pagination.Meta) (dto.GetAllCandidateRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var candidates []entity.Candidate

	tx = tx.Model(&entity.Candidate{})

	if err := WithFilters(tx, &metaReq, AddModels(entity.Candidate{}, "candidates")).
		Find(&candidates).Error; err != nil {
		return dto.GetAllCandidateRepositoryResponse{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return dto.GetAllCandidateRepositoryResponse{
		Candidates: candidates,
		Meta:       metaReq,
	}, nil
}

func (r *candidateRepository) CreateCandidate(ctx context.Context, tx *gorm.DB, candidate entity.Candidate) (entity.Candidate, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&candidate).Error; err != nil {
		return entity.Candidate{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return candidate, nil
}

func (r *candidateRepository) UpdateCandidate(ctx context.Context, tx *gorm.DB, candidate entity.Candidate) (entity.Candidate, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&candidate).Error; err != nil {
		return entity.Candidate{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return candidate, nil
}
