package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrDisa/Raspy/backend/internal/service"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	groups, err := h.groupService.GetGroups(r.Context(), query)
	if err != nil {
		http.Error(w, "failed to get groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}