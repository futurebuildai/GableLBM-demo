package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Fleet
	mux.HandleFunc("GET /api/v1/delivery/vehicles", h.HandleListVehicles)
	mux.HandleFunc("POST /api/v1/delivery/vehicles", h.HandleCreateVehicle)
	mux.HandleFunc("GET /api/v1/delivery/drivers", h.HandleListDrivers)
	mux.HandleFunc("POST /api/v1/delivery/drivers", h.HandleCreateDriver)

	// Routes
	mux.HandleFunc("GET /api/v1/delivery/routes", h.HandleListRoutes)
	mux.HandleFunc("POST /api/v1/delivery/routes", h.HandleCreateRoute)
	mux.HandleFunc("POST /api/v1/delivery/routes/{id}/dispatch", h.HandleDispatchRoute)

	// Deliveries
	mux.HandleFunc("GET /api/v1/delivery/routes/{id}/deliveries", h.HandleListDeliveries)
	mux.HandleFunc("POST /api/v1/delivery/deliveries", h.HandleAssignOrder)                     // Assign Order to Route
	mux.HandleFunc("PUT /api/v1/delivery/deliveries/{id}/status", h.HandleUpdateDeliveryStatus) // Complete Delivery
}

// Fleet

func (h *Handler) HandleListVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.service.ListVehicles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func (h *Handler) HandleCreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	v, err := h.service.CreateVehicle(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) HandleListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := h.service.ListDrivers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

func (h *Handler) HandleCreateDriver(w http.ResponseWriter, r *http.Request) {
	var req CreateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	d, err := h.service.CreateDriver(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

// Routes

func (h *Handler) HandleListRoutes(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	var datePtr *string
	if dateStr != "" {
		datePtr = &dateStr
	}

	routes, err := h.service.ListRoutes(r.Context(), datePtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

func (h *Handler) HandleCreateRoute(w http.ResponseWriter, r *http.Request) {
	var req CreateRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	route, err := h.service.CreateRoute(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(route)
}

func (h *Handler) HandleDispatchRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.service.DispatchRoute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Deliveries

func (h *Handler) HandleListDeliveries(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // Route ID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	deliveries, err := h.service.ListDeliveries(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

func (h *Handler) HandleAssignOrder(w http.ResponseWriter, r *http.Request) {
	var req AssignOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	d, err := h.service.AssignOrderToRoute(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) HandleUpdateDeliveryStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var req UpdateDeliveryStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CompleteDelivery(r.Context(), id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
