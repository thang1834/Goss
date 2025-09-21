package authentication

type RegisterRequest struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type AssignRoleRequest struct {
	UserID uint64 `json:"user_id"`
	RoleID uint64   `json:"role_id"`
}

type AssignPermissionRequest struct {
	UserID       uint64 `json:"user_id"`
	PermissionID uint64   `json:"permission_id"`
}
