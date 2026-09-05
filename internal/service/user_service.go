package service

import (
	"context"
	"errors"
	"sync"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/repository"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/password"
	"gorm.io/gorm"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserRegisterResponse, error)
		Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
	}

	userService struct {
		UserRepo   repository.UserRepository
		jwtService JWTService
	}
)

var (
	mu sync.Mutex

	// ROUTE_FE = PRODUCTION LINK, use frontend link
	// ROUTE_BE = BE DEV LINK
	// VERIFY_EMAIL_ROUTE_FE = "verify-email"
	// VERIFY_EMAIL_ROUTE_BE = "api/auth/verify"
	// VERIFY_EMAIL_TEMPLATE = "utils/mailer/email-template/verification_email.html"
	// FORGET_EMAIL_TEMPLATE = "utils/mailer/email-template/forget_password_email.html"
	// FORGET_EMAIL_ROUTE    = "reset-password"
)

func NewUserService(userRepo repository.UserRepository, jwtService JWTService) UserService {
	return &userService{
		UserRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserRegisterResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	userCheck, _ := s.UserRepo.GetUserByNIA(ctx, nil, req.NIA)
	if userCheck.NIA != "" {
		return dto.UserRegisterResponse{}, dto.ErrNIAAlreadyExists
	}

	var UserRole entity.UserRole
	if req.Role == "admin" {
		UserRole = entity.RoleAdmin
	} else if req.Role == "bawaslu" {
		UserRole = entity.RoleBawaslu
	} else if req.Role == "user" {
		UserRole = entity.RoleUser
	} else {
		return dto.UserRegisterResponse{}, dto.ErrInvalidRole
	}
	user := entity.User{
		Name:     req.Name,
		NIA:      req.NIA,
		Password: req.Password,
		Angkatan: req.Angkatan,
		Role:     UserRole,
	}

	userReg, err := s.UserRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserRegisterResponse{}, dto.ErrCreateUser
	}

	return dto.UserRegisterResponse{
		ID:       userReg.ID.String(),
		Name:     userReg.Name,
		NIA:      userReg.NIA,
		Angkatan: userReg.Angkatan,
		Role:     string(userReg.Role),
	}, nil
}

func (s *userService) Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	check, err := s.UserRepo.GetUserByNIA(ctx, nil, req.NIA)
	if err != nil || check.NIA == "" {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserLoginResponse{}, dto.ErrInvalidCredentials
		}

		return dto.UserLoginResponse{}, err
	}

	checkPassword, err := password.CheckPassword(check.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		return dto.UserLoginResponse{}, dto.ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(check.ID, string(check.Role))

	return dto.UserLoginResponse{
		Token: token,
		Role:  string(check.Role),
	}, nil
}
