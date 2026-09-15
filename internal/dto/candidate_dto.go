package dto

import (
	"net/http"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
)

const (
	// Failed
	MESSAGE_FAILED_CREATE_CANDIDATE   = "failed create candidate"
	MESSAGE_FAILED_GET_LIST_CANDIDATE = "failed get list candidate"
	MESSAGE_FAILED_UPDATE_CANDIDATE   = "failed update candidate"
	MESSAGE_FAILED_GET_CANDIDATE_ID   = "failed get candidate id"

	// Success
	MESSAGE_SUCCESS_CREATE_CANDIDATE   = "success create candidate"
	MESSAGE_SUCCESS_GET_LIST_CANDIDATE = "success get list candidate"
	MESSAGE_SUCCESS_UPDATE_CANDIDATE   = "success update candidate"
)

var (
	ErrInvalidCandidateID        = MyError.New("invalid candidate id", http.StatusBadRequest)
	ErrCreateCandidate           = MyError.New("gagal membuat candidate", http.StatusInternalServerError)
	ErrInvalidCandidateStatus    = MyError.New("invalid candidate status", http.StatusBadRequest)
	ErrCandidateNotFound         = MyError.New("candidate tidak ditemukan", http.StatusNotFound)
	ErrCandidateStatusNotAllowed = MyError.New("candidate status tidak dapat diubah", http.StatusConflict)
)

type (
	CandidatePaginationResponse struct {
		Data []ElectionResponse `json:"data"`
		pagination.Meta
	}

	GetAllCandidateRepositoryResponse struct {
		Candidates []entity.Candidate
		pagination.Meta
	}

	CandidateResponse struct {
		ID          string           `json:"id"`
		User        UserResponse     `json:"user"`
		Election    ElectionResponse `json:"election"`
		Name        string           `json:"name"`
		CandidateNo int64            `json:"candidate_no"`
		Status      string           `json:"status"`
	}

	CandidateCreateRequest struct {
		Name        string `json:"name" form:"name" binding:"required,max=150"`
		Description string `json:"description" form:"description" binding:"required,max=200"`
	}

	CandidateUpdateRequest struct {
		Name        string `json:"name" form:"name" binding:"required"`
		Description string `json:"description" form:"description" binding:"required"`
	}
)
