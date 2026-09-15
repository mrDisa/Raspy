package auth

import (
	"context"
	"net/http"

	"github.com/mrDisa/Raspy/backend/internal/model"
	"github.com/mrDisa/Raspy/backend/internal/repository"
)

type contextKey string

const userContextKey contextKey = "telegram_user"

func Middleware(userRepo *repository.UserRepository, botToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initData := r.Header.Get("X-Telegram-Init-Data")
			if initData == "" {
				http.Error(w, "missing init data", http.StatusUnauthorized)
				return
			}

			values, err := ValidateInitData(initData, botToken)
			if err != nil {
				http.Error(w, "invalid init data", http.StatusUnauthorized)
				return
			}

			tgUser, err := ParseTelegramUser(values)
			if err != nil {
				http.Error(w, "invalid user data", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.FindByTelegramID(tgUser.ID)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if user == nil {
				newUser := &model.User{
					TelegramID:           tgUser.ID,
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

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	return user, ok
}
