package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    models "student-performance-report/app/models/postgresql"
    "github.com/google/uuid"
)

type AchievementRepoPostgres interface {
    Create(ctx context.Context, ref models.AchievementReference) (uuid.UUID, error)
    GetStudentByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
    GetAllReferences(ctx context.Context, filter map[string]interface{}) ([]models.AchievementReference, error)
    GetReferenceByID(ctx context.Context, id uuid.UUID) (models.AchievementReference, error)
    DeleteReference(ctx context.Context, id uuid.UUID) error
    UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy *uuid.UUID, note string) error
    SubmitReference(ctx context.Context, id uuid.UUID) error
}

type achievementRepoPostgres struct {
    db *sql.DB
}

func NewAchievementRepoPostgres(db *sql.DB) AchievementRepoPostgres {
    return &achievementRepoPostgres{db: db}
}

func (r *achievementRepoPostgres) GetStudentByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
    query := `SELECT id FROM students WHERE user_id = $1`
    var studentID uuid.UUID
    err := r.db.QueryRowContext(ctx, query, userID).Scan(&studentID)
    return studentID, err
}

func (r *achievementRepoPostgres) Create(ctx context.Context, ref models.AchievementReference) (uuid.UUID, error) {
    query := `
        INSERT INTO achievement_references (
            student_id, mongo_achievement_id, status, created_at, updated_at
        ) VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id
    `
    var newID uuid.UUID
    err := r.db.QueryRowContext(ctx, query, 
        ref.StudentID, 
        ref.MongoAchievementID, 
        ref.Status, 
    ).Scan(&newID)

    return newID, err
}

func (r *achievementRepoPostgres) GetAllReferences(ctx context.Context, filter map[string]interface{}) ([]models.AchievementReference, error) {
    query := `SELECT id, student_id, mongo_achievement_id, status, created_at FROM achievement_references WHERE 1=1`
    var args []interface{}
    argCount := 1

    // Filter Student ID
    if val, ok := filter["student_id"]; ok {
        query += fmt.Sprintf(" AND student_id = $%d", argCount)
        args = append(args, val)
        argCount++
    }

    query += ` ORDER BY created_at DESC`

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []models.AchievementReference
    for rows.Next() {
        var ref models.AchievementReference
        // Pastikan urutan Scan sesuai SELECT
        if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &ref.CreatedAt); err != nil {
            return nil, err
        }
        results = append(results, ref)
    }
    return results, nil
}

func (r *achievementRepoPostgres) GetReferenceByID(ctx context.Context, id uuid.UUID) (models.AchievementReference, error) {
    // PERBAIKAN 1: Tambahkan kolom submitted_at, verified_at, verified_by di SELECT
    query := `
        SELECT 
            id, student_id, mongo_achievement_id, status, rejection_note, 
            created_at, submitted_at, verified_at, verified_by 
        FROM achievement_references 
        WHERE id = $1
    `
    
    var ref models.AchievementReference
    // Menggunakan NullString untuk rejection_note karena bisa null
    var rejectionNote sql.NullString

    // PERBAIKAN 2: Tambahkan scan variable untuk field baru tersebut
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &ref.ID, 
        &ref.StudentID, 
        &ref.MongoAchievementID, 
        &ref.Status, 
        &rejectionNote,
        &ref.CreatedAt,
        &ref.SubmittedAt, 
        &ref.VerifiedAt,  
        &ref.VerifiedBy,  
    )

    if rejectionNote.Valid {
        note := rejectionNote.String
        ref.RejectionNote = &note
    }

    return ref, err
}

// [BARU] Hapus Referensi
func (r *achievementRepoPostgres) DeleteReference(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM achievement_references WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

// [BARU] Update Status (Verify/Reject)
func (r *achievementRepoPostgres) UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy *uuid.UUID, note string) error {
    // Kita update status, verified_by, verified_at, dan rejection_note sekaligus
    query := `
        UPDATE achievement_references 
        SET status = $1, verified_by = $2, verified_at = $3, rejection_note = $4, updated_at = NOW()
        WHERE id = $5
    `
    _, err := r.db.ExecContext(ctx, query, status, verifiedBy, time.Now(), note, id)
    return err
}

// 2. Implementasi Function
func (r *achievementRepoPostgres) SubmitReference(ctx context.Context, id uuid.UUID) error {
    // Query ini khusus update status jadi 'submitted' DAN mengisi submitted_at dengan NOW()
    query := `
        UPDATE achievement_references 
        SET status = 'submitted', 
            submitted_at = NOW(), 
            updated_at = NOW()
        WHERE id = $1
    `
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}