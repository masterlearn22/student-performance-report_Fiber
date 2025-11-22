package models

import (
    "time"
    "github.com/google/uuid"
)

type Student struct {
    ID            uuid.UUID  `json:"id" db:"id"`
    UserID        uuid.UUID  `json:"user_id" db:"user_id"`
    StudentID     string     `json:"student_id" db:"student_id"` // NIM
    ProgramStudy  string     `json:"program_study" db:"program_study"`
    AcademicYear  string     `json:"academic_year" db:"academic_year"`
    AdvisorID     *uuid.UUID `json:"advisor_id" db:"advisor_id"` // nullable FK
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type StudentResp struct {
    ID           uuid.UUID  `json:"id"`
    StudentID    string     `json:"studentId"`
    FullName     string     `json:"fullName,omitempty"`
    ProgramStudy string     `json:"programStudy"`
    AcademicYear string     `json:"academicYear"`
    AdvisorID    *uuid.UUID `json:"advisorId,omitempty"`
}

