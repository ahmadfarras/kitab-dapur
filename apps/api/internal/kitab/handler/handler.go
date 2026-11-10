// Package handler serves the kitab operations of the OpenAPI contract.
package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"kitab-dapur/api/internal/kitab"
	"kitab-dapur/api/internal/openapi"
)

type service interface {
	Create(ctx context.Context, name string) (kitab.Kitab, error)
}

// Kitabs implements the kitab part of openapi.StrictServerInterface.
type Kitabs struct {
	svc service
}

// NewKitabs returns a Kitabs on svc.
func NewKitabs(svc service) *Kitabs {
	return &Kitabs{svc: svc}
}

// CreateKitab handles POST /kitabs.
func (h *Kitabs) CreateKitab(ctx context.Context, req openapi.CreateKitabRequestObject) (openapi.CreateKitabResponseObject, error) {
	k, err := h.svc.Create(ctx, req.Body.Name)
	var vErr *kitab.ValidationError
	if errors.As(err, &vErr) {
		return openapi.CreateKitab400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{Message: vErr.Message, Field: &vErr.Field}}, nil
	}
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(k.ID)
	if err != nil {
		return nil, fmt.Errorf("kitab id %q: %w", k.ID, err)
	}
	return openapi.CreateKitab201JSONResponse{Id: id, Name: k.Name, CreatedAt: k.CreatedAt}, nil
}
