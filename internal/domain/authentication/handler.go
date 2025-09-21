package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gmhafiz/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/thang1834/go-goss/internal/middleware"
	"github.com/thang1834/go-goss/internal/utility/param"
	"github.com/thang1834/go-goss/internal/utility/request"
	"github.com/thang1834/go-goss/internal/utility/respond"
)

const (
	minPasswordLength = 8
)

var (
	ErrEmailRequired     = errors.New("email is required")
	ErrPasswordLength    = fmt.Errorf("password must be at least %d characters", minPasswordLength)
	ErrFirstNameRequired = errors.New("first name is required")
	ErrLastNameRequired  = errors.New("last name is required")
)

type Handler struct {
	repo        Repo
	session     *scs.SessionManager
	redisClient *redis.Client
}

// NewHandler creates new handler with Redis caching
func NewHandler(session *scs.SessionManager, repo Repo, redisAddr string) (*Handler, error) {
	// Setup Redis client for permission caching (different DB from sessions)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   1, // Different DB from sessions
	})

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Handler{
		repo:        repo,
		session:     session,
		redisClient: redisClient,
	}, nil
}

// Register handles user registration following go8 pattern
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	// Trim whitespace
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	// Validate request
	if req.Email == "" {
		respond.Error(w, http.StatusBadRequest, ErrEmailRequired)
		return
	}
	if req.FirstName == "" {
		respond.Error(w, http.StatusBadRequest, ErrFirstNameRequired)
		return
	}
	if req.LastName == "" {
		respond.Error(w, http.StatusBadRequest, ErrLastNameRequired)
		return
	}
	if len(req.Password) < minPasswordLength {
		respond.Error(w, http.StatusBadRequest, ErrPasswordLength)
		return
	}

	user, err := h.repo.Register(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	// Auto-login after registration
	if err := h.session.RenewToken(r.Context()); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	h.session.Put(r.Context(), string(middleware.KeyID), user.ID)

	respond.Status(w, http.StatusCreated)
}

// Login handles user login following go8 pattern
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	ctx := r.Context()

	user, match, err := h.repo.Login(ctx, req)
	if err != nil || !match {
		respond.Status(w, http.StatusUnauthorized)
		return
	}

	if err := h.session.RenewToken(ctx); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	h.session.Put(ctx, string(middleware.KeyID), user.ID)

	respond.Status(w, http.StatusOK)
}

// Protected endpoint following go8 pattern
func (h *Handler) Protected(w http.ResponseWriter, _ *http.Request) {
	respond.Json(w, http.StatusOK, map[string]string{"success": "yup!"})
}

// Me returns current user info with roles and permissions
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, errors.New("authentication required"))
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, errors.New("user not found"))
		return
	}

	roles, _ := h.repo.GetUserRoles(r.Context(), userID)
	permissions, _ := h.repo.GetUserPermissions(r.Context(), userID)

	userResp := AuthStatusResponse{
		UserID:      userID,
		Email:       user.Email,
		Roles:       roles,
		Permissions: permissions,
	}

	respond.Json(w, http.StatusOK, userResp)
}

// ChangePassword handles password change
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, errors.New("authentication required"))
		return
	}

	if len(req.NewPassword) < minPasswordLength {
		respond.Error(w, http.StatusBadRequest, ErrPasswordLength)
		return
	}

	err = h.repo.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			respond.Error(w, http.StatusBadRequest, errors.New("current password is incorrect"))
			return
		}
		respond.Error(w, http.StatusInternalServerError, errors.New("failed to change password"))
		return
	}

	// Clear user cache after password change
	h.invalidateUserCache(r.Context(), userID)

	respond.Status(w, http.StatusOK)
}

// Logout handles user logout following go8 pattern
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear user cache
	if userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64); ok {
		h.invalidateUserCache(r.Context(), userID)
	}

	err := h.session.Destroy(r.Context())
	if err != nil {
		respond.Status(w, http.StatusBadRequest)
		return
	}

	respond.Status(w, http.StatusOK)
}

// ForceLogout allows admin to log out other users following go8 pattern
func (h *Handler) ForceLogout(w http.ResponseWriter, r *http.Request) {
	// Authorization check - only super admin can force logout
	currUser := h.session.Get(r.Context(), string(middleware.KeyID))
	currUserID, ok := currUser.(uint64)
	if !ok {
		respond.Status(w, http.StatusUnauthorized)
		return
	}

	// Check if current user has admin permissions
	if !h.hasPermission(r.Context(), currUserID, "admin:users") {
		respond.Status(w, http.StatusForbidden)
		return
	}

	userID, err := param.UInt64(r, "userID")
	if err != nil {
		respond.Status(w, http.StatusBadRequest)
		return
	}

	// Clear user cache first
	h.invalidateUserCache(r.Context(), userID)

	ok, err = h.repo.Logout(r.Context(), userID)
	if err != nil {
		respond.Status(w, http.StatusInternalServerError)
		return
	}

	if !ok {
		respond.Json(w, http.StatusInternalServerError, map[string]string{"message": "unable to log out"})
		return
	}

	respond.Status(w, http.StatusOK)
}

