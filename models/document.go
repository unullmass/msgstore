package models

import (
	"time"

	"github.com/google/uuid"
)

// Document is a unique entity in a time-series datastore
// that comprises one-or more key-value pairs along with temporal data
type Document struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid"`
	Attributes []Attribute
	Timestamp  time.Time `gorm:"index:idx_ts"`
}
