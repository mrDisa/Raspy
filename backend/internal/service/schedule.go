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

func NewScheduleService(provider ScheduleProvider) *ScheduleService {
	return &ScheduleService{
		provider: provider,
	}
}

func (s *ScheduleService) getFilteredSchedule(group string, subgroup model.Subgroup, startDate *time.Time) ([]model.ScheduleDay, error) {
	schedule, err := s.provider.GetSchedule(group, startDate)
	if err != nil {
		return []model.ScheduleDay{}, fmt.Errorf("failed to get schedule: %w", err)
	}
	for i := range schedule {
		schedule[i].List = filterBySubgroup(schedule[i].List, subgroup)
	}
	return schedule, err
}

func (s *ScheduleService) getScheduleForDate(group string, subgroup model.Subgroup, date time.Time) (model.ScheduleDay, error) {
	schedule, err := s.getFilteredSchedule(group, subgroup, nil)
	if err != nil {
		return model.ScheduleDay{}, fmt.Errorf("failed to get schedule: %w", err)
	}

	targetDate := date.Format("2006-01-02")
	for _, day := range schedule {
		if day.Date == targetDate {
			return day, nil
		}
	}
	return model.ScheduleDay{}, fmt.Errorf("schedule for %s not found", targetDate)
}

func (s *ScheduleService) GetTodaySchedule(group string, subgroup model.Subgroup) (model.ScheduleDay, error) {
	return s.getScheduleForDate(group, subgroup, time.Now())
}
func (s *ScheduleService) GetTomorrowSchedule(group string, subgroup model.Subgroup) (model.ScheduleDay, error) {
	return s.getScheduleForDate(group, subgroup, time.Now().AddDate(0, 0, 1))
}

func (s *ScheduleService) GetNextWeekSchedule(group string, subgroup model.Subgroup) ([]model.ScheduleDay, error) {
	weekday := int(time.Now().Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysUntilNextMonday := 8 - weekday
	rawDate := time.Now().AddDate(0, 0, daysUntilNextMonday)

	loc, err := time.LoadLocation(("Europe/Moscow"))
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	date := time.Date(rawDate.Year(), rawDate.Month(), rawDate.Day(), 0, 0, 0, 0, loc)

	fmt.Println(date)
	return s.getFilteredSchedule(group, subgroup, &date)
}
