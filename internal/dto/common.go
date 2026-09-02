package dto

import (
	"net/http"

	MyError "github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/error"
)

const (
	MESSAGE_FAILED_PANIC_OCCURED      = "Server mengalami panic"
	MESSAGE_TO_MANY_REQUEST           = "Terlalu banyak permintaan"
	MESSAGE_FAILED_PARSE_TIME         = "Gagal mengurai waktu"
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "Gagal mengambil data dari body"
	MESSAGE_INVALID_NAME_FORMAT       = "Format nama tidak valid, hanya karakter dan huruf yang diperbolehkan"
	MESSAGE_INVALID_TEAM_NAME_FORMAT  = "Format nama tim tidak valid"
)

var (
	ErrGeneral               = MyError.New("Terjadi kesalahan", http.StatusInternalServerError)
	ErrInvalidInput          = MyError.New("Input tidak valid", http.StatusBadRequest)
	ErrToManyRequest         = MyError.New("Terlalu banyak permintaan", http.StatusTooManyRequests)
	ErrInvalidNameFormat     = MyError.New(MESSAGE_INVALID_NAME_FORMAT, http.StatusBadRequest)
	ErrInvalidTeamNameFormat = MyError.New(MESSAGE_INVALID_TEAM_NAME_FORMAT, http.StatusBadRequest)
)
