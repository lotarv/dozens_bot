package service

import (
	"context"
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/jomei/notionapi"
	"github.com/lotarv/dozens_bot/internal/domains/bot/helpers"
	bot_types "github.com/lotarv/dozens_bot/internal/domains/bot/types/bot"
)

type repository interface {
	GetDocumentsAmount() (int, error)
	GetDocumentNotionId(db_uuid string) (string, error)
	GetMemberNotionId(username string) (string, error)
}

type BotService struct {
	repo         repository
	bot          *tgbotapi.BotAPI
	notionClient *notionapi.Client
	notionConfig bot_types.NotionConfig
}

func New(repo repository, bot *tgbotapi.BotAPI, notionClient *notionapi.Client, config bot_types.NotionConfig) *BotService {

	return &BotService{
		repo:         repo,
		bot:          bot,
		notionClient: notionClient,
		notionConfig: config,
	}
}

func (s *BotService) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			s.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			s.handleCallback(update.CallbackQuery)
		}

	}
}

func (s *BotService) handleMessage(msg *tgbotapi.Message) {
	text := strings.ToLower(msg.Text)

	switch {
	case text == "/start":
		s.handleStart(msg)
	case strings.Contains(text, "#отчет"):
		s.handleReport(msg)
	default:
		s.handleUnknown(msg)
	}
}

func (s *BotService) handleCallback(cb *tgbotapi.CallbackQuery) {
	slog.Info("Callback data", "cb", cb.Data)

	callback := tgbotapi.NewCallback(cb.ID, "")
	s.bot.Request(callback)
}

func (s *BotService) handleStart(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := "Привет! Пока я не знаю, в какой ты десятке 😅"
	message := tgbotapi.NewMessage(chatID, text)
	s.bot.Send(message)
}

func (s *BotService) handleReport(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	user := msg.From
	if user == nil {
		slog.Warn("Got report without user", "chat_id", chatID)
		return
	}

	username := user.UserName
	if username == "" {
		slog.Warn("Report without username", "chat_id", chatID)
		return
	}

	username = "incetro"

	text := msg.Text
	reportParts := strings.SplitN(text, "#отчет", 2)
	if len(reportParts) < 2 {
		s.bot.Send(tgbotapi.NewMessage(chatID, "Неверный формат отчета"))
		return
	}

	reportText := strings.TrimSpace(reportParts[1])
	reportTime := time.Unix(int64(msg.Date), 0).Format("02/01/2006")
	uuidStr := uuid.New().String()

	slog.Info("New report", "user", username, "uuid", uuidStr, "text", reportText, "time", reportTime)

	// 1. Получаем Notion ID автора
	authorNotionID, err := s.repo.GetMemberNotionId(username)
	if err != nil {
		slog.Error("Failed to get author notion ID", "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при определении автора."))
		return
	}
	//2.Создание документа
	if err := s.createDocument(uuidStr, reportText); err != nil {
		slog.Error("Failed to create document", "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при создании документа"))
		return
	}

	//3.Синхронизация документов
	if err := helpers.TriggerSyncDocuments(); err != nil {
		slog.Error("Failed to sync documents", "err", err)
	}

	//4.Получаем Notion ID документа
	docNotionID, err := s.repo.GetDocumentNotionId(uuidStr)
	if err != nil {
		slog.Error("Failed to get document notion ID", "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при поиске документа."))
		return
	}

	// 5. Создание отчёта
	if err := s.createReport(docNotionID, authorNotionID, reportTime); err != nil {
		slog.Error("Failed to create report", "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при создании отчёта."))
		return
	}

	// 6. Синхронизация отчётов
	if err := helpers.TriggerSyncReports(); err != nil {
		slog.Error("Failed to sync reports", "err", err)
	}

	s.bot.Send(tgbotapi.NewMessage(chatID, "Отчёт успешно принят ✅"))
}

func (s *BotService) handleUnknown(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	message := tgbotapi.NewMessage(chatID, "Команда не распознана. Используйте /start.")
	s.bot.Send(message)
}

func (s *BotService) createDocument(id, text string) error {
	page := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(s.notionConfig.DocumentsDBID),
		},
		Properties: notionapi.Properties{
			"ID": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{Content: id},
					},
				},
			},
			"Текст": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{Content: text},
					},
				},
			},
		},
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
