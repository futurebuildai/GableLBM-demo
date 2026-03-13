package salesteam

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /sales-team", h.HandleList)
	mux.HandleFunc("GET /sales-team/{id}", h.HandleGet)
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	people, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch sales team", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid salesperson ID", http.StatusBadRequest)
		return
	}

	person, err := h.repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "Salesperson not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}
