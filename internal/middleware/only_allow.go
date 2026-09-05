package middleware

import (
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func (m Middleware) OnlyAllow(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.MustGet("role").(string)

		for _, role := range roles {
			if userRole == role {
				ctx.Next()
				return
			}
		}

		response.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, dto.ErrRoleNotAllowed).
			SendWithAbort(ctx)
	}
}
