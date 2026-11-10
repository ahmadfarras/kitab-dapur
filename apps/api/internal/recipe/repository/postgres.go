// Package repository stores recipes in PostgreSQL.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"kitab-dapur/api/internal/recipe"
)

const foreignKeyViolation = "23503"

const columns = `id::text, kitab_id::text, title, author, ingredients, steps, created_at`

// Postgres implements recipe.Repository.
type Postgres struct {
	db *pgxpool.Pool
}

// NewPostgres returns a Postgres repository on db.
func NewPostgres(db *pgxpool.Pool) *Postgres {
	return &Postgres{db: db}
}

// Create inserts d, or returns recipe.ErrKitabNotFound when its kitab does not exist.
func (p *Postgres) Create(ctx context.Context, d recipe.Draft) (recipe.Recipe, error) {
	row := p.db.QueryRow(ctx,
		`INSERT INTO recipes (kitab_id, title, author, ingredients, steps) VALUES ($1, $2, $3, $4, $5) RETURNING `+columns,
		d.KitabID, d.Title, d.Author, d.Ingredients, d.Steps)
	r, err := scan(row)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation {
		return recipe.Recipe{}, recipe.ErrKitabNotFound
	}
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("insert recipe: %w", err)
	}
	return r, nil
}

// List returns up to limit recipes of a kitab, newest first.
func (p *Postgres) List(ctx context.Context, kitabID string, limit int) ([]recipe.Recipe, error) {
	rows, err := p.db.Query(ctx,
		`SELECT `+columns+` FROM recipes WHERE kitab_id = $1 ORDER BY created_at DESC, id DESC LIMIT $2`,
		kitabID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	recipes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (recipe.Recipe, error) { return scan(row) })
	if err != nil {
		return nil, fmt.Errorf("scan recipes: %w", err)
	}
	return recipes, nil
}

// Get returns a recipe of a kitab, or recipe.ErrNotFound.
func (p *Postgres) Get(ctx context.Context, kitabID, id string) (recipe.Recipe, error) {
	r, err := scan(p.db.QueryRow(ctx, `SELECT `+columns+` FROM recipes WHERE kitab_id = $1 AND id = $2`, kitabID, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return recipe.Recipe{}, recipe.ErrNotFound
	}
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("select recipe: %w", err)
	}
	return r, nil
}

func scan(row pgx.Row) (recipe.Recipe, error) {
	var r recipe.Recipe
	err := row.Scan(&r.ID, &r.KitabID, &r.Title, &r.Author, &r.Ingredients, &r.Steps, &r.CreatedAt)
	return r, err
}
