package controller

import (
	"net/http"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/service"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	ElectionController interface {
		GetAllElection(ctx *gin.Context)
		CreateElection(ctx *gin.Context)
		UpdateElection(ctx *gin.Context)
		StartElection(ctx *gin.Context)
		EndElection(ctx *gin.Context)
		GetElectionVoters(ctx *gin.Context)
		EnrollVoters(ctx *gin.Context)
		UnenrollVoters(ctx *gin.Context)
	}

	electionController struct {
		electionService service.ElectionService
	}
)

func NewElectionController(electionService service.ElectionService) ElectionController {
	return &electionController{
		electionService: electionService,
	}
}

func (c *electionController) GetAllElection(ctx *gin.Context) {
	result, err := c.electionService.GetAllElectionWithPagination(ctx.Request.Context(), pagination.New(ctx))
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_ELECTION, err, nil).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_ELECTION, result.Data, result.Meta).Send(ctx)
}

func (c *electionController) CreateElection(ctx *gin.Context) {
	election := dto.ElectionCreateRequest{}
	if err := ctx.ShouldBind(&election); err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	result, err := c.electionService.CreateElection(ctx.Request.Context(), election)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ELECTION, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ELECTION, result).Send(ctx)
}

func (c *electionController) UpdateElection(ctx *gin.Context) {
	idParam := ctx.Param("id")

	electionID, err := uuid.Parse(idParam)
	if err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidElectionID, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ELECTION_ID, err).SendWithAbort(ctx)
		return
	}

	election := dto.ElectionUpdateRequest{}
	if err := ctx.ShouldBind(&election); err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	result, err := c.electionService.UpdateElection(ctx.Request.Context(), election, electionID)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ELECTION, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ELECTION, result).Send(ctx)
}

func (c *electionController) StartElection(ctx *gin.Context) {
	idParam := ctx.Param("id")

	electionID, err := uuid.Parse(idParam)
	if err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidElectionID, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ELECTION_ID, err).SendWithAbort(ctx)
		return
	}

	result, err := c.electionService.StartElection(ctx.Request.Context(), electionID)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_START_ELECTION, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_START_ELECTION, result).Send(ctx)
}

func (c *electionController) EndElection(ctx *gin.Context) {
	idParam := ctx.Param("id")

	electionID, err := uuid.Parse(idParam)
	if err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidElectionID, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ELECTION_ID, err).SendWithAbort(ctx)
		return
	}

	result, err := c.electionService.EndElection(ctx.Request.Context(), electionID)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_END_ELECTION, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_END_ELECTION, result).Send(ctx)
}

func (c *electionController) GetElectionVoters(ctx *gin.Context) {
	electionID, err := parseElectionID(ctx)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ELECTION_ID, err).SendWithAbort(ctx)
		return
	}

	result, err := c.electionService.GetElectionVotersWithPagination(ctx.Request.Context(), electionID, pagination.New(ctx))
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_ELECTION_VOTER, err, nil).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_ELECTION_VOTER, result.Data, result.Meta).Send(ctx)
}

func (c *electionController) EnrollVoters(ctx *gin.Context) {
	electionID, err := parseElectionID(ctx)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ELECTION_ID, err).SendWithAbort(ctx)
		return
	}

	userIDs, err := bindElectionVoteBatchRequest(ctx)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	result, err := c.electionService.EnrollVoters(ctx.Request.Context(), electionID, userIDs)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_ENROLL_ELECTION_VOTER, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_ENROLL_ELECTION_VOTER, result).Send(ctx)
}

func (c *electionController) UnenrollVoters(ctx *gin.Context) {
	electionID, err := parseElectionID(ctx)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ELECTION_ID, err).SendWithAbort(ctx)
		return
	}

	userIDs, err := bindElectionVoteBatchRequest(ctx)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	if err := c.electionService.UnenrollVoters(ctx.Request.Context(), electionID, userIDs); err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_UNENROLL_ELECTION_VOTER, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UNENROLL_ELECTION_VOTER, nil).Send(ctx)
}

func parseElectionID(ctx *gin.Context) (uuid.UUID, error) {
	electionID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return uuid.Nil, MyError.NewWrap(err, dto.ErrInvalidElectionID, http.StatusBadRequest)
	}
	return electionID, nil
}

func bindElectionVoteBatchRequest(ctx *gin.Context) ([]uuid.UUID, error) {
	request := dto.ElectionVoteBatchRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		return nil, MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
	}

	userIDs := make([]uuid.UUID, 0, len(request.UserIDs))
	for _, userID := range request.UserIDs {
		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			return nil, MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
		}
		userIDs = append(userIDs, parsedUserID)
	}

	return userIDs, nil
}
