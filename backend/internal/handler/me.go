package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrDisa/Raspy/backend/internal/auth"
	"github.com/mrDisa/Raspy/backend/internal/model"
	"github.com/mrDisa/Raspy/backend/internal/repository"
)

type UpdateGroupRequest struct {
    GroupID  int            `json:"group_id"`
    Subgroup model.Subgroup `json:"subgroup"`
}

func Me(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateGroup(userRepo *repository.UserRepository, groupRepo *repository.GroupRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.UserFromContext(r.Context())
		if !ok {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		var request UpdateGroupRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if request.GroupID <= 0 {
			http.Error(w, "invalid group_id", http.StatusBadRequest)
			return
		}

		if request.Subgroup != model.SubgroupAll &&
			request.Subgroup != model.SubgroupFirst &&
			request.Subgroup != model.SubgroupSecond {
			http.Error(w, "invalid subgroup", http.StatusBadRequest)
			return
		}

		group, err := groupRepo.FindByID(
			r.Context(),
			request.GroupID,
		)
		if err != nil {
			http.Error(w, "failed to find group", http.StatusInternalServerError)
			return
		}

		if group == nil {
			http.Error(w, "group not found", http.StatusNotFound)
			return
		}

		if err := userRepo.UpdateGroupAndSubgroup(
			r.Context(),
			user.ID,
			group.ID,
			request.Subgroup,
		); err != nil {
			http.Error(
				w,
				"failed to update group and subgroup",
				http.StatusInternalServerError,
			)
			return
		}

		user.GroupID = &group.ID
		user.Subgroup = request.Subgroup

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}