package auth

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func ParseTelegramUser(values url.Values) (TelegramUser, error) {
	rawUser := values.Get("user")

	if rawUser == "" {
		return TelegramUser{}, fmt.Errorf("user field is missing")
	}

	var user TelegramUser

	if err := json.Unmarshal([]byte(rawUser), &user); err != nil {
		return TelegramUser{}, fmt.Errorf("failed to parse telegram user: %w", err)
	}

	if user.ID == 0 {
		return TelegramUser{}, fmt.Errorf("telegram user id is missing")
	}

	return user, nil
}