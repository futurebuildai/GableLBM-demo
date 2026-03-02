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
	mux.HandleFunc("GET /api/reports/ar-aging", h.HandleARAgingReport)
	mux.HandleFunc("GET /api/reports/customer-statement/{id}", h.HandleCustomerStatement)
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

func (h *Handler) HandleARAgingReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetARAgingReport(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *Handler) HandleCustomerStatement(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	if customerID == "" {
		http.Error(w, "customer ID required", http.StatusBadRequest)
		return
	}
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	stmt, err := h.service.GetCustomerStatement(r.Context(), customerID, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stmt)
}

func (h *Handler) RegisterBuilderRoutes(mux *http.ServeMux) {
mux.HandleFunc("POST /api/reporting/builder/preview", h.HandleBuilderPreview)
mux.HandleFunc("POST /api/reporting/builder/export", h.HandleBuilderExport)
mux.HandleFunc("POST /api/reporting/save", h.HandleSaveReport)
mux.HandleFunc("GET /api/reporting/saved", h.HandleListSavedReports)
mux.HandleFunc("GET /api/reporting/saved/{id}", h.HandleGetSavedReport)
mux.HandleFunc("PUT /api/reporting/saved/{id}", h.HandleUpdateSavedReport)
mux.HandleFunc("DELETE /api/reporting/saved/{id}", h.HandleDeleteSavedReport)
}

func (h *Handler) HandleBuilderPreview(w http.ResponseWriter, r *http.Request) {
var req struct {
EntityType string           `json:"entity_type"`
Definition ReportDefinition `json:"definition"`
}
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

results, err := h.service.ExecuteReportDefinition(r.Context(), &req.Definition, req.EntityType)
if err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(results)
}

func (h *Handler) HandleBuilderExport(w http.ResponseWriter, r *http.Request) {
var req struct {
EntityType string           `json:"entity_type"`
Format     string           `json:"format"` // csv, xlsx
Definition ReportDefinition `json:"definition"`
}
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

results, err := h.service.ExecuteReportDefinition(r.Context(), &req.Definition, req.EntityType)
if err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}

switch req.Format {
case "csv":
w.Header().Set("Content-Type", "text/csv")
w.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)
if err := ExportCSV(w, req.Definition.Columns, results); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
}
case "xlsx":
w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
w.Header().Set("Content-Disposition", `attachment; filename="report.xlsx"`)
if err := ExportXLSX(w, req.Definition.Columns, results); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
}
default:
http.Error(w, "unsupported format", http.StatusBadRequest)
}
}

func (h *Handler) HandleSaveReport(w http.ResponseWriter, r *http.Request) {
var report SavedReport
if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

// Assume user ID is extracted from context/auth middleware ideally
report.CreatedBy = "system" 

if err := h.service.CreateSavedReport(r.Context(), &report); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(report)
}

func (h *Handler) HandleListSavedReports(w http.ResponseWriter, r *http.Request) {
reports, err := h.service.ListSavedReports(r.Context())
if err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(reports)
}

func (h *Handler) HandleGetSavedReport(w http.ResponseWriter, r *http.Request) {
id := r.PathValue("id")
report, err := h.service.GetSavedReport(r.Context(), id)
if err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(report)
}

func (h *Handler) HandleUpdateSavedReport(w http.ResponseWriter, r *http.Request) {
id := r.PathValue("id")
var report SavedReport
if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}
report.ID = id

if err := h.service.UpdateSavedReport(r.Context(), &report); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleDeleteSavedReport(w http.ResponseWriter, r *http.Request) {
id := r.PathValue("id")
if err := h.service.DeleteSavedReport(r.Context(), id); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterBIIntegrationRoutes(mux *http.ServeMux) {
mux.HandleFunc("GET /api/v1/reporting/export/{entity}", h.HandleBIEntityExport)
}

func (h *Handler) HandleBIEntityExport(w http.ResponseWriter, r *http.Request) {
entity := r.PathValue("entity")

// Validate entity via existing schemata from builder
_, ok := entitySchemas[entity]
if !ok {
http.Error(w, "Invalid entity requested for BI export", http.StatusBadRequest)
return
}

// Create a "SELECT *" equivalent definition for the BI tool
def := &ReportDefinition{
Columns: []ReportColumn{},
}

for fieldName := range entitySchemas[entity] {
def.Columns = append(def.Columns, ReportColumn{
Field: fieldName,
Label: fieldName,
})
}

// Fetch raw data
results, err := h.service.ExecuteReportDefinition(r.Context(), def, entity)
if err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}

// Output as structured JSON dump
w.Header().Set("Content-Type", "application/json")
if err := json.NewEncoder(w).Encode(results); err != nil {
http.Error(w, "Failed to encode BI output", http.StatusInternalServerError)
}
}
