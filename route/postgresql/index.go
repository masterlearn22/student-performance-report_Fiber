package route

import (
	"database/sql"

	repoMongo "student-performance-report/app/repository/mongodb"
	repoPostgre "student-performance-report/app/repository/postgresql"
	service "student-performance-report/app/service/postgresql"
	"student-performance-report/database"
	"student-performance-report/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupPostgresRoutes(app *fiber.App, db *sql.DB) {

    userRepo := repoPostgre.NewUserRepository(db)
    authService := service.NewAuthService(userRepo)
    adminRepo := repoPostgre.NewAdminRepository(db)
    adminService := service.NewAdminService(adminRepo, userRepo)
    studentRepo := repoPostgre.NewStudentRepository(db)
    achievementRepo := repoMongo.NewAchievementRepository(database.MongoDB)
    studentService := service.NewStudentService(studentRepo, achievementRepo)
    lecturerRepo := repoPostgre.NewLecturerRepository(db)
    lecturerService := service.NewLecturerService(lecturerRepo)

    // AUTH
    auth := app.Group("/api/v1/auth")
    auth.Post("/login", authService.Login)
    auth.Post("/refresh", authService.Refresh)
    auth.Post("/logout", middleware.AuthRequired(), authService.Logout)
    auth.Get("/profile", middleware.AuthRequired(), authService.Profile)

    // USERS (protected)
    users := app.Group("/api/v1/users", middleware.AuthRequired())
    users.Get("/", middleware.RoleAllowed("admin"), adminService.GetAllUsers)
    users.Get("/:id", adminService.GetUserByID)
    users.Post("/", middleware.RoleAllowed("admin"), adminService.CreateUser)
    users.Put("/:id", adminService.UpdateUser)
    users.Delete("/:id", adminService.DeleteUser)
    users.Put("/:id/role", middleware.RoleAllowed("admin"), adminService.AssignRole)


     // Students and Lecturers
    auth.Get("/students", studentService.GetAllStudents)
    auth.Get("/students/:id", studentService.GetStudentByID)
    auth.Get("/students/:id/achievements", studentService.GetStudentAchievements)
    auth.Put("/students/:id/advisor", studentService.UpdateAdvisor)
    auth.Get("/lecturers", lecturerService.GetAllLecturers)
    auth.Get("/lecturers/:id/advisees", lecturerService.GetAdvisees)
}
