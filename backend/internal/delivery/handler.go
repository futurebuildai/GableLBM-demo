package delivery

import (
	"encoding/json"
	"log/slog"
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
	mux.HandleFunc("GET /api/v1/delivery/vehicles/{id}", h.HandleGetVehicle)
	mux.HandleFunc("PUT /api/v1/delivery/vehicles/{id}", h.HandleUpdateVehicle)
	mux.HandleFunc("DELETE /api/v1/delivery/vehicles/{id}", h.HandleDeleteVehicle)
	mux.HandleFunc("GET /api/v1/delivery/drivers", h.HandleListDrivers)
	mux.HandleFunc("POST /api/v1/delivery/drivers", h.HandleCreateDriver)
	mux.HandleFunc("GET /api/v1/delivery/drivers/{id}", h.HandleGetDriver)
	mux.HandleFunc("PUT /api/v1/delivery/drivers/{id}", h.HandleUpdateDriver)
	mux.HandleFunc("DELETE /api/v1/delivery/drivers/{id}", h.HandleDeleteDriver)

	// Routes
	mux.HandleFunc("GET /api/v1/delivery/routes", h.HandleListRoutes)
	mux.HandleFunc("POST /api/v1/delivery/routes", h.HandleCreateRoute)
	mux.HandleFunc("POST /api/v1/delivery/routes/{id}/dispatch", h.HandleDispatchRoute)
	mux.HandleFunc("POST /api/v1/delivery/routes/{id}/reorder", h.HandleReorderStops)
	mux.HandleFunc("POST /api/v1/delivery/routes/{id}/optimize", h.HandleOptimizeRoute)
	mux.HandleFunc("POST /api/v1/delivery/routes/{id}/complete", h.HandleCompleteRoute)

	// Deliveries
	mux.HandleFunc("GET /api/v1/delivery/routes/{id}/deliveries", h.HandleListDeliveries)
	mux.HandleFunc("GET /api/v1/delivery/deliveries/{id}", h.HandleGetDelivery)
	mux.HandleFunc("POST /api/v1/delivery/deliveries", h.HandleAssignOrder)                     // Assign Order to Route
	mux.HandleFunc("PUT /api/v1/delivery/deliveries/{id}/status", h.HandleUpdateDeliveryStatus) // Complete Delivery
	mux.HandleFunc("POST /api/v1/delivery/deliveries/{id}/adjust-qty", h.HandleAdjustQuantity)
}

// Fleet

func (h *Handler) HandleListVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.service.ListVehicles(r.Context())
	if err != nil {
		slog.Error("ListVehicles failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		slog.Error("CreateVehicle failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) HandleListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := h.service.ListDrivers(r.Context())
	if err != nil {
		slog.Error("ListDrivers failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		slog.Error("CreateDriver failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	driverIDStr := r.URL.Query().Get("driver_id")
	var driverID *uuid.UUID
	if driverIDStr != "" {
		id, err := uuid.Parse(driverIDStr)
		if err != nil {
			http.Error(w, "Invalid driver_id UUID", http.StatusBadRequest)
			return
		}
		driverID = &id
	}

	routes, err := h.service.ListRoutes(r.Context(), datePtr, driverID)
	if err != nil {
		slog.Error("ListRoutes failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		slog.Error("CreateRoute failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		slog.Error("DispatchRoute failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleReorderStops(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var req ReorderStopsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ReorderStops(r.Context(), id, req.OrderedDeliveryIDs); err != nil {
		slog.Error("ReorderStops failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		slog.Error("ListDeliveries failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

func (h *Handler) HandleGetDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	d, err := h.service.GetDelivery(r.Context(), id)
	if err != nil {
		slog.Error("GetDelivery failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) HandleAssignOrder(w http.ResponseWriter, r *http.Request) {
	var req AssignOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	d, capacityWarning, err := h.service.AssignOrderToRoute(r.Context(), req)
	if err != nil {
		slog.Error("AssignOrderToRoute failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Delivery        *Delivery        `json:"delivery"`
		CapacityWarning *CapacityWarning `json:"capacity_warning,omitempty"`
	}{
		Delivery:        d,
		CapacityWarning: capacityWarning,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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
		slog.Error("CompleteDelivery failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleOptimizeRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	result, err := h.service.OptimizeRoute(r.Context(), id)
	if err != nil {
		slog.Error("OptimizeRoute failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) HandleGetVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	v, err := h.service.GetVehicle(r.Context(), id)
	if err != nil {
		slog.Error("GetVehicle failed", "error", err)
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) HandleUpdateVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	var req UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	v, err := h.service.UpdateVehicle(r.Context(), id, req)
	if err != nil {
		slog.Error("UpdateVehicle failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) HandleDeleteVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteVehicle(r.Context(), id); err != nil {
		slog.Error("DeleteVehicle failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleGetDriver(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	d, err := h.service.GetDriver(r.Context(), id)
	if err != nil {
		slog.Error("GetDriver failed", "error", err)
		http.Error(w, "Driver not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) HandleUpdateDriver(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	var req UpdateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	d, err := h.service.UpdateDriver(r.Context(), id, req)
	if err != nil {
		slog.Error("UpdateDriver failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) HandleDeleteDriver(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteDriver(r.Context(), id); err != nil {
		slog.Error("DeleteDriver failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleCompleteRoute(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	if err := h.service.CompleteRoute(r.Context(), id); err != nil {
		slog.Error("CompleteRoute failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
}

func (h *Handler) HandleAdjustQuantity(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var req QtyAdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.DeliveryID = deliveryID

	if err := h.service.AdjustDeliveryQuantity(r.Context(), req); err != nil {
		slog.Error("AdjustDeliveryQuantity failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "adjusted"})
}
