package models

import (
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	maxFieldLength = 256
)

// Attribute is a key-value pair consisting of variable length strings
// upto `maxFieldLength` chars in length
type Attribute struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid"`
	Key        string
	Value      string
	DocumentID uuid.UUID
}

func ValidateAttribute(key, value string) error {
	lk := len(key)
	lv := len(value)
	if lk == 0 || lk > maxFieldLength {
		return errors.New("invalid key length")
	}

	if lv == 0 || lv > maxFieldLength {
		return errors.New("invalid value length")
	}

	return nil
}

// NewAttribute returns a new `Attribute` struct after validation checks
func NewAttribute(key, value string) (*Attribute, error) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	if err := ValidateAttribute(key, value); err != nil {
		return nil, err
	}

	return &Attribute{
		ID:    uuid.New(),
		Key:   key,
		Value: value,
	}, nil
}
