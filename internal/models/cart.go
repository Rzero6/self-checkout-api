package models

import "github.com/google/uuid"

type Cart struct {
	ID        uint64    `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	Status    string    `json:"status"`
}
