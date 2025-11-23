package service

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    // HAPUS IMPORT MODEL YANG TIDAK TERPAKAI
    
    repoMongo "student-performance-report/app/repository/mongodb"
    repoPg "student-performance-report/app/repository/postgresql"
)

type ReportService struct {
    mongoRepo   repoMongo.AchievementRepository
    studentRepo repoPg.StudentRepository
}

func NewReportService(m repoMongo.AchievementRepository, s repoPg.StudentRepository) *ReportService {
    return &ReportService{mongoRepo: m, studentRepo: s}
}

// 1. GET /api/v1/reports/statistics (Global Stats)
func (s *ReportService) GetStatistics(c *fiber.Ctx) error {
    ctx := c.Context()
    
    // A. Ambil Data Agregat dari Mongo
    stats, err := s.mongoRepo.GetGlobalStats(ctx)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to generate stats"})
    }

    // B. Ambil ID Top Student dari hasil Mongo
    var studentIDs []string
    for _, top := range stats.PointsDistribution {
        studentIDs = append(studentIDs, top.StudentID)
    }

    // C. Ambil Nama Mahasiswa dari Postgres berdasarkan IDs
    studentsWithDetails, _ := s.studentRepo.GetStudentsByIDs(ctx, studentIDs) 

    // D. Mapping Nama ke Hasil Statistik
    for i, top := range stats.PointsDistribution {
        for _, stud := range studentsWithDetails {
            if top.StudentID == stud.ID.String() {
                stats.PointsDistribution[i].Name = stud.FullName
                stats.PointsDistribution[i].ProgramStudy = stud.ProgramStudy
            }
        }
    }

    return c.JSON(stats)
}

// 2. GET /api/v1/reports/student/:id (Student Detail Stats)
func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
    ctx := c.Context()
    targetStudentID := c.Params("id")

    // A. Ambil Stats dari Mongo (Agregasi)
    stats, err := s.mongoRepo.GetStudentStats(ctx, targetStudentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get student stats"})
    }

    // B. Ambil Nama Mahasiswa
    studentUUID, err := uuid.Parse(targetStudentID)
    if err != nil {
         return c.Status(400).JSON(fiber.Map{"error": "Invalid UUID"})
    }

    // [FIX] Tambahkan parameter 'ctx' di sini
    studentProfile, err := s.studentRepo.GetStudentByID(ctx, studentUUID) 
    
    if err == nil && studentProfile != nil {
        stats.StudentName = studentProfile.FullName
    }

    return c.JSON(stats)
}
