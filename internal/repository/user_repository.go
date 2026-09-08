package repository

import (
	"context"
	"errors"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, metaReq pagination.Meta) (dto.GetAllUserRepositoryResponse, error)
	GetUserById(ctx context.Context, tx *gorm.DB, userId uuid.UUID) (entity.User, error)
	GetUserByNIA(ctx context.Context, tx *gorm.DB, userNIA string) (entity.User, error)
	Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	//DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error
	//ResetPassword(context.Context, string, string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return user, nil
}

func (r *userRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, metaReq pagination.Meta) (dto.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []entity.User

	tx = tx.Model(&entity.User{})

	if err := WithFilters(tx, &metaReq, AddModels(entity.User{}, "users")).
		Find(&users).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return dto.GetAllUserRepositoryResponse{
		Users: users,
		Meta:  metaReq,
	}, nil
}

func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId uuid.UUID) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("id = ?", userId).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, dto.ErrUserNotFound
		}
		return entity.User{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return user, nil
}

func (r *userRepository) GetUserByNIA(ctx context.Context, tx *gorm.DB, userNIA string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("nia = ?", userNIA).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, dto.ErrUserNotFound
		}
		return entity.User{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&user).Error; err != nil {
		return entity.User{}, MyError.Wrap(err, dto.ErrGeneral)
	}

	return user, nil
}
