package gl

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Handler exposes GL REST endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates a new GL Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all GL routes on the mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Accounts
	mux.HandleFunc("GET /gl/accounts", h.HandleListAccounts)
	mux.HandleFunc("POST /gl/accounts", h.HandleCreateAccount)
	mux.HandleFunc("PUT /gl/accounts/{id}", h.HandleUpdateAccount)

	// Journal Entries
	mux.HandleFunc("GET /gl/journal-entries", h.HandleListJournalEntries)
	mux.HandleFunc("GET /gl/journal-entries/{id}", h.HandleGetJournalEntry)
	mux.HandleFunc("POST /gl/journal-entries", h.HandleCreateJournalEntry)
	mux.HandleFunc("POST /gl/journal-entries/{id}/post", h.HandlePostJournalEntry)
	mux.HandleFunc("POST /gl/journal-entries/{id}/void", h.HandleVoidJournalEntry)

	// Trial Balance
	mux.HandleFunc("GET /gl/trial-balance", h.HandleTrialBalance)

	// Fiscal Periods
	mux.HandleFunc("GET /gl/fiscal-periods", h.HandleListFiscalPeriods)
	mux.HandleFunc("POST /gl/fiscal-periods/{id}/close", h.HandleCloseFiscalPeriod)
}

// --- Account Handlers ---

func (h *Handler) HandleListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.svc.ListAccounts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

type createAccountRequest struct {
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Subtype       string     `json:"subtype"`
	ParentID      *uuid.UUID `json:"parent_id,omitempty"`
	NormalBalance string     `json:"normal_balance"`
	Description   string     `json:"description"`
}

func (h *Handler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	acct := &GLAccount{
		Code:          req.Code,
		Name:          req.Name,
		Type:          req.Type,
		Subtype:       req.Subtype,
		ParentID:      req.ParentID,
		NormalBalance: req.NormalBalance,
		Description:   req.Description,
	}

	if err := h.svc.CreateAccount(r.Context(), acct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acct)
}

func (h *Handler) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid account ID", http.StatusBadRequest)
		return
	}

	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	acct := &GLAccount{
		ID:            id,
		Code:          req.Code,
		Name:          req.Name,
		Type:          req.Type,
		Subtype:       req.Subtype,
		ParentID:      req.ParentID,
		NormalBalance: req.NormalBalance,
		Description:   req.Description,
		IsActive:      true,
	}

	if err := h.svc.UpdateAccount(r.Context(), acct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acct)
}

// --- Journal Entry Handlers ---

func (h *Handler) HandleListJournalEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.ListJournalEntries(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (h *Handler) HandleGetJournalEntry(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid journal entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.svc.GetJournalEntry(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

type createJournalEntryRequest struct {
	EntryDate string           `json:"entry_date"` // YYYY-MM-DD
	Memo      string           `json:"memo"`
	Lines     []journalLineReq `json:"lines"`
}

type journalLineReq struct {
	AccountID   string  `json:"account_id"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}

func (h *Handler) HandleCreateJournalEntry(w http.ResponseWriter, r *http.Request) {
	var req createJournalEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	entryDate := time.Now()
	if req.EntryDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EntryDate)
		if err != nil {
			http.Error(w, "invalid entry_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		entryDate = parsed
	}

	var lines []JournalLine
	for _, lr := range req.Lines {
		accountID, err := uuid.Parse(lr.AccountID)
		if err != nil {
			http.Error(w, "invalid account_id: "+lr.AccountID, http.StatusBadRequest)
			return
		}
		lines = append(lines, JournalLine{
			AccountID:   accountID,
			Description: lr.Description,
			Debit:       int64(lr.Debit*100 + 0.5),
			Credit:      int64(lr.Credit*100 + 0.5),
		})
	}

	entry := &JournalEntry{
		EntryDate: entryDate,
		Memo:      req.Memo,
		Source:    SourceManual,
		Status:    StatusDraft,
		Lines:     lines,
	}

	if err := h.svc.CreateJournalEntry(r.Context(), entry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func (h *Handler) HandlePostJournalEntry(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid journal entry ID", http.StatusBadRequest)
		return
	}

	if err := h.svc.PostJournalEntry(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "posted"})
}

func (h *Handler) HandleVoidJournalEntry(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid journal entry ID", http.StatusBadRequest)
		return
	}

	if err := h.svc.VoidJournalEntry(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "voided"})
}

// --- Trial Balance ---

func (h *Handler) HandleTrialBalance(w http.ResponseWriter, r *http.Request) {
	asOfStr := r.URL.Query().Get("as_of")
	asOf := time.Now()
	if asOfStr != "" {
		parsed, err := time.Parse("2006-01-02", asOfStr)
		if err != nil {
			http.Error(w, "invalid as_of date (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		asOf = parsed
	}

	rows, err := h.svc.GetTrialBalance(r.Context(), asOf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

// --- Fiscal Periods ---

func (h *Handler) HandleListFiscalPeriods(w http.ResponseWriter, r *http.Request) {
	periods, err := h.svc.ListFiscalPeriods(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(periods)
}

func (h *Handler) HandleCloseFiscalPeriod(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid fiscal period ID", http.StatusBadRequest)
		return
	}

	if err := h.svc.CloseFiscalPeriod(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "closed"})
}
