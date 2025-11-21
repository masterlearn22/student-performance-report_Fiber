package models

import (
	"time"
	"github.com/google/uuid"
)

type Lecturer struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"userId" db:"user_id"`
	LecturerID string    `json:"lecturerId" db:"lecturer_id"` // NIP/NIDN
	Department string    `json:"department" db:"department"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}