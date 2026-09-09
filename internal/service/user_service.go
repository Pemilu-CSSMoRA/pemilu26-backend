package service

import (
	"context"
	"errors"
	"sync"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/repository"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/password"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	UserService interface {
		GetUserById(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req pagination.Meta) (dto.UserPaginationResponse, error)
		RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserResponse, error)
		Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
		UpdateMe(ctx context.Context, req dto.UserUpdateMeRequest, userId uuid.UUID) (dto.UserResponse, error)
		UpdateAdmin(ctx context.Context, req dto.UserUpdateAdminRequest, userID uuid.UUID) (dto.UserResponse, error)
		ResetPasswordMe(ctx context.Context, req dto.UserUpdatePasswordMeRequest, userID uuid.UUID) error
		ResetPasswordAdmin(ctx context.Context, req dto.UserUpdatePasswordAdminRequest, userID uuid.UUID) error
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

func (s *userService) GetUserById(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error) {
	check, err := s.UserRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       check.ID.String(),
		Name:     check.Name,
		NIA:      check.NIA,
		Angkatan: check.Angkatan,
		Role:     string(check.Role),
	}, nil
}

func (s *userService) GetAllUserWithPagination(ctx context.Context, req pagination.Meta) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := s.UserRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.UserResponse{
			ID:       user.ID.String(),
			Name:     user.Name,
			NIA:      user.NIA,
			Angkatan: user.Angkatan,
			Role:     string(user.Role),
		}
		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		Meta: dataWithPaginate.Meta,
	}, nil
}

func (s *userService) RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	userCheck, _ := s.UserRepo.GetUserByNIA(ctx, nil, req.NIA)
	if userCheck.NIA != "" {
		return dto.UserResponse{}, dto.ErrNIAAlreadyExists
	}

	var UserRole entity.UserRole
	UserRole = entity.RoleUser
	user := entity.User{
		Name:     req.Name,
		NIA:      req.NIA,
		Password: req.Password,
		Angkatan: req.Angkatan,
		Role:     UserRole,
	}

	userReg, err := s.UserRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
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

func (s *userService) UpdateMe(ctx context.Context, req dto.UserUpdateMeRequest, userId uuid.UUID) (dto.UserResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	user, err := s.UserRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	if req.Name != user.Name {
		user.Name = req.Name
	}
	if req.NIA != user.NIA {
		check, err := s.UserRepo.GetUserByNIA(ctx, nil, req.NIA)
		if err != nil {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}

		if check.NIA == req.NIA {
			return dto.UserResponse{}, dto.ErrNIAAlreadyExists
		}
		user.NIA = req.NIA
	}
	if req.Angkatan != user.Angkatan {
		user.Angkatan = req.Angkatan
	}

	update, err := s.UserRepo.Update(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       update.ID.String(),
		Name:     update.Name,
		NIA:      update.NIA,
		Angkatan: update.Angkatan,
		Role:     string(update.Role),
	}, nil
}

func (s *userService) UpdateAdmin(ctx context.Context, req dto.UserUpdateAdminRequest, userID uuid.UUID) (dto.UserResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	user, err := s.UserRepo.GetUserById(ctx, nil, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	if req.Name != user.Name {
		user.Name = req.Name
	}
	if req.NIA != user.NIA {
		check, err := s.UserRepo.GetUserByNIA(ctx, nil, req.NIA)
		if err != nil {
			return dto.UserResponse{}, dto.ErrGeneral
		}

		if check.NIA == req.NIA {
			return dto.UserResponse{}, dto.ErrNIAAlreadyExists
		}
		user.NIA = req.NIA
	}
	if req.Angkatan != user.Angkatan {
		user.Angkatan = req.Angkatan
	}

	switch req.Role {
	case "admin":
		user.Role = entity.RoleAdmin
	case "bawaslu":
		user.Role = entity.RoleBawaslu
	case "user":
		user.Role = entity.RoleUser
	default:
		return dto.UserResponse{}, dto.ErrInvalidRole

	}

	update, err := s.UserRepo.Update(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       update.ID.String(),
		Name:     update.Name,
		NIA:      update.NIA,
		Angkatan: update.Angkatan,
		Role:     string(update.Role),
	}, nil
}

func (s *userService) ResetPasswordMe(ctx context.Context, req dto.UserUpdatePasswordMeRequest, userID uuid.UUID) error {
	// get user data
	user, err := s.UserRepo.GetUserById(ctx, nil, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrUserNotFound
		}

		return err
	}

	// check old password
	checkPassword, err := password.CheckPassword(user.Password, []byte(req.OldPassword))
	if err != nil || !checkPassword {
		return dto.ErrInvalidCredentials
	}

	// check match new password
	if req.NewPassword != req.ConfirmPassword {
		return dto.ErrMissMatchNewPassword
	}

	// check old and new password
	if req.NewPassword == req.OldPassword {
		return dto.ErrSamePassword
	}

	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return dto.ErrHashPasswordFailed
	}

	user.Password = hashedPassword
	if _, err := s.UserRepo.Update(ctx, nil, user); err != nil {
		return err
	}

	return nil

}

func (s *userService) ResetPasswordAdmin(ctx context.Context, req dto.UserUpdatePasswordAdminRequest, userID uuid.UUID) error {
	// get user data
	user, err := s.UserRepo.GetUserById(ctx, nil, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrUserNotFound
		}

		return err
	}

	// check prevent same password
	samePassword, err := password.CheckPassword(
		user.Password,
		[]byte(req.NewPassword),
	)

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			samePassword = false
		} else {
			return err
		}
	}

	if samePassword {
		return dto.ErrSamePassword
	}

	// check match new password
	if req.NewPassword != req.ConfirmPassword {
		return dto.ErrMissMatchNewPassword
	}

	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return dto.ErrHashPasswordFailed
	}

	user.Password = hashedPassword
	if _, err := s.UserRepo.Update(ctx, nil, user); err != nil {
		return err
	}

	return nil
}