// Csrf handles CSRF token generation following go8 pattern
func (h *Handler) Csrf(w http.ResponseWriter, r *http.Request) {
	_, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(int64)
	if !ok {
		respond.Error(w, http.StatusBadRequest, errors.New("you need to be logged in"))
		return
	}

	token, err := h.repo.Csrf(r.Context())
	if err != nil {
		respond.Status(w, http.StatusInternalServerError)
		return
	}

	respond.Json(w, http.StatusOK, &RespondCsrf{CsrfToken: token})
}

// Admin endpoints

// ListUsers lists all users (admin only)
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
	if !ok {
		respond.Status(w, http.StatusUnauthorized)
		return
	}

	if !h.hasPermission(r.Context(), userID, "user:read_all") {
		respond.Status(w, http.StatusForbidden)
		return
	}

	// Implementation would query all users
	respond.Json(w, http.StatusOK, map[string]string{"message": "List of users"})
}

// AssignRole assigns role to user (admin only)
func (h *Handler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req AssignRoleRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	adminUserID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
	if !ok {
		respond.Status(w, http.StatusUnauthorized)
		return
	}

	if !h.hasPermission(r.Context(), adminUserID, "role:assign") {
		respond.Status(w, http.StatusForbidden)
		return
	}

	err = h.repo.AssignRoleToUser(r.Context(), req.UserID, req.RoleID, adminUserID)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	// Clear user cache
	h.invalidateUserCache(r.Context(), req.UserID)

	respond.Status(w, http.StatusOK)
}

// GetUserRoles gets user roles (admin only)
func (h *Handler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userIDParam, err := param.UInt64(r, "userID")
	if err != nil {
		respond.Status(w, http.StatusBadRequest)
		return
	}

	adminUserID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
	if !ok {
		respond.Status(w, http.StatusUnauthorized)
		return
	}

	if !h.hasPermission(r.Context(), adminUserID, "user:read") {
		respond.Status(w, http.StatusForbidden)
		return
	}

	roles, err := h.repo.GetUserRoles(r.Context(), userIDParam)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.Json(w, http.StatusOK, map[string]interface{}{"roles": roles})
}

// RBAC Helper methods

// hasPermission checks if user has specific permission (with caching)
func (h *Handler) hasPermission(ctx context.Context, userID uint64, permission string) bool {
	permissions, err := h.loadUserPermissions(ctx, userID)
	if err != nil {
		return false
	}

	for _, p := range permissions.Permissions {
		if p == permission || p == "*:*" {
			return true
		}
	}
	return false
}

// hasRole checks if user has specific role
func (h *Handler) hasRole(ctx context.Context, userID uint64, role string) bool {
	permissions, err := h.loadUserPermissions(ctx, userID)
	if err != nil {
		return false
	}

	for _, r := range permissions.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// loadUserPermissions loads user permissions with Redis caching
func (h *Handler) loadUserPermissions(ctx context.Context, userID uint64) (*UserPermissionsCache, error) {
	cacheKey := fmt.Sprintf("user_perms:%d", userID)

	// Try Redis cache first
	cached, err := h.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var userPerms UserPermissionsCache
		if json.Unmarshal([]byte(cached), &userPerms) == nil {
			return &userPerms, nil
		}
	}

	// Cache miss - load from database
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles, err := h.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissions, err := h.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	userPerms := &UserPermissionsCache{
		UserID:      userID,
		Email:       user.Email,
		Roles:       roles,
		Permissions: permissions,
		CachedAt:    time.Now(),
	}

	// Cache for 10 minutes
	permsJSON, _ := json.Marshal(userPerms)
	h.redisClient.Set(ctx, cacheKey, permsJSON, 10*time.Minute)

	return userPerms, nil
}

// invalidateUserCache clears user permissions cache
func (h *Handler) invalidateUserCache(ctx context.Context, userID uint64) {
	cacheKey := fmt.Sprintf("user_perms:%d", userID)
	h.redisClient.Del(ctx, cacheKey)
}

// Middleware functions following go8 pattern

// RequirePermission middleware for permission-based access control
func (h *Handler) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
			if !ok {
				respond.Status(w, http.StatusUnauthorized)
				return
			}

			if !h.hasPermission(r.Context(), userID, permission) {
				respond.Status(w, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole middleware for role-based access control
func (h *Handler) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
			if !ok {
				respond.Status(w, http.StatusUnauthorized)
				return
			}

			if !h.hasRole(r.Context(), userID, role) {
				respond.Status(w, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole middleware that requires any of the specified roles
func (h *Handler) RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
			if !ok {
				respond.Status(w, http.StatusUnauthorized)
				return
			}

			hasRole := false
			for _, role := range roles {
				if h.hasRole(r.Context(), userID, role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				respond.Status(w, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission middleware that requires any of the specified permissions
func (h *Handler) RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
			if !ok {
				respond.Status(w, http.StatusUnauthorized)
				return
			}

			hasPermission := false
			for _, perm := range permissions {
				if h.hasPermission(r.Context(), userID, perm) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				respond.Status(w, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
