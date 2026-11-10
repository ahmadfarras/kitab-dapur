package recipe_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"kitab-dapur/api/internal/recipe"
)

func validDraft() recipe.Draft {
	return recipe.Draft{
		KitabID:     "k1",
		Title:       "  Rendang Padang ",
		Author:      " Ibu ",
		Ingredients: []string{"1 kg daging sapi", "  ", " santan "},
		Steps:       []string{"Masak santan", ""},
	}
}

func TestDraftNormalize(t *testing.T) {
	got, err := validDraft().Normalize()
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	want := recipe.Draft{
		KitabID:     "k1",
		Title:       "Rendang Padang",
		Author:      "Ibu",
		Ingredients: []string{"1 kg daging sapi", "santan"},
		Steps:       []string{"Masak santan"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Normalize() = %+v, want %+v", got, want)
	}
}

func TestDraftNormalizeRejects(t *testing.T) {
	long := func(n int) string { return strings.Repeat("a", n) }
	many := func(n int) []string { return strings.Split(strings.Repeat("garam,", n-1)+"garam", ",") }
	tests := []struct {
		name  string
		edit  func(*recipe.Draft)
		field string
	}{
		{"empty title", func(d *recipe.Draft) { d.Title = " " }, "title"},
		{"title too long", func(d *recipe.Draft) { d.Title = long(recipe.TitleMaxLength + 1) }, "title"},
		{"empty author", func(d *recipe.Draft) { d.Author = "" }, "author"},
		{"author too long", func(d *recipe.Draft) { d.Author = long(recipe.AuthorMaxLength + 1) }, "author"},
		{"no ingredients", func(d *recipe.Draft) { d.Ingredients = []string{" ", ""} }, "ingredients"},
		{"nil ingredients", func(d *recipe.Draft) { d.Ingredients = nil }, "ingredients"},
		{"too many ingredients", func(d *recipe.Draft) { d.Ingredients = many(recipe.ItemsMaxCount + 1) }, "ingredients"},
		{"ingredient too long", func(d *recipe.Draft) { d.Ingredients = []string{long(recipe.ItemMaxLength + 1)} }, "ingredients"},
		{"no steps", func(d *recipe.Draft) { d.Steps = nil }, "steps"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := validDraft()
			tt.edit(&d)
			_, err := d.Normalize()
			var vErr *recipe.ValidationError
			if !errors.As(err, &vErr) || vErr.Field != tt.field {
				t.Fatalf("err = %v, want validation error on %q", err, tt.field)
			}
			if vErr.Error() == "" {
				t.Fatal("empty error message")
			}
		})
	}
}

func TestDraftNormalizeAcceptsLimits(t *testing.T) {
	d := validDraft()
	d.Title = strings.Repeat("é", recipe.TitleMaxLength)
	d.Ingredients = make([]string, recipe.ItemsMaxCount)
	for i := range d.Ingredients {
		d.Ingredients[i] = strings.Repeat("a", recipe.ItemMaxLength)
	}
	if _, err := d.Normalize(); err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
}

type fakeRepo struct {
	created   []recipe.Draft
	limits    []int
	recipes   map[string]recipe.Recipe
	createErr error
}

func (f *fakeRepo) Create(_ context.Context, d recipe.Draft) (recipe.Recipe, error) {
	if f.createErr != nil {
		return recipe.Recipe{}, f.createErr
	}
	f.created = append(f.created, d)
	return recipe.Recipe{ID: "r1", KitabID: d.KitabID, Title: d.Title}, nil
}

func (f *fakeRepo) List(_ context.Context, _ string, limit int) ([]recipe.Recipe, error) {
	f.limits = append(f.limits, limit)
	return nil, nil
}

func (f *fakeRepo) Get(_ context.Context, kitabID, id string) (recipe.Recipe, error) {
	r, ok := f.recipes[kitabID+"/"+id]
	if !ok {
		return recipe.Recipe{}, recipe.ErrNotFound
	}
	return r, nil
}

func TestServiceAdd(t *testing.T) {
	t.Run("stores the normalized draft", func(t *testing.T) {
		repo := &fakeRepo{}
		got, err := recipe.NewService(repo).Add(context.Background(), validDraft())
		if err != nil || got.Title != "Rendang Padang" || len(repo.created) != 1 {
			t.Fatalf("Add() = %+v, %v; created %d", got, err, len(repo.created))
		}
	})
	t.Run("stores nothing when invalid", func(t *testing.T) {
		repo := &fakeRepo{}
		d := validDraft()
		d.Title = ""
		if _, err := recipe.NewService(repo).Add(context.Background(), d); err == nil || len(repo.created) != 0 {
			t.Fatalf("err = %v, created %d", err, len(repo.created))
		}
	})
	t.Run("keeps ErrKitabNotFound visible", func(t *testing.T) {
		repo := &fakeRepo{createErr: recipe.ErrKitabNotFound}
		if _, err := recipe.NewService(repo).Add(context.Background(), validDraft()); !errors.Is(err, recipe.ErrKitabNotFound) {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestServiceListClampsLimit(t *testing.T) {
	tests := []struct{ in, want int }{
		{0, recipe.DefaultListLimit},
		{-5, recipe.DefaultListLimit},
		{1, 1},
		{recipe.MaxListLimit, recipe.MaxListLimit},
		{recipe.MaxListLimit + 1, recipe.MaxListLimit},
	}
	for _, tt := range tests {
		repo := &fakeRepo{}
		if _, err := recipe.NewService(repo).List(context.Background(), "k1", tt.in); err != nil {
			t.Fatal(err)
		}
		if repo.limits[0] != tt.want {
			t.Errorf("List(limit=%d) asked repo for %d, want %d", tt.in, repo.limits[0], tt.want)
		}
	}
}

func TestServiceGet(t *testing.T) {
	repo := &fakeRepo{recipes: map[string]recipe.Recipe{"k1/r1": {ID: "r1"}}}
	svc := recipe.NewService(repo)
	if got, err := svc.Get(context.Background(), "k1", "r1"); err != nil || got.ID != "r1" {
		t.Fatalf("Get() = %+v, %v", got, err)
	}
	if _, err := svc.Get(context.Background(), "k2", "r1"); !errors.Is(err, recipe.ErrNotFound) {
		t.Fatalf("Get() from another kitab err = %v, want ErrNotFound", err)
	}
}
