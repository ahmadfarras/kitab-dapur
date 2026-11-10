// Package repository stores kitabs in PostgreSQL.
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"kitab-dapur/api/internal/kitab"
)

// Postgres implements kitab.Repository.
type Postgres struct {
	db *pgxpool.Pool
}

// NewPostgres returns a Postgres repository on db.
func NewPostgres(db *pgxpool.Pool) *Postgres {
	return &Postgres{db: db}
}

// Create inserts a kitab and returns it with its generated id.
func (p *Postgres) Create(ctx context.Context, name string) (kitab.Kitab, error) {
	k := kitab.Kitab{Name: name}
	err := p.db.QueryRow(ctx, `INSERT INTO kitabs (name) VALUES ($1) RETURNING id::text, created_at`, name).
		Scan(&k.ID, &k.CreatedAt)
	if err != nil {
		return kitab.Kitab{}, fmt.Errorf("insert kitab: %w", err)
	}
	return k, nil
}
