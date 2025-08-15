package repoutil

import (
	appErr "NeliQuiz/internal/shared/errorx"
	"errors"
	gormErr "gorm.io/gorm"
)

func TranslateGormError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gormErr.ErrRecordNotFound):
		return appErr.NotFound("data tidak ditemukan")
	case errors.Is(err, gormErr.ErrMissingWhereClause):
		return appErr.BadRequest("missing WHERE clause dalam query")
	case errors.Is(err, gormErr.ErrInvalidData):
		return appErr.BadRequest("data tidak valid")
	case errors.Is(err, gormErr.ErrInvalidTransaction):
		return appErr.InternalError(err)
	default:
		return appErr.InternalError(err)
	}
}
