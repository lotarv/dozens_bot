package repository

import (
	"context"

	"github.com/lotarv/dozens_bot/internal/domains/meetings/types"
	member_types "github.com/lotarv/dozens_bot/internal/domains/members/types"
	user_types "github.com/lotarv/dozens_bot/internal/domains/user/types"
	"github.com/lotarv/dozens_bot/internal/storage"
	global_types "github.com/lotarv/dozens_bot/internal/types"
)

type MeetingsRepository struct {
	*storage.Storage
	UsersRepository
}

type UsersRepository interface {
	GetUserByID(ctx context.Context, userID int64) (*user_types.User, error)
	GetMemberByUsername(username string) (member_types.Member, error)
	GetUserDozenByUsername(username string) (global_types.Dozen, error)
}

func New(storage *storage.Storage, userRepo UsersRepository) *MeetingsRepository {
	return &MeetingsRepository{
		Storage:         storage,
		UsersRepository: userRepo,
	}
}

func (r *MeetingsRepository) GetDozenMeetings(dozen_id int) ([]types.Meeting, error) {
	var meetings []types.Meeting
	err := r.DB().Select(&meetings, `SELECT * FROM meetings WHERE dozen_id=$1 ORDER BY start_time ASC`, dozen_id)
	if err != nil {
		return nil, err
	}
	return meetings, nil
}
