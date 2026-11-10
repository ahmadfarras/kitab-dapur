package handler_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"kitab-dapur/api/internal/openapi"
	"kitab-dapur/api/internal/recipe"
	"kitab-dapur/api/internal/recipe/handler"
)

var (
	kitabID  = uuid.MustParse("7c9e6679-7425-40de-944b-e07fc1f90ae7")
	recipeID = uuid.MustParse("1b4e28ba-2fa1-11d2-883f-0016d3cca427")
	opor     = recipe.Recipe{ID: recipeID.String(), KitabID: kitabID.String(), Title: "Opor", Author: "Ibu", Ingredients: []string{"ayam"}, Steps: []string{"Masak"}}
	errDB    = errors.New("db down")
)

type fakeService struct {
	recipe    recipe.Recipe
	recipes   []recipe.Recipe
	err       error
	gotLimit  int
	gotDraft  recipe.Draft
	gotKitab  string
	gotRecipe string
}

func (f *fakeService) Add(_ context.Context, d recipe.Draft) (recipe.Recipe, error) {
	f.gotDraft = d
	return f.recipe, f.err
}

func (f *fakeService) List(_ context.Context, kitabID string, limit int) ([]recipe.Recipe, error) {
	f.gotKitab, f.gotLimit = kitabID, limit
	return f.recipes, f.err
}

func (f *fakeService) Get(_ context.Context, kitabID, id string) (recipe.Recipe, error) {
	f.gotKitab, f.gotRecipe = kitabID, id
	return f.recipe, f.err
}

func createRequest() openapi.CreateRecipeRequestObject {
	return openapi.CreateRecipeRequestObject{
		KitabId: kitabID,
		Body:    &openapi.NewRecipe{Title: "Opor", Author: "Ibu", Ingredients: []string{"ayam"}, Steps: []string{"Masak"}},
	}
}

func TestCreateRecipe(t *testing.T) {
	tests := []struct {
		name     string
		svc      *fakeService
		wantType any
		wantErr  bool
	}{
		{name: "created", svc: &fakeService{recipe: opor}, wantType: openapi.CreateRecipe201JSONResponse{}},
		{name: "invalid", svc: &fakeService{err: &recipe.ValidationError{Field: "title", Message: "wajib"}}, wantType: openapi.CreateRecipe400JSONResponse{}},
		{name: "unknown kitab", svc: &fakeService{err: recipe.ErrKitabNotFound}, wantType: openapi.CreateRecipe404JSONResponse{}},
		{name: "unexpected error", svc: &fakeService{err: errDB}, wantErr: true},
		{name: "corrupt id", svc: &fakeService{recipe: recipe.Recipe{ID: "nope"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.NewRecipes(tt.svc).CreateRecipe(context.Background(), createRequest())
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !sameType(got, tt.wantType) {
				t.Fatalf("got %T, want %T", got, tt.wantType)
			}
			if tt.svc.gotDraft.KitabID != kitabID.String() || tt.svc.gotDraft.Title != "Opor" {
				t.Fatalf("draft passed to service = %+v", tt.svc.gotDraft)
			}
		})
	}
}

func TestListRecipes(t *testing.T) {
	limit := 5
	tests := []struct {
		name      string
		limit     *int
		svc       *fakeService
		wantLimit int
		wantItems int
		wantErr   bool
	}{
		{name: "no limit lets the service pick", svc: &fakeService{recipes: []recipe.Recipe{opor, opor}}, wantLimit: 0, wantItems: 2},
		{name: "passes the limit", limit: &limit, svc: &fakeService{}, wantLimit: 5, wantItems: 0},
		{name: "unexpected error", svc: &fakeService{err: errDB}, wantErr: true},
		{name: "corrupt row", svc: &fakeService{recipes: []recipe.Recipe{{ID: "nope"}}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.NewRecipes(tt.svc).ListRecipes(context.Background(), openapi.ListRecipesRequestObject{KitabId: kitabID, Params: openapi.ListRecipesParams{Limit: tt.limit}})
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v", err)
			}
			if tt.wantErr {
				return
			}
			list, ok := got.(openapi.ListRecipes200JSONResponse)
			if !ok || len(list.Items) != tt.wantItems || list.Items == nil {
				t.Fatalf("got %#v, want %d items (never null)", got, tt.wantItems)
			}
			if tt.svc.gotLimit != tt.wantLimit || tt.svc.gotKitab != kitabID.String() {
				t.Fatalf("service got kitab %q limit %d", tt.svc.gotKitab, tt.svc.gotLimit)
			}
		})
	}
}

func TestGetRecipe(t *testing.T) {
	tests := []struct {
		name     string
		svc      *fakeService
		wantType any
		wantErr  bool
	}{
		{name: "found", svc: &fakeService{recipe: opor}, wantType: openapi.GetRecipe200JSONResponse{}},
		{name: "not found", svc: &fakeService{err: recipe.ErrNotFound}, wantType: openapi.GetRecipe404JSONResponse{}},
		{name: "unexpected error", svc: &fakeService{err: errDB}, wantErr: true},
		{name: "corrupt kitab id", svc: &fakeService{recipe: recipe.Recipe{ID: recipeID.String(), KitabID: "nope"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.NewRecipes(tt.svc).GetRecipe(context.Background(), openapi.GetRecipeRequestObject{KitabId: kitabID, RecipeId: recipeID})
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v", err)
			}
			if !tt.wantErr && !sameType(got, tt.wantType) {
				t.Fatalf("got %T, want %T", got, tt.wantType)
			}
			if tt.svc.gotKitab != kitabID.String() || tt.svc.gotRecipe != recipeID.String() {
				t.Fatalf("service got %q/%q", tt.svc.gotKitab, tt.svc.gotRecipe)
			}
		})
	}
}

func sameType(a, b any) bool {
	return a != nil && b != nil && typeName(a) == typeName(b)
}

func typeName(v any) string { return fmt.Sprintf("%T", v) }
