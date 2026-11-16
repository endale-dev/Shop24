package interfaces

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrConflict          = errors.New("conflict")
)

type DuplicateKey struct {
	Field string
}

func (e *DuplicateKey) Error() string { return "duplicate key" }
