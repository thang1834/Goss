package authentication

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gmhafiz/scs/v2"

	"github.com/thang1834/go-goss/ent/gen"
	"github.com/thang1834/go-goss/ent/gen/role"
	"github.com/thang1834/go-goss/ent/gen/session"
	"github.com/thang1834/go-goss/ent/gen/user"
	"github.com/thang1834/go-goss/ent/gen/userpermission"
	"github.com/thang1834/go-goss/ent/gen/userrole"
)

type repo struct {
	ent     *gen.Client
	db      *sql.DB
	session *scs.SessionManager
}

var (
	ErrEmailNotAvailable = errors.New("email is not available")
	ErrNotLoggedIn       = errors.New("you are not logged in yet")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserInactive      = errors.New("user account is inactive")
)

type Repo interface {
	// Authentication
	Register(ctx context.Context, req RegisterRequest) (*gen.User, error)
	Login(ctx context.Context, req LoginRequest) (*gen.User, bool, error)
	Logout(ctx context.Context, userID uint64) (bool, error)
	ChangePassword(ctx context.Context, userID uint64, currentPassword, newPassword string) error

	// User operations
	GetUserByID(ctx context.Context, userID uint64) (*gen.User, error)
	GetUserByEmail(ctx context.Context, email string) (*gen.User, error)
	UpdateUserStatus(ctx context.Context, userID uint64, status string) error

	// Role and Permission operations
	GetUserRoles(ctx context.Context, userID uint64) ([]string, error)
	GetUserPermissions(ctx context.Context, userID uint64) ([]string, error)
	GetUserWithRoles(ctx context.Context, userID uint64) (*gen.User, error)

	// Role management
	GetAllRoles(ctx context.Context) ([]*gen.Role, error)
	AssignRoleToUser(ctx context.Context, userID uint64, roleID uint64, assignedBy uint64) error
	RemoveRoleFromUser(ctx context.Context, userID uint64, roleID uint64) error

	// Permission management
	GetAllPermissions(ctx context.Context) ([]*gen.Permission, error)
	AssignPermissionToUser(ctx context.Context, userID uint64, permissionID uint64, grantedBy uint64) error
	RemovePermissionFromUser(ctx context.Context, userID uint64, permissionID uint64) error

	// CSRF token
	Csrf(ctx context.Context) (string, error)
}

// Authentication methods
func (r *repo) Register(ctx context.Context, req RegisterRequest) (*gen.User, error) {
	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	builder := r.ent.User.Create().
		SetFirstName(req.FirstName).
		SetLastName(req.LastName).
		SetEmail(req.Email).
		SetPasswordHash(hashedPassword).
		SetStatus("active").
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now())

	if req.MiddleName != "" {
		builder = builder.SetMiddleName(req.MiddleName)
	}
	if req.Phone != "" {
		builder = builder.SetPhone(req.Phone)
	}

	user, err := builder.Save(ctx)
	if err != nil {
		if gen.IsConstraintError(err) {
			return nil, ErrEmailNotAvailable
		}
		return nil, err
	}

	// Auto-assign customer role
	customerRole, err := r.ent.Role.Query().Where(role.NameEQ("customer")).First(ctx)
	if err == nil {
		r.AssignRoleToUser(ctx, user.ID, customerRole.ID, user.ID)
	}

	return user, nil
}

func (r *repo) Login(ctx context.Context, req LoginRequest) (*gen.User, bool, error) {
	u, err := r.ent.User.Query().Where(user.EmailEqualFold(req.Email)).First(ctx)
	if err != nil {
		return nil, false, ErrUserNotFound
	}

	if u.Status != "active" {
		return nil, false, ErrUserInactive
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, u.PasswordHash)
	if err != nil {
		return nil, false, ErrInvalidPassword
	}

	return u, match, nil
}

