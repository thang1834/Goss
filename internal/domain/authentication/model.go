package authentication

import "time"

type User struct {
	ID           uint64    `json:"id"`
	FirstName    string    `json:"first_name"`
	MiddleName   string    `json:"middle_name,omitempty"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	Status       string    `json:"status"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	VerifiedAt   time.Time `json:"verified_at,omitempty"`
}

type Role struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Permission struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRole struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	RoleID     uint64    `json:"role_id"`
	AssignedBy uint64    `json:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	IsActive   bool      `json:"is_active"`
}

type UserPermission struct {
	ID           uint64    `json:"id"`
	UserID       uint64    `json:"user_id"`
	PermissionID uint64    `json:"permission_id"`
	GrantedBy    uint64    `json:"granted_by"`
	GrantedAt    time.Time `json:"granted_at"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	IsActive     bool      `json:"is_active"`
}

// UserPermissionsCache represents cached user permissions data
type UserPermissionsCache struct {
	UserID      uint64    `json:"user_id"`
	Email       string    `json:"email"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	CachedAt    time.Time `json:"cached_at"`
}
