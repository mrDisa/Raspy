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

func moscowNow() (time.Time, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"failed to load timezone: %w",
			err,
		)
	}

	return time.Now().In(loc), nil
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
	loc, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        return model.ScheduleDay{}, fmt.Errorf(
            "failed to load timezone: %w",
            err,
        )
    }
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)

	schedule, err := s.getFilteredSchedule(group, subgroup, &targetDate)
	if err != nil {
		return model.ScheduleDay{}, fmt.Errorf("failed to get schedule: %w", err)
	}

	targetDateString := targetDate.Format("2006-01-02")
	for _, day := range schedule {
		if day.Date == targetDateString {
			return day, nil
		}
	}
	return model.ScheduleDay{Date: targetDateString, List: []model.Lesson{}}, nil
}

func (s *ScheduleService) GetTodaySchedule(group string, subgroup model.Subgroup) (model.ScheduleDay, error) {
	now, err := moscowNow()
	if err != nil {
		return model.ScheduleDay{}, err
	}

	return s.getScheduleForDate(group, subgroup, now)
}
func (s *ScheduleService) GetTomorrowSchedule(group string, subgroup model.Subgroup) (model.ScheduleDay, error) {
	now, err := moscowNow()
	if err != nil {
		return model.ScheduleDay{}, err
	}

	return s.getScheduleForDate(group, subgroup, now.AddDate(0, 0, 1))
}

func (s *ScheduleService) GetWeekSchedule(
	group string,
	subgroup model.Subgroup,
	weekOffset int,
) ([]model.ScheduleDay, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	now := time.Now().In(loc)

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	monday := now.AddDate(0, 0, -(weekday - 1))

	targetMonday := monday.AddDate(
		0,
		0,
		weekOffset*7,
	)

	startDate := time.Date(
		targetMonday.Year(),
		targetMonday.Month(),
		targetMonday.Day(),
		0,
		0,
		0,
		0,
		loc,
	)

	return s.getFilteredSchedule(
		group,
		subgroup,
		&startDate,
	)
}

