package service

import (
	"context"
	"log/slog"

	"github.com/lotarv/dozens_bot/internal/domains/meetings/types"
	member_types "github.com/lotarv/dozens_bot/internal/domains/members/types"
	user_types "github.com/lotarv/dozens_bot/internal/domains/user/types"
	global_types "github.com/lotarv/dozens_bot/internal/types"
)

type repository interface {
	GetDozenMeetings(dozen_id int) ([]types.Meeting, error)
	GetUserByID(ctx context.Context, userID int64) (*user_types.User, error)
	GetMemberByUsername(username string) (member_types.Member, error)
	GetUserDozenByUsername(username string) (global_types.Dozen, error)
}

type MeetingsService struct {
	repo repository
}

func New(repo repository) *MeetingsService {
	return &MeetingsService{
		repo: repo,
	}
}

func (s *MeetingsService) GetDozenMeetings(ctx context.Context, userID int64) ([]types.Meeting, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		return []types.Meeting{}, err
	}

	dozen, err := s.repo.GetUserDozenByUsername(user.Username)
	if err != nil {
		slog.Error("failed to get user dozen by username", "error", err)
		return []types.Meeting{}, err
	}

	meetings, err := s.repo.GetDozenMeetings(dozen.ID)
	if err != nil {
		slog.Error("failed to get dozen meetings", "error", err)
		return []types.Meeting{}, err
	}

	return meetings, nil
}
