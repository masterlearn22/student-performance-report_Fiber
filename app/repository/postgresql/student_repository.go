package repository

import (
    "context"
    "database/sql"
    "errors"
    
    models "student-performance-report/app/models/postgresql"
    
    "github.com/google/uuid"
    "github.com/lib/pq" // [PENTING] Pastikan sudah go get github.com/lib/pq
)

// Update Interface: Tambahkan Context dan Method Baru
type StudentRepository interface {
    GetAllStudents(ctx context.Context) ([]models.Student, error)
    GetStudentByID(ctx context.Context, id uuid.UUID) (*models.Student, error)
    UpdateAdvisor(ctx context.Context, studentID, lecturerID uuid.UUID) error
    
    // [BARU] Method untuk Report Service
    GetStudentsByIDs(ctx context.Context, ids []string) ([]models.StudentWithUser, error)
}

type studentRepository struct {
    pg *sql.DB
}

func NewStudentRepository(pg *sql.DB) StudentRepository {
    return &studentRepository{pg: pg}
}

func (r *studentRepository) GetAllStudents(ctx context.Context) ([]models.Student, error) {
    // Tambahkan JOIN users jika ingin menampilkan nama di list, 
    // tapi query standar juga oke.
    query := `
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students
        ORDER BY created_at DESC
    `
    rows, err := r.pg.QueryContext(ctx, query)
    if err != nil { return nil, err }
    defer rows.Close()

    var list []models.Student
    for rows.Next() {
        var s models.Student
        // AdvisorID bisa null, gunakan scanner yang aman jika driver support, 
        // atau gunakan NullUUID intermediate.
        err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
            &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
        if err != nil { return nil, err }
        list = append(list, s)
    }
    return list, nil
}

// [DIPERBAIKI] Sekarang menggunakan JOIN untuk mengambil FullName
func (r *studentRepository) GetStudentByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
    var s models.Student
    
    // Query Join: Ambil data student + nama dari user
    query := `
        SELECT 
            s.id, s.user_id, s.student_id, s.program_study, s.academic_year, s.advisor_id, s.created_at,
            u.full_name
        FROM students s
        JOIN users u ON s.user_id = u.id
        WHERE s.id = $1
    `
    
    // Kita perlu handling AdvisorID karena bisa NULL di database
    var advisorID sql.NullString 

    err := r.pg.QueryRowContext(ctx, query, id).Scan(
        &s.ID, &s.UserID, &s.StudentID,
        &s.ProgramStudy, &s.AcademicYear,
        &advisorID, // Scan ke variable sementara
        &s.CreatedAt,
        &s.FullName, // [BARU] Isi field nama
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("student not found")
    } else if err != nil {
        return nil, err
    }

    // Convert NullString kembali ke *UUID
    if advisorID.Valid {
        uid, _ := uuid.Parse(advisorID.String)
        s.AdvisorID = &uid
    } else {
        s.AdvisorID = nil
    }

    return &s, nil
}

func (r *studentRepository) UpdateAdvisor(ctx context.Context, studentID, lecturerID uuid.UUID) error {
    query := `UPDATE students SET advisor_id=$1 WHERE id=$2`
    _, err := r.pg.ExecContext(ctx, query, lecturerID, studentID)
    return err
}

// [BARU] Implementasi GetStudentsByIDs untuk Report
func (r *studentRepository) GetStudentsByIDs(ctx context.Context, ids []string) ([]models.StudentWithUser, error) {
    if len(ids) == 0 {
        return []models.StudentWithUser{}, nil
    }

    // Menggunakan pq.Array untuk handling array di query Postgres
    query := `
        SELECT s.id, u.full_name, s.program_study
        FROM students s
        JOIN users u ON s.user_id = u.id
        WHERE s.id::text = ANY($1)
    `

    rows, err := r.pg.QueryContext(ctx, query, pq.Array(ids))
    if err != nil { return nil, err }
    defer rows.Close()

    var results []models.StudentWithUser
    for rows.Next() {
        var data models.StudentWithUser
        if err := rows.Scan(&data.ID, &data.FullName, &data.ProgramStudy); err != nil {
            return nil, err
        }
        results = append(results, data)
    }
    
    return results, nil
}