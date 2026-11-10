package kitab_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"kitab-dapur/api/internal/kitab"
)

type fakeRepo struct {
	names []string
	err   error
}

func (f *fakeRepo) Create(_ context.Context, name string) (kitab.Kitab, error) {
	if f.err != nil {
		return kitab.Kitab{}, f.err
	}
	f.names = append(f.names, name)
	return kitab.Kitab{ID: "k1", Name: name}, nil
}

func TestServiceCreate(t *testing.T) {
	errDB := errors.New("db down")
	tests := []struct {
		name      string
		input     string
		repoErr   error
		wantName  string
		wantField string
		wantErr   error
	}{
		{name: "trims the name", input: "  Keluarga Syafrin ", wantName: "Keluarga Syafrin"},
		{name: "accepts the max length in characters", input: strings.Repeat("é", kitab.NameMaxLength), wantName: strings.Repeat("é", kitab.NameMaxLength)},
		{name: "rejects an empty name", input: "   ", wantField: "name"},
		{name: "rejects a name that is too long", input: strings.Repeat("a", kitab.NameMaxLength+1), wantField: "name"},
		{name: "wraps a repository error", input: "Keluarga", repoErr: errDB, wantErr: errDB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{err: tt.repoErr}
			got, err := kitab.NewService(repo).Create(context.Background(), tt.input)

			var vErr *kitab.ValidationError
			switch {
			case tt.wantField != "":
				if !errors.As(err, &vErr) || vErr.Field != tt.wantField {
					t.Fatalf("err = %v, want validation error on %q", err, tt.wantField)
				}
				if len(repo.names) != 0 {
					t.Fatal("invalid kitab was stored")
				}
			case tt.wantErr != nil:
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
			default:
				if err != nil || got.Name != tt.wantName {
					t.Fatalf("got %+v, %v; want name %q", got, err, tt.wantName)
				}
			}
		})
	}
}

func TestValidationErrorMessage(t *testing.T) {
	err := &kitab.ValidationError{Field: "name", Message: "wajib"}
	if err.Error() != "name: wajib" {
		t.Fatalf("Error() = %q", err.Error())
	}
}
