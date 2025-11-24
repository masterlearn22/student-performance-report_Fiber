package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	// Import Models
	models "student-performance-report/app/models/postgresql"

	// PENTING: Import package tempat Interface AdminRepository & UserRepository berada
	// Sesuaikan path ini jika interface Anda ada di folder lain (misal: app/repository)
	repo "student-performance-report/app/repository/postgresql"
)

// =========================================================
// MOCK ADMIN REPOSITORY
// =========================================================

type MockAdminRepo struct {
	mock.Mock
}

// Pastikan struct ini mengimplementasikan repo.AdminRepository
var _ repo.AdminRepository = (*MockAdminRepo)(nil)

func (m *MockAdminRepo) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAdminRepo) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAdminRepo) DeleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminRepo) GetUserByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	// Safety check agar tidak panic jika return nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAdminRepo) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockAdminRepo) AssignRole(userID, roleID uuid.UUID) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *MockAdminRepo) SetStudentProfile(profile *models.Student) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockAdminRepo) SetLecturerProfile(profile *models.Lecturer) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockAdminRepo) SetAdvisor(studentID, lecturerID uuid.UUID) error {
	args := m.Called(studentID, lecturerID)
	return args.Error(0)
}

// =========================================================
// MOCK USER REPOSITORY
// =========================================================

type MockUserRepo struct {
	mock.Mock
}

// Pastikan struct ini mengimplementasikan repo.UserRepository
var _ repo.UserRepository = (*MockUserRepo)(nil)

func (m *MockUserRepo) GetByUsername(username string) (*models.User, string, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func (m *MockUserRepo) GetPermissionsByRoleID(roleID uuid.UUID) ([]string, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRepo) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}