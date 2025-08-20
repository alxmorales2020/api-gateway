package admin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/alxmorales2020/api-gateway/config"
	"github.com/alxmorales2020/api-gateway/router"
)

type AdminHandler struct {
	store    config.RouteStore
	reloader router.Reloader
}

func NewAdminHandler(store config.RouteStore, reloader router.Reloader) *AdminHandler {
	return &AdminHandler{store: store, reloader: reloader}
}

// Routes registers admin endpoints
func (h *AdminHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/routes", func(r chi.Router) {
		r.Get("/", h.GetRoutes)          // GET    /admin/routes
		r.Post("/", h.CreateRoute)       // POST   /admin/routes
		r.Delete("/{id}", h.DeleteRoute) // DELETE /admin/routes/{id}
	})

	// Helpful: see 405 vs 404 clearly
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("ADMIN 405: %s %s", req.Method, req.URL.Path)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("ADMIN 404: %s %s", req.Method, req.URL.Path)
		http.Error(w, "admin route not found", http.StatusNotFound)
	})

	return r
}

// GET /admin/routes
func (h *AdminHandler) GetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := h.store.LoadRoutes()
	if err != nil {
		http.Error(w, "Failed to load routes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

// POST /admin/routes
func (h *AdminHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var route config.RouteConfig
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		http.Error(w, "Invalid route data", http.StatusBadRequest)
		return
	}

	// Basic validation
	if route.Path == "" || route.Upstream == "" || len(route.Methods) == 0 {
		http.Error(w, "Missing required route fields", http.StatusBadRequest)
		return
	}

	if err := h.store.SaveRoute(&route); err != nil {
		http.Error(w, "Failed to save route", http.StatusInternalServerError)
		return
	}

	// Hot-reload the router after save
	if err := h.reloader.Reload(); err != nil {
		http.Error(w, "saved but reload failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":      route.ID,
		"message": "Route saved",
	})
}

// DELETE /admin/routes/{id}
func (h *AdminHandler) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing route id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteRoute(id); err != nil {
		// You can map specific errors to 404 if your store returns them
		http.Error(w, "route not found", http.StatusNotFound)
		return
	}

	// Hot-reload the router after delete
	if err := h.reloader.Reload(); err != nil {
		http.Error(w, "deleted but reload failed", http.StatusInternalServerError)
		return
	}

	// Either 204 No Content (common for delete)...
	w.WriteHeader(http.StatusNoContent)
}
