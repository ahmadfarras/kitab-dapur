// Package recipe holds the recipe rules and use cases.
package recipe

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Limits on what one recipe may hold, in characters and lines.
const (
	TitleMaxLength  = 120
	AuthorMaxLength = 60
	ItemMaxLength   = 500
	ItemsMaxCount   = 100
)

var (
	// ErrNotFound means the recipe does not exist in that kitab.
	ErrNotFound = errors.New("recipe not found")
	// ErrKitabNotFound means the kitab the recipe belongs to does not exist.
	ErrKitabNotFound = errors.New("kitab not found")
)

// Recipe is a stored recipe.
type Recipe struct {
	ID      string
	KitabID string
	Title   string
	// Author is who the recipe comes from, e.g. "Ibu".
	Author      string
	Ingredients []string
	Steps       []string
	CreatedAt   time.Time
}

// Draft is a recipe that has not been stored yet.
type Draft struct {
	KitabID     string
	Title       string
	Author      string
	Ingredients []string
	Steps       []string
}

// ValidationError reports input that breaks a recipe rule.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string { return e.Field + ": " + e.Message }

// Normalize trims every text, drops blank lines and checks the recipe rules.
func (d Draft) Normalize() (Draft, error) {
	var err error
	if d.Title, err = requireText("title", d.Title, TitleMaxLength, "Judul resep wajib diisi"); err != nil {
		return Draft{}, err
	}
	if d.Author, err = requireText("author", d.Author, AuthorMaxLength, "Nama pemilik resep wajib diisi"); err != nil {
		return Draft{}, err
	}
	if d.Ingredients, err = requireItems("ingredients", d.Ingredients, "Isi minimal satu bahan"); err != nil {
		return Draft{}, err
	}
	if d.Steps, err = requireItems("steps", d.Steps, "Isi minimal satu langkah"); err != nil {
		return Draft{}, err
	}
	return d, nil
}

func requireText(field, value string, maxLength int, emptyMessage string) (string, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return "", &ValidationError{Field: field, Message: emptyMessage}
	}
	if utf8.RuneCountInString(text) > maxLength {
		return "", &ValidationError{Field: field, Message: fmt.Sprintf("Maksimal %d karakter", maxLength)}
	}
	return text, nil
}

func requireItems(field string, values []string, emptyMessage string) ([]string, error) {
	if len(values) > ItemsMaxCount {
		return nil, &ValidationError{Field: field, Message: fmt.Sprintf("Maksimal %d baris", ItemsMaxCount)}
	}
	items := make([]string, 0, len(values))
	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" {
			continue
		}
		if utf8.RuneCountInString(item) > ItemMaxLength {
			return nil, &ValidationError{Field: field, Message: fmt.Sprintf("Tiap baris maksimal %d karakter", ItemMaxLength)}
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, &ValidationError{Field: field, Message: emptyMessage}
	}
	return items, nil
}