func (r *repo) Logout(ctx context.Context, userID uint64) (bool, error) {
	var found bool
	rows := r.db.QueryRowContext(ctx, `
		SELECT CASE
		WHEN EXISTS(SELECT *
					FROM sessions
					WHERE sessions.user_id = $1)
			THEN true
		ELSE false
		END
	;`, userID)

	err := rows.Scan(&found)
	if err != nil {
		return false, err
	}

	if !found {
		return false, ErrNotLoggedIn
	}

	_, err = r.ent.Session.Delete().Where(session.UserIDEQ(userID)).Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repo) ChangePassword(ctx context.Context, userID uint64, currentPassword, newPassword string) error {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	match, err := argon2id.ComparePasswordAndHash(currentPassword, user.PasswordHash)
	if err != nil || !match {
		return ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := argon2id.CreateHash(newPassword, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	// Update password
	return r.ent.User.UpdateOneID(userID).
		SetPasswordHash(hashedPassword).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

// User operations
func (r *repo) GetUserByID(ctx context.Context, userID uint64) (*gen.User, error) {
	return r.ent.User.Query().Where(user.IDEQ(userID)).Only(ctx)
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*gen.User, error) {
	return r.ent.User.Query().Where(user.EmailEQ(email)).Only(ctx)
}

func (r *repo) UpdateUserStatus(ctx context.Context, userID uint64, status string) error {
	return r.ent.User.UpdateOneID(userID).
		SetStatus(status).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

// Role and Permission operations
func (r *repo) GetUserWithRoles(ctx context.Context, userID uint64) (*gen.User, error) {
	return r.ent.User.
		Query().
		Where(user.IDEQ(userID)).
		WithUserRoles(func(q *gen.UserRoleQuery) {
			q.Where(userrole.IsActiveEQ(true)).
				Where(func(s *gen.Selector) {
					s.Where(
						gen.Or(
							gen.IsNull("expires_at"),
							gen.GT("expires_at", time.Now()),
						),
					)
				}).
				WithRole(func(rq *gen.RoleQuery) {
					rq.Where(role.IsActiveEQ(true))
				})
		}).
		Only(ctx)
}

func (r *repo) GetUserRoles(ctx context.Context, userID uint64) ([]string, error) {
	userRoles, err := r.ent.User.
		Query().
		Where(user.IDEQ(userID)).
		QueryUserRoles().
		Where(
			userrole.IsActiveEQ(true),
			func(s *gen.Selector) {
				s.Where(
					gen.Or(
						gen.IsNull("expires_at"),
						gen.GT("expires_at", time.Now()),
					),
				)
			},
		).
		QueryRole().
		Where(role.IsActiveEQ(true)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	roles := make([]string, len(userRoles))
	for i, r := range userRoles {
		roles[i] = r.Name
	}

	return roles, nil
}

func (r *repo) GetUserPermissions(ctx context.Context, userID uint64) ([]string, error) {
	// Get direct user permissions
	directPermissions, err := r.ent.User.
		Query().
		Where(user.IDEQ(userID)).
		QueryUserPermissions().
		Where(
			userpermission.IsActiveEQ(true),
			func(s *gen.Selector) {
				s.Where(
					gen.Or(
						gen.IsNull("expires_at"),
						gen.GT("expires_at", time.Now()),
					),
				)
			},
		).
		QueryPermission().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Get role-based permissions
	rolePermissions, err := r.ent.User.
		Query().
		Where(user.IDEQ(userID)).
		QueryUserRoles().
		Where(
			userrole.IsActiveEQ(true),
			func(s *gen.Selector) {
				s.Where(
					gen.Or(
						gen.IsNull("expires_at"),
						gen.GT("expires_at", time.Now()),
					),
				)
			},
		).
		QueryRole().
		Where(role.IsActiveEQ(true)).
		QueryRolePermissions().
		QueryPermission().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Combine and deduplicate permissions
	permissionMap := make(map[string]bool)

	for _, p := range directPermissions {
		permissionMap[p.Name] = true
	}

	for _, p := range rolePermissions {
		permissionMap[p.Name] = true
	}

	permissions := make([]string, 0, len(permissionMap))
	for p := range permissionMap {
		permissions = append(permissions, p)
	}

	return permissions, nil
}

// Role management
func (r *repo) GetAllRoles(ctx context.Context) ([]*gen.Role, error) {
	return r.ent.Role.Query().Where(role.IsActiveEQ(true)).All(ctx)
}

func (r *repo) AssignRoleToUser(ctx context.Context, userID uint64, roleID uint64, assignedBy uint64) error {
	// Check if assignment already exists
	exists, err := r.ent.UserRole.
		Query().
		Where(
			userrole.UserIDEQ(userID),
			userrole.RoleIDEQ(roleID),
			userrole.IsActiveEQ(true),
		).
		Exist(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil // Already assigned
	}

	// Create new assignment
	return r.ent.UserRole.
		Create().
		SetUserID(userID).
		SetRoleID(roleID).
		SetAssignedBy(assignedBy).
		SetIsActive(true).
		SetAssignedAt(time.Now()).
		Exec(ctx)
}

func (r *repo) RemoveRoleFromUser(ctx context.Context, userID uint64, roleID uint64) error {
	return r.ent.UserRole.
		Update().
		Where(
			userrole.UserIDEQ(userID),
			userrole.RoleIDEQ(roleID),
		).
		SetIsActive(false).
		Exec(ctx)
}

// Permission management
func (r *repo) GetAllPermissions(ctx context.Context) ([]*gen.Permission, error) {
	return r.ent.Permission.Query().All(ctx)
}

func (r *repo) AssignPermissionToUser(ctx context.Context, userID uint64, permissionID uint64, grantedBy uint64) error {
	// Check if assignment already exists
	exists, err := r.ent.UserPermission.
		Query().
		Where(
			userpermission.UserIDEQ(userID),
			userpermission.PermissionIDEQ(permissionID),
			userpermission.IsActiveEQ(true),
		).
		Exist(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil // Already assigned
	}

	// Create new assignment
	return r.ent.UserPermission.
		Create().
		SetUserID(userID).
		SetPermissionID(permissionID).
		SetGrantedBy(grantedBy).
		SetIsActive(true).
		SetGrantedAt(time.Now()).
		Exec(ctx)
}

func (r *repo) RemovePermissionFromUser(ctx context.Context, userID uint64, permissionID uint64) error {
	return r.ent.UserPermission.
		Update().
		Where(
			userpermission.UserIDEQ(userID),
			userpermission.PermissionIDEQ(permissionID),
		).
		SetIsActive(false).
		Exec(ctx)
}

// CSRF token
func (r *repo) Csrf(ctx context.Context) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	err = r.session.CtxStore.CommitCtx(ctx, token, []byte("csrf_token"), time.Now().Add(r.session.Lifetime))
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewRepo(ent *gen.Client, db *sql.DB, manager *scs.SessionManager) *repo {
	return &repo{
		ent:     ent,
		db:      db,
		session: manager,
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
