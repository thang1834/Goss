package authentication

type UserResponse struct {
	ID         uint64   `json:"id"`
	FirstName  string   `json:"first_name"`
	MiddleName string   `json:"middle_name,omitempty"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone,omitempty"`
	Status     string   `json:"status"`
	Roles      []string `json:"roles"`
	CreatedAt  string   `json:"created_at"`
}

type RespondCsrf struct {
	CsrfToken string `json:"csrf_token"`
}

type AuthStatusResponse struct {
	UserID      uint64   `json:"user_id"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type RoleResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type PermissionResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}
