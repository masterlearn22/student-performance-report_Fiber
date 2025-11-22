package repository

import (
	"database/sql"

	models "student-performance-report/app/models/postgresql"

	"github.com/google/uuid"
)

type StudentRepository interface {
	GetAllStudents() ([]models.Student, error)
	GetStudentByID(id uuid.UUID) (*models.Student, error)
	UpdateAdvisor(studentID, lecturerID uuid.UUID) error
}

type studentRepository struct {
	pg  *sql.DB
}

func NewStudentRepository(pg *sql.DB) StudentRepository {
	return &studentRepository{pg: pg}
}

func (r *studentRepository) GetAllStudents() ([]models.Student, error) {
	rows, err := r.pg.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Student
	for rows.Next() {
		var s models.Student
		err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
			&s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *studentRepository) GetStudentByID(id uuid.UUID) (*models.Student, error) {
	var s models.Student
	err := r.pg.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students WHERE id=$1
	`, id).Scan(
		&s.ID, &s.UserID, &s.StudentID,
		&s.ProgramStudy, &s.AcademicYear,
		&s.AdvisorID, &s.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	return &s, err
}

func (r *studentRepository) UpdateAdvisor(studentID, lecturerID uuid.UUID) error {
	_, err := r.pg.Exec(`
		UPDATE students SET advisor_id=$1 WHERE id=$2
	`, lecturerID, studentID)
	return err
}

