package repository

import (
	"database/sql"

	models "student-performance-report/app/models/postgresql"
	"github.com/google/uuid"
)

type LecturerRepository interface {
	GetAllLecturers() ([]models.Lecturer, error)
	GetLecturerByID(id uuid.UUID) (*models.Lecturer, error)
	GetAdvisees(lecturerID uuid.UUID) ([]models.Student, error)
}

type lecturerRepository struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepository{db: db}
}

func (r *lecturerRepository) GetAllLecturers() ([]models.Lecturer, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, lecturer_id, department, created_at 
		FROM lecturers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Lecturer
	for rows.Next() {
		var l models.Lecturer
		rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
		list = append(list, l)
	}
	return list, nil
}


func (r *lecturerRepository) GetLecturerByID(id uuid.UUID) (*models.Lecturer, error) {
	var l models.Lecturer
	err := r.db.QueryRow(`
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers WHERE id=$1
	`, id).Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}
	return &l, err
}


func (r *lecturerRepository) GetAdvisees(lecturerID uuid.UUID) ([]models.Student, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students WHERE advisor_id=$1
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Student
	for rows.Next() {
		var s models.Student
		rows.Scan(&s.ID, &s.UserID, &s.StudentID,
			&s.ProgramStudy, &s.AcademicYear,
			&s.AdvisorID, &s.CreatedAt,
		)
		list = append(list, s)
	}
	return list, nil
}


