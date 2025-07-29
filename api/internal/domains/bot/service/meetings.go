package service

import (
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *BotService) handleMeeting(msg *tgbotapi.Message) {
	lines := strings.Split(strings.TrimSpace(msg.Text), "\n")

	if len(lines) < 4 {
		s.replyTo(msg, "Формат встречи неверный. Нужно минимум 4 строки.")
		return
	}

	dateTimeLine := strings.TrimSpace(lines[1]) // "02.09 13:00-20:00" или "02.09.25 13:00-20:00"
	dateTimeParts := strings.Split(dateTimeLine, " ")
	if len(dateTimeParts) != 2 {
		s.replyTo(msg, "Не удалось распарсить дату и время.")
		return
	}

	dateStr := dateTimeParts[0]
	timeStr := dateTimeParts[1]

	// Определяем формат даты
	var dateParsed time.Time
	var err error
	switch {
	case strings.Count(dateStr, ".") == 2:
		// указан год: 02.09.25
		dateParsed, err = time.Parse("02.01.06", dateStr)
	case strings.Count(dateStr, ".") == 1:
		// без года: 02.09 — подставляем текущий
		dateParsed, err = time.Parse("02.01", dateStr)
		if err == nil {
			dateParsed = dateParsed.AddDate(time.Now().Year(), 0, 0)
		}
	default:
		s.replyTo(msg, "Формат даты должен быть ДД.ММ или ДД.ММ.ГГ")
		return
	}

	if err != nil {
		s.replyTo(msg, "Ошибка при разборе даты: "+err.Error())
		return
	}

	// Время
	timeRange := strings.Split(timeStr, "-")
	if len(timeRange) < 1 {
		s.replyTo(msg, "Неверный формат времени. Пример: 13:00-20:00")
		return
	}

	startTimeOnly, err := time.Parse("15:04", timeRange[0])
	if err != nil {
		s.replyTo(msg, "Ошибка при разборе начала времени")
		return
	}

	startTime := time.Date(
		dateParsed.Year(), dateParsed.Month(), dateParsed.Day(),
		startTimeOnly.Hour(), startTimeOnly.Minute(), 0, 0, time.Local,
	)
	startTimeStr := startTime.Format("2006-01-02 15:04:05")

	var endTimeStr *string
	if len(timeRange) == 2 {
		endTimeOnly, err := time.Parse("15:04", timeRange[1])
		if err == nil {
			endTime := time.Date(
				dateParsed.Year(), dateParsed.Month(), dateParsed.Day(),
				endTimeOnly.Hour(), endTimeOnly.Minute(), 0, 0, time.Local,
			)
			formatted := endTime.Format("2006-01-02 15:04:05")
			endTimeStr = &formatted
		}
	}

	// Место и карта
	locationName := strings.TrimSpace(lines[2])
	mapURL := strings.Trim(lines[3], "\" ")

	// TODO: Получать ID десятки пользователя

	// Создаём встречу
	err = s.repo.CreateMeeting(1, startTimeStr, endTimeStr, locationName, mapURL)
	if err != nil {
		s.replyTo(msg, "Ошибка при сохранении встречи.")
		slog.Error("Failed to create meeting: ", "error", err, "startTimeStr", startTimeStr, "endTimeStr", endTimeStr)
		return
	}

	s.replyTo(msg, "Встреча успешно добавлена ✅")
}
