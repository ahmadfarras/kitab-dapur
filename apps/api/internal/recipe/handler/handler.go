// Package handler serves the recipe operations of the OpenAPI contract.
package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"kitab-dapur/api/internal/openapi"
	"kitab-dapur/api/internal/recipe"
)

type service interface {
	Add(ctx context.Context, d recipe.Draft) (recipe.Recipe, error)
	List(ctx context.Context, kitabID string, limit int) ([]recipe.Recipe, error)
	Get(ctx context.Context, kitabID, id string) (recipe.Recipe, error)
}

// Recipes implements the recipe part of openapi.StrictServerInterface.
type Recipes struct {
	svc service
}

// NewRecipes returns a Recipes on svc.
func NewRecipes(svc service) *Recipes {
	return &Recipes{svc: svc}
}

var kitabNotFound = openapi.NotFoundJSONResponse{Message: "Kitab tidak ditemukan"}

// CreateRecipe handles POST /kitabs/{kitabId}/recipes.
func (h *Recipes) CreateRecipe(ctx context.Context, req openapi.CreateRecipeRequestObject) (openapi.CreateRecipeResponseObject, error) {
	r, err := h.svc.Add(ctx, recipe.Draft{
		KitabID:     req.KitabId.String(),
		Title:       req.Body.Title,
		Author:      req.Body.Author,
		Ingredients: req.Body.Ingredients,
		Steps:       req.Body.Steps,
	})
	var vErr *recipe.ValidationError
	switch {
	case errors.As(err, &vErr):
		return openapi.CreateRecipe400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{Message: vErr.Message, Field: &vErr.Field}}, nil
	case errors.Is(err, recipe.ErrKitabNotFound):
		return openapi.CreateRecipe404JSONResponse{NotFoundJSONResponse: kitabNotFound}, nil
	case err != nil:
		return nil, err
	}
	body, err := toAPI(r)
	if err != nil {
		return nil, err
	}
	return openapi.CreateRecipe201JSONResponse(body), nil
}

// ListRecipes handles GET /kitabs/{kitabId}/recipes.
func (h *Recipes) ListRecipes(ctx context.Context, req openapi.ListRecipesRequestObject) (openapi.ListRecipesResponseObject, error) {
	limit := 0
	if req.Params.Limit != nil {
		limit = *req.Params.Limit
	}
	recipes, err := h.svc.List(ctx, req.KitabId.String(), limit)
	if err != nil {
		return nil, err
	}
	items := make([]openapi.Recipe, 0, len(recipes))
	for _, r := range recipes {
		item, err := toAPI(r)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return openapi.ListRecipes200JSONResponse{Items: items}, nil
}

// GetRecipe handles GET /kitabs/{kitabId}/recipes/{recipeId}.
func (h *Recipes) GetRecipe(ctx context.Context, req openapi.GetRecipeRequestObject) (openapi.GetRecipeResponseObject, error) {
	r, err := h.svc.Get(ctx, req.KitabId.String(), req.RecipeId.String())
	if errors.Is(err, recipe.ErrNotFound) {
		return openapi.GetRecipe404JSONResponse{NotFoundJSONResponse: openapi.NotFoundJSONResponse{Message: "Resep tidak ditemukan"}}, nil
	}
	if err != nil {
		return nil, err
	}
	body, err := toAPI(r)
	if err != nil {
		return nil, err
	}
	return openapi.GetRecipe200JSONResponse(body), nil
}

func toAPI(r recipe.Recipe) (openapi.Recipe, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return openapi.Recipe{}, fmt.Errorf("recipe id %q: %w", r.ID, err)
	}
	kitabID, err := uuid.Parse(r.KitabID)
	if err != nil {
		return openapi.Recipe{}, fmt.Errorf("recipe %s kitab id %q: %w", r.ID, r.KitabID, err)
	}
	return openapi.Recipe{
		Id:          id,
		KitabId:     kitabID,
		Title:       r.Title,
		Author:      r.Author,
		Ingredients: r.Ingredients,
		Steps:       r.Steps,
		CreatedAt:   r.CreatedAt,
	}, nil
}
