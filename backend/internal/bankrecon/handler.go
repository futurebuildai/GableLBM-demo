package bankrecon

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Handler handles bank reconciliation HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new bank reconciliation handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers bank reconciliation API routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Bank Accounts
	mux.HandleFunc("POST /api/bankrecon/accounts", h.CreateBankAccount)
	mux.HandleFunc("GET /api/bankrecon/accounts", h.ListBankAccounts)

	// CSV Import
	mux.HandleFunc("POST /api/bankrecon/import", h.ImportCSV)

	// Reconciliation Sessions
	mux.HandleFunc("POST /api/bankrecon/sessions", h.CreateSession)
	mux.HandleFunc("GET /api/bankrecon/sessions", h.ListSessions)
	mux.HandleFunc("GET /api/bankrecon/sessions/{id}", h.GetSession)
	mux.HandleFunc("POST /api/bankrecon/sessions/{id}/complete", h.CompleteSession)

	// Manual Match/Unmatch
	mux.HandleFunc("POST /api/bankrecon/match", h.ManualMatch)
	mux.HandleFunc("POST /api/bankrecon/unmatch", h.ManualUnmatch)
}

func (h *Handler) CreateBankAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateBankAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	acct, err := h.service.CreateBankAccount(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acct)
}

func (h *Handler) ListBankAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListBankAccounts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if accounts == nil {
		accounts = []BankAccount{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *Handler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	var req ImportCSVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.ImportCSV(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, err := h.service.CreateSession(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	session, err := h.service.GetSession(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	var bankAccountID *uuid.UUID
	if bid := r.URL.Query().Get("bank_account_id"); bid != "" {
		parsed, err := uuid.Parse(bid)
		if err == nil {
			bankAccountID = &parsed
		}
	}

	sessions, err := h.service.ListSessions(r.Context(), bankAccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sessions == nil {
		sessions = []ReconciliationSession{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func (h *Handler) CompleteSession(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	session, err := h.service.CompleteSession(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (h *Handler) ManualMatch(w http.ResponseWriter, r *http.Request) {
	var req ManualMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ManualMatch(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "matched"})
}

func (h *Handler) ManualUnmatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BankTransactionID uuid.UUID `json:"bank_transaction_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ManualUnmatch(r.Context(), req.BankTransactionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "unmatched"})
}
