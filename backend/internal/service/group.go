package service

import (
	"context"
	"fmt"

	"github.com/mrDisa/Raspy/backend/internal/collegeapi"
	"github.com/mrDisa/Raspy/backend/internal/model"
	"github.com/mrDisa/Raspy/backend/internal/repository"
)

type GroupService struct {
	collegeClient *collegeapi.Client
	groupRepo     *repository.GroupRepository
}

func NewGroupService(
	collegeClient *collegeapi.Client,
	groupRepo *repository.GroupRepository,
) *GroupService {
	return &GroupService{
		collegeClient: collegeClient,
		groupRepo:     groupRepo,
	}
}

func (s *GroupService) GetGroups(ctx context.Context) ([]model.Group, error) {
	groups, err := s.collegeClient.GetGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get groups from college api: %w", err)
	}

	result := make([]model.Group, 0, len(groups))

	for _, group := range groups {
		dbGroup, err := s.groupRepo.Upsert(
			ctx,
			group.ID,
			group.Text,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to save group: %w", err)
		}

		result = append(result, *dbGroup)
	}

	return result, nil
}