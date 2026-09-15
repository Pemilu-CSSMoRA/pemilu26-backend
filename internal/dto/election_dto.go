package dto

import (
	"net/http"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/google/uuid"
)

const (
	// Election

	// Failed
	MESSAGE_FAILED_CREATE_ELECTION         = "failed create election"
	MESSAGE_FAILED_GET_LIST_ELECTION       = "failed get list election"
	MESSAGE_FAILED_UPDATE_ELECTION         = "failed update election"
	MESSAGE_FAILED_GET_ELECTION_ID         = "failed get election id"
	MESSAGE_FAILED_START_ELECTION          = "failed start election"
	MESSAGE_FAILED_END_ELECTION            = "failed end election"
	MESSAGE_FAILED_GET_LIST_ELECTION_VOTER = "failed get list election voter"
	MESSAGE_FAILED_ENROLL_ELECTION_VOTER   = "failed enroll election voter"
	MESSAGE_FAILED_UNENROLL_ELECTION_VOTER = "failed unenroll election voter"

	// Success
	MESSAGE_SUCCESS_CREATE_ELECTION         = "success create election"
	MESSAGE_SUCCESS_GET_LIST_ELECTION       = "success get list election"
	MESSAGE_SUCCESS_UPDATE_ELECTION         = "success update election"
	MESSAGE_SUCCESS_START_ELECTION          = "success start election"
	MESSAGE_SUCCESS_END_ELECTION            = "success end election"
	MESSAGE_SUCCESS_GET_LIST_ELECTION_VOTER = "success get list election voter"
	MESSAGE_SUCCESS_ENROLL_ELECTION_VOTER   = "success enroll election voter"
	MESSAGE_SUCCESS_UNENROLL_ELECTION_VOTER = "success unenroll election voter"
)

var (
	ErrInvalidElectionID           = MyError.New("invalid election id", http.StatusBadRequest)
	ErrCreateElection              = MyError.New("gagal membuat election", http.StatusInternalServerError)
	ErrInvalidElectionStatus       = MyError.New("invalid election status", http.StatusBadRequest)
	ErrElectionNotFound            = MyError.New("election tidak ditemukan", http.StatusNotFound)
	ErrElectionStatusNotAllowed    = MyError.New("election status tidak dapat diubah", http.StatusConflict)
	ErrElectionVoteAlreadyEnrolled = MyError.New("user sudah terdaftar pada election ini", http.StatusConflict)
	ErrElectionVoteNotFound        = MyError.New("data voter election tidak ditemukan", http.StatusNotFound)
	ErrElectionVoteAlreadyVoted    = MyError.New("voter yang sudah memilih tidak dapat dihapus", http.StatusConflict)
	ErrDuplicateUserIDs            = MyError.New("user_ids tidak boleh duplikat", http.StatusBadRequest)
)

type (
	ElectionPaginationResponse struct {
		Data []ElectionResponse `json:"data"`
		pagination.Meta
	}

	GetAllElectionRepositoryResponse struct {
		Elections []entity.Election
		pagination.Meta
	}

	ElectionResponse struct {
		ID          string `json:"id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Status      string `json:"status" binding:"required"`
		Year        string `json:"year" binding:"required"`
	}

	ElectionCreateRequest struct {
		Name        string `json:"name" form:"name" binding:"required,max=150"`
		Description string `json:"description" form:"description" binding:"required,max=200"`
		Year        string `json:"year" binding:"required"`
	}

	ElectionUpdateRequest struct {
		Name        string `json:"name" form:"name" binding:"required"`
		Description string `json:"description" form:"description" binding:"required"`
		Year        string `json:"year" binding:"required"`
	}

	ElectionVoteResponse struct {
		ID         string        `json:"id"`
		ElectionID uuid.UUID     `json:"election_id"`
		UserID     uuid.UUID     `json:"user_id"`
		Status     bool          `json:"status"`
		Voter      *UserResponse `json:"voter,omitempty"`
	}

	ElectionVotePaginationResponse struct {
		Data []ElectionVoteResponse `json:"data"`
		pagination.Meta
	}

	GetAllElectionVoteRepositoryResponse struct {
		ElectionVotes []entity.ElectionVote
		pagination.Meta
	}

	ElectionVoteBatchRequest struct {
		UserIDs []string `json:"user_ids" binding:"required,min=1,dive,uuid"`
	}
)
