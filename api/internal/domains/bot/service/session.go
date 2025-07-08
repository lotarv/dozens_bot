package service

import (
	bot_types "github.com/lotarv/dozens_bot/internal/domains/bot/types/bot"
	user_types "github.com/lotarv/dozens_bot/internal/domains/user/types"
)

type userSession struct {
	User            user_types.User
	Dozen           bot_types.Dozen
	IsCreatingDozen bool
}

var sessions = make(map[int64]*userSession)

func getSession(userID int64) *userSession {
	if session, ok := sessions[userID]; ok {
		return session
	}
	session := &userSession{}
	sessions[userID] = session
	return session
}

func resetSession(userID int64) {
	delete(sessions, userID)
}
