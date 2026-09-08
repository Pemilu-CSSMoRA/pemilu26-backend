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
	UserController interface {
		GetMe(ctx *gin.Context)
		GetAllUser(ctx *gin.Context)
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
		UpdateMe(ctx *gin.Context)
		UpdateAdmin(ctx *gin.Context)
	}

	userController struct {
		userService service.UserService
	}
)

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) GetAllUser(ctx *gin.Context) {
	result, err := c.userService.GetAllUserWithPagination(ctx.Request.Context(), pagination.New(ctx))
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER, err, nil).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER, result.Data, result.Meta).Send(ctx)
}
func (c *userController) GetMe(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "user id not found",
		})
		return
	}

	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user id",
		})
		return
	}

	result, err := c.userService.GetUserById(ctx.Request.Context(), userId)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result).Send(ctx)
}
func (c *userController) Register(ctx *gin.Context) {
	user := dto.UserRegisterRequest{}
	if err := ctx.ShouldBind(&user); err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	result, err := c.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result).Send(ctx)
}

func (c *userController) Login(ctx *gin.Context) {
	user := dto.UserLoginRequest{}
	if err := ctx.ShouldBind(&user); err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	result, err := c.userService.Verify(ctx.Request.Context(), user)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_LOGIN, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGIN, result).Send(ctx)
}

func (c *userController) UpdateMe(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "user id not found",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user id",
		})
		return
	}

	user := dto.UserUpdateMeRequest{}
	if err := ctx.ShouldBind(&user); err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	result, err := c.userService.UpdateMe(ctx.Request.Context(), user, userID)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result).Send(ctx)
}

func (c *userController) UpdateAdmin(ctx *gin.Context) {
	idParam := ctx.Param("id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidUserID, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER_ID, err).SendWithAbort(ctx)
		return
	}
	user := dto.UserUpdateAdminRequest{}
	if err := ctx.ShouldBind(&user); err != nil {
		err = MyError.NewWrap(err, dto.ErrInvalidInput, http.StatusBadRequest)
		response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err).SendWithAbort(ctx)
		return
	}

	result, err := c.userService.UpdateAdmin(ctx.Request.Context(), user, userID)
	if err != nil {
		response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err).Send(ctx)
		return
	}

	response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result).Send(ctx)
}
