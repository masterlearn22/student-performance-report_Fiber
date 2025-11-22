package route

import (
    "database/sql"

    // Import Repositories
    repoMongo "student-performance-report/app/repository/mongodb"
    repoPostgre "student-performance-report/app/repository/postgresql"
    
    // Import Services (Sesuaikan path package service Anda)
    postgreService "student-performance-report/app/service/postgresql" 
    mongoService "student-performance-report/app/service/mongodb"
    
    "student-performance-report/database"
    "student-performance-report/middleware"

    "github.com/gofiber/fiber/v2"
)

func SetupPostgresRoutes(app *fiber.App, db *sql.DB) {
    // === Initialization ===
    
    // Repositories
    userRepo := repoPostgre.NewUserRepository(db)
    adminRepo := repoPostgre.NewAdminRepository(db)
    studentRepo := repoPostgre.NewStudentRepository(db)
    lecturerRepo := repoPostgre.NewLecturerRepository(db)
    
    // Achievement Repos (Dual DB)
    achRepoPg := repoPostgre.NewAchievementRepoPostgres(db)
    achRepoMongo := repoMongo.NewAchievementRepository(database.MongoDB)

    // Services
    authService := postgreService.NewAuthService(userRepo)
    adminService := postgreService.NewAdminService(adminRepo, userRepo)
    lecturerService := postgreService.NewLecturerService(lecturerRepo)
    
    // Student Service & Achievement Service
    // Perhatikan: StudentService butuh achRepoMongo untuk "GetStudentAchievements" (view only)
    studentService := postgreService.NewStudentService(studentRepo, achRepoMongo)
    
    // AchievementService untuk CRUD Logic lengkap (Create/Update/Dual Write)
    achievementService := mongoService.NewAchievementService(achRepoMongo, achRepoPg, lecturerRepo)
    app.Static("/uploads", "./uploads")


    // === Route Definitions ===

    // 1. Authentication
    auth := app.Group("/api/v1/auth")
    auth.Post("/login", authService.Login)
    auth.Post("/refresh", authService.Refresh)
    auth.Post("/logout", middleware.AuthRequired(), authService.Logout)
    auth.Get("/profile", middleware.AuthRequired(), authService.Profile)

    // 2. Users (Admin)
    users := app.Group("/api/v1/users", middleware.AuthRequired())
    users.Get("/", middleware.RoleAllowed("admin"), adminService.GetAllUsers)
    users.Get("/:id", adminService.GetUserByID)
    users.Post("/", middleware.RoleAllowed("admin"), adminService.CreateUser)
    users.Put("/:id", adminService.UpdateUser)
    users.Delete("/:id", adminService.DeleteUser)
    users.Put("/:id/role", middleware.RoleAllowed("admin"), adminService.AssignRole)

    // 3. Students & Lecturers
    student := app.Group("/api/v1", middleware.AuthRequired())
    student.Get("/students", studentService.GetAllStudents)
    student.Get("/students/:id", studentService.GetStudentByID)
    student.Get("/students/:id/achievements", studentService.GetStudentAchievements)
    student.Put("/students/:id/advisor", studentService.UpdateAdvisor)
    
    student.Get("/lecturers", lecturerService.GetAllLecturers)
    student.Get("/lecturers/:id/advisees", lecturerService.GetAdvisees)

    
    // Group Achievements
    ach := app.Group("/api/v1/achievements")

    // ==========================================
    // 1. GENERAL / COMMON ROUTES (List & Detail)
    // ==========================================
    // Diakses oleh Mahasiswa (milik sendiri), Dosen (bimbingan), Admin (semua)
    // Permission: achievement:read [cite: 60]
    
    // GET /api/v1/achievements (List)
    ach.Get("/", 
        middleware.PermissionRequired("achievement:read"), 
        achievementService.GetAllAchievements, // Service ini sudah kita buat
    )

    // GET /api/v1/achievements/:id (Detail)
    ach.Get("/:id", 
        middleware.PermissionRequired("achievement:read"), 
        achievementService.GetAchievementDetail, // Service ini sudah kita buat
    )

    // GET /api/v1/achievements/:id/history (Status History)
    ach.Get("/:id/history",
        middleware.PermissionRequired("achievement:read"),
        achievementService.GetAchievementHistory, 
    )


    // ==========================================
    // 2. MAHASISWA ROUTES (CRUD & Submission)
    // ==========================================
    
    // POST /api/v1/achievements (Create)
    // Permission: achievement:create [cite: 59]
    ach.Post("/", 
        middleware.RoleAllowed("mahasiswa"), 
        middleware.PermissionRequired("achievement:create"), 
        achievementService.CreateAchievement, // Service ini sudah kita buat
    )

    // PUT /api/v1/achievements/:id (Update Data - Draft only)
    // Permission: achievement:update [cite: 61]
    ach.Put("/:id",
        middleware.RoleAllowed("mahasiswa"),
        middleware.PermissionRequired("achievement:update"),
        achievementService.UpdateAchievement, 
    )

    // DELETE /api/v1/achievements/:id (Delete - Draft only)
    // Permission: achievement:delete [cite: 62]
    ach.Delete("/:id", 
        middleware.RoleAllowed("mahasiswa"), 
        middleware.PermissionRequired("achievement:delete"), 
        achievementService.DeleteAchievement, // Service ini sudah kita buat
    )

    // POST /api/v1/achievements/:id/submit (Submit for verification)
    // Permission: achievement:update (Mengubah status adalah bentuk update) [cite: 61]
    ach.Post("/:id/submit",
        middleware.RoleAllowed("mahasiswa"),
        middleware.PermissionRequired("achievement:update"),
        achievementService.SubmitAchievement, 
    )

    // POST /api/v1/achievements/:id/attachments (Upload Files)
    // Permission: achievement:update [cite: 61]
    ach.Post("/:id/attachments",
        middleware.RoleAllowed("mahasiswa"),
        middleware.PermissionRequired("achievement:update"),
        achievementService.UploadAttachments, 
    )


    // ==========================================
    // 3. DOSEN WALI ROUTES (Verification)
    // ==========================================
    
    // POST /api/v1/achievements/:id/verify
    // Permission: achievement:verify [cite: 65]
    ach.Post("/:id/verify", 
        middleware.RoleAllowed("dosen_wali"), 
        middleware.PermissionRequired("achievement:verify"), 
        achievementService.VerifyAchievement, // Service ini sudah kita buat
    )
    
    // POST /api/v1/achievements/:id/reject
    // Permission: achievement:verify (Reject bagian dari proses verifikasi) [cite: 65]
    ach.Post("/:id/reject", 
        middleware.RoleAllowed("dosen_wali"), 
        middleware.PermissionRequired("achievement:verify"), 
        achievementService.RejectAchievement, // Service ini sudah kita buat
    )
}