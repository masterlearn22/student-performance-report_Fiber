package models

import (
	"time"
	"github.com/google/uuid"
)

type Student struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"userId" db:"user_id"`
	StudentID    string     `json:"studentId" db:"student_id"` // NIM
	ProgramStudy string     `json:"programStudy" db:"program_study"`
	AcademicYear string     `json:"academicYear" db:"academic_year"`
	AdvisorID    *uuid.UUID `json:"advisorId" db:"advisor_id"` // Pointer karena bisa NULL
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`

	// Helper field untuk join user
	FullName     string     `json:"fullName,omitempty" db:"full_name"`
}