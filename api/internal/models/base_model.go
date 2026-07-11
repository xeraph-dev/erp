package models

import (
	"time"

	"github.com/google/uuid"
)

type base struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedByID uuid.UUID `json:"created_by_id"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedByID uuid.UUID `json:"updated_by_id"`
	DeletedAt   time.Time `json:"deleted_at"`
	DeletedByID uuid.UUID `json:"deleted_by_id"`
}
