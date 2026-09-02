package MyError

import (
	"errors"
	"os"
)

type FieldError struct {
	Field   string `json:"field"`
	Label   string `json:"label"`
	Message string `json:"message"`
	Rule    string `json:"rule,omitempty"`
}

type Error struct {
	Message    string        `json:"Message"`
	StatusCode int           `json:"StatusCode"`
	ErrorCode  string        `json:"ErrorCode,omitempty"`
	Fields     *[]FieldError `json:"Fields,omitempty"`
}

func New(msg string, statusCode int) Error {
	return Error{
		Message:    msg,
		StatusCode: statusCode,
	}
}

func NewValidation(fields []FieldError) Error {
	return Error{
		Message:    "Input tidak valid",
		StatusCode: 400,
		ErrorCode:  "VALIDATION_ERROR",
		Fields:     &fields,
	}
}

// ForField keeps an existing public error message/status and adds the exact
// form field that caused it. Non-public errors are left untouched so internal
// failures are not accidentally exposed as validation errors.
func ForField(err error, field, label string) error {
	var myErr Error
	if !errors.As(err, &myErr) {
		return err
	}

	myErr.ErrorCode = "VALIDATION_ERROR"
	fields := []FieldError{{
		Field:   field,
		Label:   label,
		Message: myErr.Message,
	}}
	myErr.Fields = &fields
	return myErr
}

// digunakan untuk mengembalikan error berdasarkan mode aplikasi sekarang
// jika mode aplikasi adalah development, maka akan mengembalikan error asli
// jika mode aplikasi adalah production, maka akan mengembalikan error yang sudah di wrap
// bedanya ini digunakan untuk deklarasi awal error
func NewWrap(deverror error, proderror error, statusCode int) Error {
	mode := os.Getenv("APP_MODE")
	switch mode {
	case "development":
		return Error{
			Message:    deverror.Error(),
			StatusCode: statusCode,
		}
	default:
		return Error{
			Message:    proderror.Error(),
			StatusCode: statusCode,
		}
	}
}

// digunakan untuk mengembalikan error berdasarkan mode aplikasi sekarang
// jika mode aplikasi adalah development, maka akan mengembalikan error asli
// jika mode aplikasi adalah production, maka akan mengembalikan error yang sudah di wrap
func Wrap(actErr error, prodErr error) error {
	mode := os.Getenv("APP_MODE")
	switch mode {
	case "development":
		return actErr
	default:
		return prodErr
	}
}

func (e Error) Error() string {
	return e.Message
}
