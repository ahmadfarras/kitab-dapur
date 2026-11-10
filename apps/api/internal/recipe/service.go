package recipe

import (
	"context"
	"fmt"
)

// Page sizes for List.
const (
	DefaultListLimit = 50
	MaxListLimit     = 100
)

// Repository stores recipes. Every method is scoped to one kitab.
type Repository interface {
	// Create returns ErrKitabNotFound when the kitab does not exist.
	Create(ctx context.Context, d Draft) (Recipe, error)
	// List returns the newest recipes first.
	List(ctx context.Context, kitabID string, limit int) ([]Recipe, error)
	// Get returns ErrNotFound when the recipe is not in that kitab.
	Get(ctx context.Context, kitabID, id string) (Recipe, error)
}

// Service holds the recipe use cases.
type Service struct {
	repo Repository
}

// NewService returns a Service that stores recipes in repo.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Add validates d and stores it.
func (s *Service) Add(ctx context.Context, d Draft) (Recipe, error) {
	d, err := d.Normalize()
	if err != nil {
		return Recipe{}, err
	}
	r, err := s.repo.Create(ctx, d)
	if err != nil {
		return Recipe{}, fmt.Errorf("add recipe: %w", err)
	}
	return r, nil
}

// List returns one page of a kitab's recipes, newest first. limit is clamped to 1..MaxListLimit; 0 means the default.
func (s *Service) List(ctx context.Context, kitabID string, limit int) ([]Recipe, error) {
	if limit <= 0 {
		limit = DefaultListLimit
	}
	limit = min(limit, MaxListLimit)
	recipes, err := s.repo.List(ctx, kitabID, limit)
	if err != nil {
		return nil, fmt.Errorf("list recipes: %w", err)
	}
	return recipes, nil
}

// Get returns one recipe of a kitab, or ErrNotFound.
func (s *Service) Get(ctx context.Context, kitabID, id string) (Recipe, error) {
	r, err := s.repo.Get(ctx, kitabID, id)
	if err != nil {
		return Recipe{}, fmt.Errorf("get recipe: %w", err)
	}
	return r, nil
}
