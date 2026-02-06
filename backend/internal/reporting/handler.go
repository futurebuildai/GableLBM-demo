package reporting

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/daily-till", h.HandleDailyTill)
	mux.HandleFunc("GET /api/reports/sales-summary", h.HandleSalesSummary)
}

func (h *Handler) HandleDailyTill(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	report, err := h.service.GetDailyTill(r.Context(), dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *Handler) HandleSalesSummary(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	report, err := h.service.GetSalesSummary(r.Context(), start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
