package middleware

import (
    "strings"
    "student-performance-report/utils"
    "github.com/gofiber/fiber/v2"
)

// AuthRequired: Validasi token & simpan claims ke Context
func AuthRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        auth := c.Get("Authorization")
        if auth == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
        }

        parts := strings.Split(auth, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
        }

        // Pastikan utils.ValidateToken mengembalikan struct Claims yang benar
        claims, err := utils.ValidateToken(parts[1])
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
        }

        // Simpan ke Locals agar bisa dipakai di middleware selanjutnya / controller
        c.Locals("user_id", claims.UserID)
        c.Locals("role_id", claims.RoleID)
        c.Locals("role_name", claims.RoleName) 
        c.Locals("permissions", claims.Permissions) // Pastikan ini []string

        return c.Next()
    }
}

// RoleAllowed: Cek Role (Case Insensitive)
func RoleAllowed(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("role_name")
        if role == nil {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "role missing in context"})
        }

        userRole, ok := role.(string)
        if !ok {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid role type"})
        }

        for _, r := range allowedRoles {
            // PERBAIKAN: Gunakan EqualFold agar "mahasiswa" dianggap sama dengan "Mahasiswa"
            if strings.EqualFold(userRole, r) {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden: role not allowed"})
    }
}

// PermissionRequired: Cek Permission Spesifik
func PermissionRequired(needed string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        raw := c.Locals("permissions")
        if raw == nil {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permissions found"})
        }

        // PERBAIKAN: Safe Type Assertion (mencegah panic)
        perms, ok := raw.([]string)
        if !ok {
            // Jika format data salah (misal interface{} bukan []string), return error jangan panic
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid permissions format"})
        }

        for _, p := range perms {
            // Permission biasanya case-sensitive (misal: achievement:create), jadi == aman
            if p == needed {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "permission denied: needed '" + needed + "'",
        })
    }
}