// Package kitab is a family's recipe collection: its rules and use cases.
package kitab

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// NameMaxLength is the longest kitab name, in characters.
const NameMaxLength = 80

// Kitab is one family's recipe collection.
type Kitab struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

// ValidationError reports input that breaks a kitab rule.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string { return e.Field + ": " + e.Message }

// Repository stores kitabs.
type Repository interface {
	Create(ctx context.Context, name string) (Kitab, error)
}

// Service holds the kitab use cases.
type Service struct {
	repo Repository
}

// NewService returns a Service that stores kitabs in repo.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create validates name and stores a new kitab.
func (s *Service) Create(ctx context.Context, name string) (Kitab, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Kitab{}, &ValidationError{Field: "name", Message: "Nama kitab wajib diisi"}
	}
	if utf8.RuneCountInString(name) > NameMaxLength {
		return Kitab{}, &ValidationError{Field: "name", Message: fmt.Sprintf("Maksimal %d karakter", NameMaxLength)}
	}
	k, err := s.repo.Create(ctx, name)
	if err != nil {
		return Kitab{}, fmt.Errorf("create kitab: %w", err)
	}
	return k, nil
}
