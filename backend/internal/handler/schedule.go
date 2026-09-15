package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrDisa/Raspy/backend/internal/auth"
	"github.com/mrDisa/Raspy/backend/internal/repository"
	"github.com/mrDisa/Raspy/backend/internal/service"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
	groupRepo      *repository.GroupRepository
}

func NewScheduleHandler(
	scheduleService *service.ScheduleService,
	groupRepo *repository.GroupRepository,
) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
		groupRepo:       groupRepo,
	}
}

func (h *ScheduleHandler) Today(w http.ResponseWriter, r *http.Request) {
	h.getSchedule(w, r, "today")
}

func (h *ScheduleHandler) Tomorrow(w http.ResponseWriter, r *http.Request) {
	h.getSchedule(w, r, "tomorrow")
}

func (h *ScheduleHandler) getSchedule(w http.ResponseWriter, r *http.Request, day string) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	if user.GroupID == nil {
		http.Error(w, "group is not selected", http.StatusBadRequest)
		return
	}

	group, err := h.groupRepo.FindByID(
		r.Context(),
		*user.GroupID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get group",
			http.StatusInternalServerError,
		)
		return
	}

	if group == nil {
		http.Error(
			w,
			"group not found",
			http.StatusNotFound,
		)
		return
	}

	var schedule interface{}

	switch day {
	case "today":
		schedule, err = h.scheduleService.GetTodaySchedule(
			group.ExternalID,
			user.Subgroup,
		)

	case "tomorrow":
		schedule, err = h.scheduleService.GetTomorrowSchedule(
			group.ExternalID,
			user.Subgroup,
		)

	default:
		http.Error(
			w,
			"invalid schedule day",
			http.StatusBadRequest,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"failed to get schedule",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(schedule); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}