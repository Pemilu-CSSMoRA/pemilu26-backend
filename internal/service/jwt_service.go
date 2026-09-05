package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jWTService struct {
	SecretKey   string
	ExpireHours int
}

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`

	jwt.RegisteredClaims
}

func NewJWTService(secretKey string, expireHours int) JWTService {
	return &jWTService{
		SecretKey:   secretKey,
		ExpireHours: expireHours,
	}
}

func (j *jWTService) GenerateToken(userID uuid.UUID, role string) (string, error) {

	now := time.Now()

	claims := JWTClaims{
		UserID: userID,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(j.ExpireHours) * time.Hour),
			),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(j.SecretKey))
}

func (j *jWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}

			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
