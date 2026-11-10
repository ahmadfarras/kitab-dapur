package repository_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	kitabrepo "kitab-dapur/api/internal/kitab/repository"
	"kitab-dapur/api/internal/platform/pgtest"
	"kitab-dapur/api/internal/recipe"
	"kitab-dapur/api/internal/recipe/repository"
)

const unknownID = "00000000-0000-0000-0000-000000000000"

func newKitab(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	k, err := kitabrepo.NewPostgres(db).Create(context.Background(), "Keluarga")
	if err != nil {
		t.Fatal(err)
	}
	return k.ID
}

func draft(kitabID, title string) recipe.Draft {
	return recipe.Draft{KitabID: kitabID, Title: title, Author: "Ibu", Ingredients: []string{"garam", "gula"}, Steps: []string{"Aduk"}}
}

func TestPostgresRoundTrip(t *testing.T) {
	ctx := context.Background()
	db := pgtest.New(t)
	repo := repository.NewPostgres(db)
	kitabID := newKitab(t, db)

	created, err := repo.Create(ctx, draft(kitabID, "Opor"))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	got, err := repo.Get(ctx, kitabID, created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !reflect.DeepEqual(got, created) || got.Ingredients[1] != "gula" {
		t.Fatalf("Get() = %+v, want %+v", got, created)
	}
}

func TestPostgresListNewestFirstWithinLimitAndKitab(t *testing.T) {
	ctx := context.Background()
	db := pgtest.New(t)
	repo := repository.NewPostgres(db)
	kitabID, otherKitabID := newKitab(t, db), newKitab(t, db)
	for _, title := range []string{"pertama", "kedua", "ketiga"} {
		if _, err := repo.Create(ctx, draft(kitabID, title)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := repo.Create(ctx, draft(otherKitabID, "milik keluarga lain")); err != nil {
		t.Fatal(err)
	}

	got, err := repo.List(ctx, kitabID, 2)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 2 || got[0].Title != "ketiga" || got[1].Title != "kedua" {
		t.Fatalf("List() titles = %v", titles(got))
	}
	if empty, err := repo.List(ctx, unknownID, 10); err != nil || len(empty) != 0 {
		t.Fatalf("List(unknown kitab) = %v, %v", titles(empty), err)
	}
}

func TestPostgresErrors(t *testing.T) {
	ctx := context.Background()
	db := pgtest.New(t)
	repo := repository.NewPostgres(db)
	kitabID, otherKitabID := newKitab(t, db), newKitab(t, db)
	created, err := repo.Create(ctx, draft(kitabID, "Rawon"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := repo.Create(ctx, draft(unknownID, "x")); !errors.Is(err, recipe.ErrKitabNotFound) {
		t.Errorf("Create(unknown kitab) err = %v, want ErrKitabNotFound", err)
	}
	if _, err := repo.Get(ctx, kitabID, unknownID); !errors.Is(err, recipe.ErrNotFound) {
		t.Errorf("Get(unknown id) err = %v, want ErrNotFound", err)
	}
	if _, err := repo.Get(ctx, otherKitabID, created.ID); !errors.Is(err, recipe.ErrNotFound) {
		t.Errorf("Get(recipe of another kitab) err = %v, want ErrNotFound", err)
	}
}

func titles(rs []recipe.Recipe) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.Title)
	}
	return out
}
