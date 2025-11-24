package service_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	models "student-performance-report/app/models/postgresql"
	"student-performance-report/app/repository/mocks"
	"student-performance-report/app/service/postgresql"
)

func setupAdminTest() (*service.AdminService, *mocks.MockAdminRepo, *mocks.MockUserRepo) {
	mockAdminRepo := new(mocks.MockAdminRepo)
	mockUserRepo := new(mocks.MockUserRepo)
	svc := service.NewAdminService(mockAdminRepo, mockUserRepo)

	return svc, mockAdminRepo, mockUserRepo
}

func setupApp(roleName string, userID uuid.UUID) *fiber.App {
	app := fiber.New()

	// Middleware simulasi Auth
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_name", roleName)
		c.Locals("user_id", userID)
		return c.Next()
	})

	return app
}

func TestGetAllUsers(t *testing.T) {
	t.Run("Success: Admin gets all users", func(t *testing.T) {
		svc, mockRepo, _ := setupAdminTest()
		app := setupApp("admin", uuid.New())

		mockData := []models.User{
			{ID: uuid.New(), Username: "user1"},
			{ID: uuid.New(), Username: "user2"},
		}

		mockRepo.On("GetAllUsers").Return(mockData, nil)
		app.Get("/users", svc.GetAllUsers)

		req := httptest.NewRequest("GET", "/users", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Forbidden: Student cannot get users", func(t *testing.T) {
		svc, mockRepo, _ := setupAdminTest()
		app := setupApp("student", uuid.New())

		app.Get("/users", svc.GetAllUsers)

		req := httptest.NewRequest("GET", "/users", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 403, resp.StatusCode)
		mockRepo.AssertNotCalled(t, "GetAllUsers")
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Success: Admin creates user", func(t *testing.T) {
		svc, mockRepo, _ := setupAdminTest()
		app := setupApp("admin", uuid.New())

		inputPayload := models.User{
			Username:     "new_admin",
			PasswordHash: "raw_password",
			Email:        "admin@test.com",
		}

		mockRepo.On("CreateUser", mock.MatchedBy(func(u *models.User) bool {
			return u.Username == "new_admin" && u.PasswordHash != "raw_password"
		})).Return(nil)

		app.Post("/users", svc.CreateUser)

		body, _ := json.Marshal(inputPayload)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("Success: Get own profile (non-admin)", func(t *testing.T) {
		svc, mockRepo, _ := setupAdminTest()
		myID := uuid.New()
		app := setupApp("student", myID)

		mockUser := &models.User{ID: myID, Username: "me"}
		mockRepo.On("GetUserByID", myID).Return(mockUser, nil)

		app.Get("/users/:id", svc.GetUserByID)

		req := httptest.NewRequest("GET", "/users/"+myID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Forbidden: Get other profile (non-admin)", func(t *testing.T) {
		svc, mockRepo, _ := setupAdminTest()
		myID := uuid.New()
		otherID := uuid.New()
		app := setupApp("student", myID)

		app.Get("/users/:id", svc.GetUserByID)

		req := httptest.NewRequest("GET", "/users/"+otherID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 403, resp.StatusCode)
		mockRepo.AssertNotCalled(t, "GetUserByID")
	})
}