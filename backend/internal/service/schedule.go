package service

import (
	"fmt"
	"time"

	"github.com/mrDisa/Raspy/backend/internal/model"
)

type ScheduleProvider interface {
	GetSchedule(group string, startDate *time.Time) ([]model.ScheduleDay, error)
}

type ScheduleService struct {
    provider ScheduleProvider
}

func NewScheduleService(provider ScheduleProvider) ScheduleService {
    return ScheduleService{
		provider: provider,
	}
}

func (s *ScheduleService) GetTodaySchedule(group string) (model.ScheduleDay, error) {
	today := time.Now()
	schedule, err := s.provider.GetSchedule(group, nil)
	if err != nil {
		return model.ScheduleDay{}, fmt.Errorf("failed to get schedule: %v", err)
	}
	dateToday := today.Format("2006-01-02")
	for _, day := range schedule {
		if day.Date == dateToday {
			return day, nil
		}
	}
	return  model.ScheduleDay{}, fmt.Errorf("schedule for today not found")
}