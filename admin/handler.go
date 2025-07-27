package admin

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/alxmorales2020/api-gateway/config"
)

type AdminHandler struct {
	store config.RouteStore
}

func NewAdminHandler(store config.RouteStore) *AdminHandler {
	return &AdminHandler{
		store: store,
	}
}

// Routes registers admin endpoints
func (h *AdminHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/routes", h.GetRoutes)
	r.Post("/routes", h.CreateRoute)
	// Future: r.Delete("/routes/{path}", h.DeleteRoute)

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

	if err := h.store.SaveRoute(route); err != nil {
		http.Error(w, "Failed to save route", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Route saved",
	})
}
