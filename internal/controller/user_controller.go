package controller

import (
	"net/http"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/service"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type (
	UserController interface {
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
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
