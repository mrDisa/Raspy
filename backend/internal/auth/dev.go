package auth

import (
	"context"
	"net/http"

	"github.com/mrDisa/Raspy/backend/internal/model"
	"github.com/mrDisa/Raspy/backend/internal/repository"
)

const devUserID int64 = 123456789

func DevMiddleware(
	userRepo *repository.UserRepository,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := userRepo.FindByTelegramID(r.Context(), devUserID)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if user == nil {
				newUser := &model.User{
					TelegramID:           devUserID,
					GroupID:              nil,
					Subgroup:             model.SubgroupAll,
					NotificationsEnabled: false,
				}

				user, err = userRepo.Create(r.Context(), newUser)
				if err != nil {
					http.Error(w, "internal error", http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(
				r.Context(),
				userContextKey,
				user,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}