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
	GetDozenByCode(code string) (int, error)
	ResetUserState(userID int64) error
	SetUserState(userID int64, state string) error
	GetUserState(userID int64) (string, error)
}

type BotService struct {
	repo         repository
	bot          *tgbotapi.BotAPI
	notionClient *notionapi.Client
	notionConfig bot_types.NotionConfig
	userStates   map[int64]bot_types.UserState
}

func New(repo repository, bot *tgbotapi.BotAPI, notionClient *notionapi.Client, config bot_types.NotionConfig) *BotService {

	return &BotService{
		repo:         repo,
		bot:          bot,
		notionClient: notionClient,
		notionConfig: config,
		userStates:   make(map[int64]bot_types.UserState),
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
	chatID := msg.Chat.ID
	userID := msg.From.ID
	text := strings.TrimSpace(msg.Text)

	//1.Проверяем текущее состояние пользователя
	state, err := s.repo.GetUserState(userID)
	if err != nil {
		slog.Error("Failed to get user state", "user_id", userID, "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Произошла ошибка. Попробуйте позже."))
		return
	}

	switch state {
	case "join_enter_code":
		s.handleJoinCodeInput(userID, chatID, text)
		return
	}

	//2.Общие команды

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
	chatID := cb.Message.Chat.ID
	userID := cb.From.ID
	_, _ = s.bot.Request(tgbotapi.NewCallback(cb.ID, "")) //Убираем loading в UI

	// Удаляем inline-кнопки
	s.bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		cb.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
	))

	switch cb.Data {
	case "create_dozen":
		s.createDozen(userID)
	case "join_dozen":
		s.joinDozen(userID)
	default:
		s.bot.Send(tgbotapi.NewMessage(chatID, "Неизвестное действие"))
	}
}

func (s *BotService) createDozen(userID int64) {
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите название вашей десятки:"))
}

func (s *BotService) joinDozen(userID int64) {
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите код десятки, к которой хотите присоединиться:"))
	err := s.repo.SetUserState(userID, "join_enter_code")
	if err != nil {
		slog.Error("Failed to set user state", "chat_id", userID, "err", err)
	}
}

func (s *BotService) handleStart(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := "Добро пожаловать! Выберите действие:"
	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔗 Присоединиться к десятке", "join_dozen"),
			tgbotapi.NewInlineKeyboardButtonData("✨ Создать десятку", "create_dozen"),
		),
	)
	msgOut := tgbotapi.NewMessage(chatID, text)
	msgOut.ReplyMarkup = btns

	s.bot.Send(msgOut)
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

func (s *BotService) handleJoinCodeInput(userID int64, chatID int64, code string) {
	code = strings.ToLower(strings.TrimSpace(code))
	dozenID, err := s.repo.GetDozenByCode(code)
	if err != nil || dozenID == 0 {
		slog.Warn("Invalid dozen code", "code", code, "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Код десятки не найден. Проверьте правильность и попробуйте снова"))
		return
	}

	slog.Info("got dozen id", "dozenID", dozenID)

	s.bot.Send(tgbotapi.NewMessage(chatID, "Код принят! Введите ваше имя:"))

	if err := s.repo.SetUserState(userID, "join_enter_name"); err != nil {
		slog.Error("failed to set user state", "user_id", userID, "err", err)
	}

}
