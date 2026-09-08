package dto

import (
	"net/http"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/pagination"
)

const (
	// Failed
	MESSAGE_FAILED_REGISTER_USER           = "failed create user"
	MESSAGE_FAILED_GET_LIST_USER           = "failed get list user"
	MESSAGE_FAILED_GET_USER_TOKEN          = "failed get user token"
	MESSAGE_FAILED_GET_USER                = "failed get user"
	MESSAGE_FAILED_GET_USER_ID             = "failed get user id"
	MESSAGE_FAILED_LOGIN                   = "failed login"
	MESSAGE_FAILED_WRONG_NIA_OR_PASSWORD   = "wrong nia or password"
	MESSAGE_FAILED_UPDATE_USER             = "failed update user"
	MESSAGE_FAILED_DELETE_USER             = "failed delete user"
	MESSAGE_FAILED_PROSES_REQUEST          = "failed proses request"
	MESSAGE_FAILED_DENIED_ACCESS           = "denied access"
	MESSAGE_FAILED_VERIFY_EMAIL            = "failed verify email"
	MESSAGE_FAILED_RESET_PASSWORD          = "failed reset password"
	MESSAGE_FAILED_EMAIL_NOT_FOUND         = "email not found"
	MESSAGE_FAILED_FORGET_PASSWORD         = "failed handle forget password"
	MESSAGE_FAILED_GET_VERIFICATION_STATUS = "failed get verification status"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER           = "success create user"
	MESSAGE_SUCCESS_GET_LIST_USER           = "success get list user"
	MESSAGE_SUCCESS_GET_USER                = "success get user"
	MESSAGE_SUCCESS_LOGIN                   = "success login"
	MESSAGE_SUCCESS_UPDATE_USER             = "success update user"
	MESSAGE_SUCCESS_DELETE_USER             = "success delete user"
	MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS = "success send verification email"
	MESSAGE_SUCCESS_VERIFY_EMAIL            = "success verify email"
	MESSAGE_SUCCESS_RESET_PASSWORD          = "success reset password"
	MESSAGE_SUCCESS_FORGET_PASSWORD         = "success handle forget password"
	MESSAGE_SUCCESS_GET_VERIFICATION_STATUS = "success get verification status"
)

var (
	ErrInvalidUserID      = MyError.New("invalid user id", http.StatusInternalServerError)
	ErrCreateUser         = MyError.New("gagal membuat user", http.StatusInternalServerError)
	ErrInvalidRole        = MyError.New("invalid role", http.StatusBadRequest)
	ErrGetAllUser         = MyError.New("gagal mengambil semua user", http.StatusInternalServerError)
	ErrGetUserById        = MyError.New("gagal mengambil user berdasarkan ID", http.StatusBadRequest)
	ErrGetUserByNIA       = MyError.New("gagal mengambil user berdasarkan NIA", http.StatusBadRequest)
	ErrNIAAlreadyExists   = MyError.New("NIA sudah digunakan", http.StatusConflict)
	ErrUpdateUser         = MyError.New("gagal memperbarui user", http.StatusInternalServerError)
	ErrUserNotAdmin       = MyError.New("user bukan admin", http.StatusForbidden)
	ErrUserNotFound       = MyError.New("user tidak ditemukan", http.StatusBadRequest)
	ErrInvalidCredentials = MyError.New("nia atau password salah", http.StatusUnauthorized)
	ErrDeleteUser         = MyError.New("gagal menghapus user", http.StatusInternalServerError)
	// ErrPasswordNotMatch       = MyError.New("password tidak cocok", http.StatusBadRequest)
	// ErrEmailOrPassword        = MyError.New("email atau password salah", http.StatusUnauthorized)
	ErrAccountNotVerified     = MyError.New("akun belum diverifikasi", http.StatusForbidden)
	ErrAccountAlreadyVerified = MyError.New("akun sudah diverifikasi", http.StatusConflict)
	ErrHashPasswordFailed     = MyError.New("gagal melakukan hash password", http.StatusInternalServerError)
	ErrParseUUID              = MyError.New("gagal parsing UUID", http.StatusBadRequest)
	ErrUserIdEmpty            = MyError.New("ID user kosong", http.StatusBadRequest)
	ErrUserIdEmptyString      = MyError.New("ID user berupa string kosong", http.StatusBadRequest)
	ErrUserIdInvalid          = MyError.New("format ID user tidak valid", http.StatusBadRequest)
)

type (
	UserPaginationResponse struct {
		Data []UserResponse `json:"data"`
		pagination.Meta
	}

	GetAllUserRepositoryResponse struct {
		Users []entity.User
		pagination.Meta
	}
	UserResponse struct {
		ID       string `json:"id" binding:"required"`
		Name     string `json:"name" binding:"required"`
		NIA      string `json:"nia" binding:"required"`
		Angkatan string `json:"angkatan" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}
	UserRegisterRequest struct {
		Name     string `json:"name" form:"name" binding:"required,max=150"`
		NIA      string `json:"nia" form:"nia" binding:"required,max=100"`
		Angkatan string `json:"angkatan" form:"angkatan" binding:"required,max=100"`
		Password string `json:"password" form:"password" binding:"required,max=100"`
	}

	UserLoginRequest struct {
		NIA      string `json:"nia" form:"nia" binding:"required,max=100"`
		Password string `json:"password" form:"password" binding:"required,max=100"`
	}

	UserLoginResponse struct {
		Token string `json:"token" binding:"required"`
		Role  string `json:"role" binding:"required"`
	}

	UserUpdateMeRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		NIA      string `json:"nia" form:"nia" binding:"required"`
		Angkatan string `json:"angkatan" form:"angkatan" binding:"required"`
	}
	UserUpdateAdminRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		NIA      string `json:"nia" form:"nia" binding:"required"`
		Angkatan string `json:"angkatan" form:"angkatan" binding:"required"`
		Role     string `json:"role" form:"role" binding:"required"`
	}
)
