package repository_test

import (
	"context"
	"testing"

	"kitab-dapur/api/internal/kitab/repository"
	"kitab-dapur/api/internal/platform/pgtest"
)

func TestPostgresCreate(t *testing.T) {
	db := pgtest.New(t)
	k, err := repository.NewPostgres(db).Create(context.Background(), "Keluarga Syafrin")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if k.ID == "" || k.Name != "Keluarga Syafrin" || k.CreatedAt.IsZero() {
		t.Fatalf("Create() = %+v", k)
	}

	var name string
	if err := db.QueryRow(context.Background(), `SELECT name FROM kitabs WHERE id = $1`, k.ID).Scan(&name); err != nil || name != k.Name {
		t.Fatalf("stored row = %q, %v", name, err)
	}
}
