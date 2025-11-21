package models

import (
	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`         // Contoh: "achievement:create"
	Resource    string    `json:"resource" db:"resource"` // Contoh: "achievement"
	Action      string    `json:"action" db:"action"`     // Contoh: "create"
	Description string    `json:"description" db:"description"`
}

// RolePermission merepresentasikan tabel pivot/junction (Many-to-Many)
type RolePermission struct {
	RoleID       uuid.UUID `json:"roleId" db:"role_id"`
	PermissionID uuid.UUID `json:"permissionId" db:"permission_id"`
}