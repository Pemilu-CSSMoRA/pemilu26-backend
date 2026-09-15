package service

import (
	"context"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/repository"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ElectionService interface {
	GetAllElectionWithPagination(ctx context.Context, req pagination.Meta) (dto.ElectionPaginationResponse, error)
	CreateElection(ctx context.Context, election dto.ElectionCreateRequest) (dto.ElectionResponse, error)
	UpdateElection(ctx context.Context, election dto.ElectionUpdateRequest, electionID uuid.UUID) (dto.ElectionResponse, error)
	StartElection(ctx context.Context, electionID uuid.UUID) (dto.ElectionResponse, error)
	EndElection(ctx context.Context, electionID uuid.UUID) (dto.ElectionResponse, error)
	GetElectionVotersWithPagination(ctx context.Context, electionID uuid.UUID, req pagination.Meta) (dto.ElectionVotePaginationResponse, error)
	EnrollVoters(ctx context.Context, electionID uuid.UUID, userIDs []uuid.UUID) ([]dto.ElectionVoteResponse, error)
	UnenrollVoters(ctx context.Context, electionID uuid.UUID, userIDs []uuid.UUID) error
}

type electionService struct {
	ElectionRepo     repository.ElectionRepository
	UserRepo         repository.UserRepository
	ElectionVoteRepo repository.ElectionVoteRepository
}

func NewElectionService(electionRepo repository.ElectionRepository, userRepo repository.UserRepository, electionVoteRepo repository.ElectionVoteRepository) ElectionService {
	return &electionService{
		ElectionRepo:     electionRepo,
		UserRepo:         userRepo,
		ElectionVoteRepo: electionVoteRepo,
	}
}

func (s *electionService) GetAllElectionWithPagination(ctx context.Context, req pagination.Meta) (dto.ElectionPaginationResponse, error) {
	dataWithPaginate, err := s.ElectionRepo.GetAllElectionWithPagination(ctx, nil, req)
	if err != nil {
		return dto.ElectionPaginationResponse{}, err
	}

	var datas []dto.ElectionResponse
	for _, election := range dataWithPaginate.Elections {
		data := dto.ElectionResponse{
			ID:          election.ID.String(),
			Name:        election.Name,
			Description: election.Description,
			Year:        election.Year,
			Status:      string(election.Status),
		}
		datas = append(datas, data)
	}

	return dto.ElectionPaginationResponse{
		Data: datas,
		Meta: dataWithPaginate.Meta,
	}, nil
}

func (s *electionService) CreateElection(ctx context.Context, req dto.ElectionCreateRequest) (dto.ElectionResponse, error) {
	election := entity.Election{
		Name:        req.Name,
		Description: req.Description,
		Year:        req.Year,
		Status:      entity.ElectionCreated,
	}

	electionCreate, err := s.ElectionRepo.CreateElection(ctx, nil, election)
	if err != nil {
		return dto.ElectionResponse{}, dto.ErrCreateElection
	}

	return dto.ElectionResponse{
		ID:          electionCreate.ID.String(),
		Name:        electionCreate.Name,
		Description: electionCreate.Description,
		Year:        electionCreate.Year,
		Status:      string(electionCreate.Status),
	}, nil
}

func (s *electionService) UpdateElection(ctx context.Context, req dto.ElectionUpdateRequest, electionID uuid.UUID) (dto.ElectionResponse, error) {
	election, err := s.ElectionRepo.GetElectionById(ctx, nil, electionID)
	if err != nil {
		return dto.ElectionResponse{}, err
	}

	election.Name = req.Name
	election.Description = req.Description
	election.Year = req.Year

	electionUpdate, err := s.ElectionRepo.UpdateElection(ctx, nil, election)
	if err != nil {
		return dto.ElectionResponse{}, err
	}

	return dto.ElectionResponse{
		ID:          electionUpdate.ID.String(),
		Name:        electionUpdate.Name,
		Description: electionUpdate.Description,
		Year:        electionUpdate.Year,
		Status:      string(electionUpdate.Status),
	}, nil
}

func (s *electionService) StartElection(ctx context.Context, electionID uuid.UUID) (dto.ElectionResponse, error) {
	election, err := s.ElectionRepo.GetElectionById(ctx, nil, electionID)
	if err != nil {
		return dto.ElectionResponse{}, err
	}

	if election.Status != entity.ElectionCreated {
		return dto.ElectionResponse{}, dto.ErrElectionStatusNotAllowed
	}

	election.Status = entity.ElectionStarted
	electionUpdate, err := s.ElectionRepo.UpdateElection(ctx, nil, election)
	if err != nil {
		return dto.ElectionResponse{}, err
	}

	return dto.ElectionResponse{
		ID:          electionUpdate.ID.String(),
		Name:        electionUpdate.Name,
		Description: electionUpdate.Description,
		Year:        electionUpdate.Year,
		Status:      string(electionUpdate.Status),
	}, nil
}

