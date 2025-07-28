package service

import (
<<<<<<< HEAD
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
=======
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/jomei/notionapi"
	"github.com/lotarv/dozens_bot/internal/domains/bot/helpers"
	"github.com/lotarv/dozens_bot/internal/utils/crypto"
)

func (s *BotService) createDocument(id string) error {
	properties := notionapi.Properties{
		"ID": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{Content: id},
				},
			},
		},
	}

	if os.Getenv("ENV") == "DEV" {
		properties["Текст"] = notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{Content: "Тестовый документ"},
				},
			},
		}
	}

	page := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(s.notionConfig.DocumentsDBID),
		},
		Properties: properties,
	}

	_, err := s.notionClient.Page.Create(context.Background(), page)
	return err
}

func (s *BotService) createReport(documentID, authorID, date string) error {
	page := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(s.notionConfig.ReportsDBID),
		},
		Properties: notionapi.Properties{
			"ID": notionapi.RelationProperty{
				Relation: []notionapi.Relation{
					{ID: notionapi.PageID(documentID)},
				},
			},
			"Автор": notionapi.RelationProperty{
				Relation: []notionapi.Relation{
					{ID: notionapi.PageID(authorID)},
				},
			},
			"Дата создания": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{Content: date},
					},
				},
			},
		},
	}

	_, err := s.notionClient.Page.Create(context.Background(), page)
	return err
}

func (s *BotService) replyTo(msg *tgbotapi.Message, text string) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyToMessageID = msg.MessageID
	s.bot.Send(reply)
}

func (s *BotService) handleReport(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	// Определяем автора
	var username string
	if msg.ForwardFrom != nil && msg.ForwardFrom.UserName != "" {
		username = msg.ForwardFrom.UserName
	} else if msg.From != nil && msg.From.UserName != "" {
		// username = msg.From.UserName
		username = "incetro"
	} else {
		slog.Warn("Report without username", "chat_id", chatID)
		s.replyTo(msg, "Ошибка при определении автора отчета")
		return
	}

	// Нормализуем текст
	text := strings.ToLower(strings.ReplaceAll(msg.Text, "ё", "е"))
	lines := strings.Split(text, "\n")
	if len(lines) < 1 {
		s.replyTo(msg, "Неверный формат отчета: отсутствует содержимое")
		return
	}

	// Достаем тело отчета: все, кроме строки с хештегом
	reportText := helpers.ExtractReportBody(msg.Text)
	var reportTime string
	if msg.ForwardDate != 0 {
		reportTime = time.Unix(int64(msg.ForwardDate), 0).Format("02/01/2006")
	} else {
		reportTime = time.Unix(int64(msg.Date), 0).Format("02/01/2006")
	}
	uuidStr := uuid.New().String()

	slog.Info("New report", "user", username, "uuid", uuidStr, "text", reportText, "time", reportTime)

	// 1. Получаем Notion ID автора
	authorNotionID, err := s.repo.GetMemberNotionId(username)
	if err != nil {
		slog.Error("Failed to get author notion ID", "err", err)
		s.replyTo(msg, "Ошибка при определении автора отчета")
		return
	}

	//1.1 Шифруем данные
	//TODO: брать код десятки текущего пользователя
	pepper := os.Getenv("ENCRYPTION_PEPPER")
	dozenCode := os.Getenv("ONLY_DOZEN_CODE")
	passphrase := dozenCode + pepper

	encryptedText, err := crypto.Encrypt(reportText, passphrase)
	if err != nil {
		slog.Error("failed to encrypt reportText", "err", err)
		s.replyTo(msg, "Не удалось выполнить шифрование отчета. Попробуйте еще раз")
		return
	}

	//2.Создание документа
	if err := s.createDocument(uuidStr); err != nil {
		slog.Error("Failed to create document", "err", err)
		s.replyTo(msg, "Ошибка при создании документа в Notion")
		return
	}

	//3.Синхронизация документов
	if err := helpers.TriggerSyncDocuments(); err != nil {
		slog.Error("Failed to sync documents", "err", err)
	}

	//4. Подставляем текст в БД

	if err := s.repo.SetEncryptedText(uuidStr, encryptedText); err != nil {
		slog.Error("failed to set encrypted text: ", "error", err)
		s.replyTo(msg, "Ошибка при создании документа в базе")
		return
	}
	//5.Получаем Notion ID документа
	docNotionID, err := s.repo.GetDocumentNotionId(uuidStr)
	if err != nil {
		slog.Error("Failed to get document notion ID", "err", err)
		s.replyTo(msg, "Ошибка при получении notion-id документа")
		return
	}

	// 6. Создание отчёта
	if err := s.createReport(docNotionID, authorNotionID, reportTime); err != nil {
		slog.Error("Failed to create report", "err", err)
		s.replyTo(msg, "Ошибка при создании отчета в notion")

		return
	}

	// 7. Синхронизация отчётов
	if err := helpers.TriggerSyncReports(); err != nil {
		slog.Error("Failed to sync reports", "err", err)
		s.replyTo(msg, "Отчет создан в notion, но синхронизация не удалась")
	}

	s.replyTo(msg, "Отчёт успешно принят ✅")
>>>>>>> 260dcfe (refactor(backend): reports refactoring)
}
