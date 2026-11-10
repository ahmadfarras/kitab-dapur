package handler_test

import (
	"context"
	"errors"
	"testing"

	"kitab-dapur/api/internal/kitab"
	"kitab-dapur/api/internal/kitab/handler"
	"kitab-dapur/api/internal/openapi"
)

type fakeService struct {
	kitab kitab.Kitab
	err   error
}

func (f fakeService) Create(context.Context, string) (kitab.Kitab, error) { return f.kitab, f.err }

func request(name string) openapi.CreateKitabRequestObject {
	return openapi.CreateKitabRequestObject{Body: &openapi.NewKitab{Name: name}}
}

func TestCreateKitab(t *testing.T) {
	errDB := errors.New("db down")
	tests := []struct {
		name    string
		svc     fakeService
		want    any
		wantErr error
	}{
		{name: "created", svc: fakeService{kitab: kitab.Kitab{ID: "7c9e6679-7425-40de-944b-e07fc1f90ae7", Name: "Keluarga"}}, want: openapi.CreateKitab201JSONResponse{}},
		{name: "invalid input", svc: fakeService{err: &kitab.ValidationError{Field: "name", Message: "wajib"}}, want: openapi.CreateKitab400JSONResponse{}},
		{name: "unexpected error goes to the 500 handler", svc: fakeService{err: errDB}, wantErr: errDB},
		{name: "corrupt id is an error, not a response", svc: fakeService{kitab: kitab.Kitab{ID: "nope"}}, wantErr: errors.New("any")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.NewKitabs(tt.svc).CreateKitab(context.Background(), request("Keluarga"))
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("got %T, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			switch tt.want.(type) {
			case openapi.CreateKitab201JSONResponse:
				if r, ok := got.(openapi.CreateKitab201JSONResponse); !ok || r.Name != "Keluarga" {
					t.Fatalf("got %#v", got)
				}
			case openapi.CreateKitab400JSONResponse:
				if r, ok := got.(openapi.CreateKitab400JSONResponse); !ok || *r.Field != "name" || r.Message != "wajib" {
					t.Fatalf("got %#v", got)
				}
			}
		})
	}
}
