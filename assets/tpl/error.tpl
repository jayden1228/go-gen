package {{.PackageName}}

import (
	"errors"
)

var (
	ErrInsertFailed = errors.New("ErrInsertFailed")
	ErrNotFound     = errors.New("ErrNotFound")
	ErrDeleteFailed = errors.New("ErrDeleteFailed")
	ErrUpdateFailed = errors.New("ErrUpdateFailed")
)
