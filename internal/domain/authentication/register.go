package authentication

import (
	"net/http"

	"github.com/gmhafiz/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/thang1834/go-goss/internal/middleware"
)

// RegisterHTTPEndPoints registers authentication routes following go8 pattern
func RegisterHTTPEndPoints(router *chi.Mux, session *scs.SessionManager, repo Repo, redisAddr string) error {
	h, err := NewHandler(session, repo, redisAddr)
	if err != nil {
		return err
	}

	// Public authentication routes
	router.Post("/api/v1/login", h.Login)
	router.Post("/api/v1/register", h.Register)

	// Logout route
	router.Route("/api/v1/logout", func(router chi.Router) {
		router.Post("/", h.Logout)
	})

	// Protected routes
	router.Route("/api/v1/restricted", func(router chi.Router) {
		router.Use(middleware.Authenticate(session))

		// Basic protected routes
		router.Get("/", h.Protected)
		router.Get("/me", h.Me)
		router.Post("/change-password", h.ChangePassword)
		router.Get("/csrf", h.Csrf)

		// Admin routes - force logout other users
		router.Post("/logout/{userID}", h.ForceLogout)
	})

	// Admin management routes
	router.Route("/api/v1/admin", func(router chi.Router) {
		router.Use(middleware.Authenticate(session))
		router.Use(h.RequireAnyRole("admin", "super_admin"))

		// User management
		router.Get("/users", h.ListUsers)
		router.Get("/users/{userID}/roles", h.GetUserRoles)
		router.Post("/users/assign-role", h.AssignRole)
	})

	// Manager routes (example for e-commerce)
	router.Route("/api/v1/manage", func(router chi.Router) {
		router.Use(middleware.Authenticate(session))
		router.Use(h.RequireAnyRole("admin", "super_admin", "manager"))

		// Product management routes would go here
		router.Get("/products", func(w http.ResponseWriter, r *http.Request) {
			// Product management endpoint
		})
	})

	// Permission-based routes (example)
	router.Route("/api/v1/orders", func(router chi.Router) {
		router.Use(middleware.Authenticate(session))

		// Create order - requires permission
		router.With(h.RequirePermission("order:create")).Post("/", func(w http.ResponseWriter, r *http.Request) {
			// Create order endpoint
		})

		// Read all orders - admin permission
		router.With(h.RequirePermission("order:read_all")).Get("/all", func(w http.ResponseWriter, r *http.Request) {
			// List all orders endpoint
		})
	})

	return nil
}
