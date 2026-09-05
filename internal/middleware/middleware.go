package middleware

import (
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/service"
	"gorm.io/gorm"
)

type Middleware struct {
	db         *gorm.DB
	jwtService service.JWTService
}

func New(db *gorm.DB, jwtService service.JWTService) Middleware {
	return Middleware{db, jwtService}
}
