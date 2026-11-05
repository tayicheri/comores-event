package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"comores-event-api/internal/models"
	"comores-event-api/internal/repository"
)

// BilletAdminHandler routes admin billets.
type BilletAdminHandler struct {
	billets repository.BilletRepository
}

// NewBilletAdminHandler construit le handler admin.
func NewBilletAdminHandler(billets repository.BilletRepository) *BilletAdminHandler {
	return &BilletAdminHandler{billets: billets}
}

// List : GET /api/admin/billets
func (h *BilletAdminHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.billets.List(r.Context())
	if err != nil {
		slog.Error("admin list billets", "error", err)
		WriteError(w, http.StatusInternalServerError, "Erreur interne du serveur")
		return
	}
	stats, err := h.billets.Stats(r.Context())
	if err != nil {
		slog.Error("admin stats billets", "error", err)
		WriteError(w, http.StatusInternalServerError, "Erreur interne du serveur")
		return
	}
	out := make([]models.BilletResponse, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToResponse())
	}
	WriteJSON(w, http.StatusOK, models.BilletListResponse{Billets: out, Stats: stats})
}

// Update : PATCH /api/admin/billets/{id}
func (h *BilletAdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "Identifiant manquant")
		return
	}
	var req models.UpdateBilletRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Status == nil && req.Comments == nil {
		WriteError(w, http.StatusBadRequest, "Aucune modification fournie")
		return
	}
	b, err := h.billets.Update(r.Context(), id, req.Status, req.Comments)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Billet introuvable")
			return
		}
		slog.Error("admin update billet", "error", err, "id", id)
		WriteError(w, http.StatusInternalServerError, "Erreur interne du serveur")
		return
	}
	WriteJSON(w, http.StatusOK, b.ToResponse())
}
