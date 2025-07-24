package service

import (
	"context"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TransactionStep int

const (
	Idle TransactionStep = iota
	AwaitingMember
	AwaitingAmount
	AwaitingReason
)

type TransactionSession struct {
	Step           TransactionStep
	IsDeposit      bool
	IsPenalty      bool
	Amount         int
	MemberUsername string
}

var transactionSessions = make(map[int64]*TransactionSession)

func (s *BotService) StartNewTransaction(userID int64) {
	msg := tgbotapi.NewMessage(userID, "Что хотим сделать?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить штраф", "start_penalty"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить (не штраф)", "start_deposit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Потратить", "start_withdraw"),
		),
	)
	s.bot.Send(msg)
}

func (s *BotService) handleTransactionStep(msg *tgbotapi.Message, session *TransactionSession) {
	userID := msg.From.ID
	text := strings.TrimSpace(msg.Text)

	//TODO: вычислять десятку пользователя по его ID
	switch session.Step {
	case AwaitingAmount:
		amount, err := strconv.Atoi(text)
		if err != nil || amount <= 0 {
			s.bot.Send(tgbotapi.NewMessage(userID, "Введите корректную положительную сумму"))
			return
		}
		if !(session.IsDeposit || session.IsPenalty) {
			amount = -amount
		}
		session.Amount = amount
		session.Step = AwaitingReason
		if session.IsPenalty {
			s.bot.Send(tgbotapi.NewMessage(userID, "За что штраф?"))
		} else if session.IsDeposit {
			s.bot.Send(tgbotapi.NewMessage(userID, "Укажите причину"))
		} else {
			s.bot.Send(tgbotapi.NewMessage(userID, "Укажите цель списания средств:"))
		}
	case AwaitingReason:
		reason := text

		err := s.repo.ChangeBankBalance(context.Background(), 1, session.Amount, reason, session.MemberUsername)
		if err != nil {
			s.bot.Send(tgbotapi.NewMessage(userID, "Ошибка при изменении баланса"))
		} else {
			s.bot.Send(tgbotapi.NewMessage(userID, "✅ Информация внесена"))
		}
		// 🧹 Чистим
		delete(transactionSessions, userID)
	}
}

func (s *BotService) askForTransactionMember(userID int64) {
	//TODO: вычислять участников десятки человека
	members, err := s.repo.GetMembers()
	if err != nil {
		s.bot.Send(tgbotapi.NewMessage(userID, "Не удалось получить список пользователей. Попробуйте позже"))
		// 🧹 Чистим
		delete(transactionSessions, userID)
		return
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, m := range members {
		btn := tgbotapi.NewInlineKeyboardButtonData(m.FullName, "select_member_"+m.Username)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg := tgbotapi.NewMessage(userID, "Чей штраф?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	s.bot.Send(msg)

}
