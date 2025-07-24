package service

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var reportBuffers = map[string]*reportBuffer{}

type reportBuffer struct {
	text        string
	username    string
	timer       *time.Timer
	originalMsg *tgbotapi.Message
	forwardDate int
}

func (s *BotService) flushBufferedReport(key string, userID int64) {
	buf, ok := reportBuffers[key]
	if !ok {
		return
	}
	delete(reportBuffers, key)

	dummyMsg := &tgbotapi.Message{
		MessageID:   buf.originalMsg.MessageID,
		Chat:        buf.originalMsg.Chat,
		From:        buf.originalMsg.From,
		Text:        buf.text,
		Date:        buf.originalMsg.Date,
		ForwardDate: buf.forwardDate,
	}
	s.handleReport(dummyMsg)
}