func (s *electionService) EndElection(ctx context.Context, electionID uuid.UUID) (dto.ElectionResponse, error) {
	election, err := s.ElectionRepo.GetElectionById(ctx, nil, electionID)
	if err != nil {
		return dto.ElectionResponse{}, err
	}

	if election.Status != entity.ElectionStarted {
		return dto.ElectionResponse{}, dto.ErrElectionStatusNotAllowed
	}

	election.Status = entity.ElectionEnded
	electionUpdate, err := s.ElectionRepo.UpdateElection(ctx, nil, election)
	if err != nil {
		return dto.ElectionResponse{}, err
	}

	return dto.ElectionResponse{
		ID:          electionUpdate.ID.String(),
		Name:        electionUpdate.Name,
		Description: electionUpdate.Description,
		Year:        electionUpdate.Year,
		Status:      string(electionUpdate.Status),
	}, nil
}

func (s *electionService) GetElectionVotersWithPagination(ctx context.Context, electionID uuid.UUID, req pagination.Meta) (dto.ElectionVotePaginationResponse, error) {
	if _, err := s.ElectionRepo.GetElectionById(ctx, nil, electionID); err != nil {
		return dto.ElectionVotePaginationResponse{}, err
	}

	dataWithPaginate, err := s.ElectionVoteRepo.GetVotersByElectionID(ctx, nil, electionID, req)
	if err != nil {
		return dto.ElectionVotePaginationResponse{}, err
	}

	return dto.ElectionVotePaginationResponse{
		Data: mapElectionVoteResponses(dataWithPaginate.ElectionVotes),
		Meta: dataWithPaginate.Meta,
	}, nil
}

func (s *electionService) EnrollVoters(ctx context.Context, electionID uuid.UUID, userIDs []uuid.UUID) ([]dto.ElectionVoteResponse, error) {
	if err := validateUniqueUserIDs(userIDs); err != nil {
		return nil, err
	}

	var enrolled []entity.ElectionVote
	err := s.ElectionVoteRepo.WithinTransaction(ctx, func(tx *gorm.DB) error {
		if _, err := s.ElectionRepo.GetElectionById(ctx, tx, electionID); err != nil {
			return err
		}
		for _, userID := range userIDs {
			if _, err := s.UserRepo.GetUserById(ctx, tx, userID); err != nil {
				return err
			}
		}

		existing, err := s.ElectionVoteRepo.GetByElectionIDAndUserIDs(ctx, tx, electionID, userIDs)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return dto.ErrElectionVoteAlreadyEnrolled
		}

		votes := make([]entity.ElectionVote, 0, len(userIDs))
		for _, userID := range userIDs {
			votes = append(votes, entity.ElectionVote{ElectionID: electionID, UserID: userID, Status: false})
		}
		enrolled, err = s.ElectionVoteRepo.CreateElectionVotes(ctx, tx, votes)
		return err
	})
	if err != nil {
		return nil, err
	}

	return mapElectionVoteResponses(enrolled), nil
}

func (s *electionService) UnenrollVoters(ctx context.Context, electionID uuid.UUID, userIDs []uuid.UUID) error {
	if err := validateUniqueUserIDs(userIDs); err != nil {
		return err
	}

	return s.ElectionVoteRepo.WithinTransaction(ctx, func(tx *gorm.DB) error {
		if _, err := s.ElectionRepo.GetElectionById(ctx, tx, electionID); err != nil {
			return err
		}

		existing, err := s.ElectionVoteRepo.GetByElectionIDAndUserIDs(ctx, tx, electionID, userIDs)
		if err != nil {
			return err
		}
		if len(existing) != len(userIDs) {
			return dto.ErrElectionVoteNotFound
		}
		for _, electionVote := range existing {
			if electionVote.Status {
				return dto.ErrElectionVoteAlreadyVoted
			}
		}

		return s.ElectionVoteRepo.DeleteByElectionIDAndUserIDs(ctx, tx, electionID, userIDs)
	})
}

func validateUniqueUserIDs(userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return dto.ErrInvalidInput
	}

	seen := make(map[uuid.UUID]struct{}, len(userIDs))
	for _, userID := range userIDs {
		if _, exists := seen[userID]; exists {
			return dto.ErrDuplicateUserIDs
		}
		seen[userID] = struct{}{}
	}

	return nil
}

func mapElectionVoteResponses(electionVotes []entity.ElectionVote) []dto.ElectionVoteResponse {
	responses := make([]dto.ElectionVoteResponse, 0, len(electionVotes))
	for _, electionVote := range electionVotes {
		response := dto.ElectionVoteResponse{
			ID:         electionVote.ID.String(),
			ElectionID: electionVote.ElectionID,
			UserID:     electionVote.UserID,
			Status:     electionVote.Status,
		}
		if electionVote.ElectionVoter != nil {
			response.Voter = &dto.UserResponse{
				ID:       electionVote.ElectionVoter.ID.String(),
				Name:     electionVote.ElectionVoter.Name,
				NIA:      electionVote.ElectionVoter.NIA,
				Angkatan: electionVote.ElectionVoter.Angkatan,
				Role:     string(electionVote.ElectionVoter.Role),
			}
		}
		responses = append(responses, response)
	}

	return responses
}
