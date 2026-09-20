package model

import (
	"time"
)

type User struct {
	ID 						int
	TelegramID 				int64
	GroupID					*int
	Group                   *Group
	Subgroup 				Subgroup
	NotificationsEnabled 	bool
	CreatedAt 				time.Time
	UpdatedAt 				time.Time
}