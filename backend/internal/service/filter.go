package service

import "github.com/mrDisa/Raspy/backend/internal/model"

func filterBySubgroup(lessons []model.Lesson, subgroup model.Subgroup) []model.Lesson {
	filteredLessons := make([]model.Lesson, 0)
	for _, lesson := range lessons {
		if lesson.Subgroup == model.SubgroupAll || lesson.Subgroup == subgroup {
			filteredLessons = append(filteredLessons, lesson)
		}
	}
	return filteredLessons
}
