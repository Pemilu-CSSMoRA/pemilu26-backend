package repository

import (
	"context"
	"errors"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error)
	GetUserByNIA(ctx context.Context, tx *gorm.DB, userNIA string) (entity.User, error)
	//UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
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

func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error) {
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
