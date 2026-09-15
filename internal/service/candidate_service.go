package service

import (
	"context"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/repository"
)

type CandidateService interface {
	CreateCandidate(ctx context.Context)
}

type candidateService struct {
	CandidateRepo repository.CandidateRepository
}

func NewCandidateRepository(candidate repository.CandidateRepository) CandidateService {
	return &candidateService{
		CandidateRepo: candidate,
	}
}

func (s *candidateService) CreateCandidate(ctx context.Context) {

}
