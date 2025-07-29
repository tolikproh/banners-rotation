package cnst

import "errors"

var (
	ErrGetBanner      = errors.New("error get banner")
	ErrObjectNotFound = errors.New("object not found")
	ErrGetRowsError   = errors.New("get rows error")
	ErrStatusCode     = errors.New("error status code not OK (200)")
)
