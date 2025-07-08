package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/jomei/notionapi"
	"github.com/lotarv/dozens_bot/internal/domains/bot/helpers"
	bot_types "github.com/lotarv/dozens_bot/internal/domains/bot/types/bot"
	user_types "github.com/lotarv/dozens_bot/internal/domains/user/types"
)

type repository interface {
	GetDocumentsAmount() (int, error)
	GetDocumentNotionId(db_uuid string) (string, error)
	GetMemberNotionId(username string) (string, error)
	GetDozenByCode(code string) (bot_types.Dozen, error)
	ResetUserState(userID int64) error
	SetUserState(userID int64, state string) error
	GetUserState(userID int64) (string, error)
	DeleteUserState(userID int64) error

	CreateUser(ctx context.Context, user *user_types.User) error
	UpdateUser(ctx context.Context, user *user_types.User) error
	GetUserByID(ctx context.Context, userID int64) (*user_types.User, error)
	AddUserToDozen(userID int64, dozenID int) error
}

type BotService struct {
	repo         repository
	bot          *tgbotapi.BotAPI
	notionClient *notionapi.Client
	notionConfig bot_types.NotionConfig
	user         user_types.User
	dozen        bot_types.Dozen
}

func New(repo repository, bot *tgbotapi.BotAPI, notionClient *notionapi.Client, config bot_types.NotionConfig) *BotService {

	return &BotService{
		repo:         repo,
		bot:          bot,
		notionClient: notionClient,
		notionConfig: config,
		user:         user_types.User{},
		dozen:        bot_types.Dozen{},
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
		s.handleJoinCodeInput(msg, userID, chatID, text)
		return
	case "join_changed_name":
		s.handleChangedNameInput(userID, text)
		return
	case "join_enter_sphere":
		s.handleSphereInput(userID, text)
		return
	case "join_enter_income":
		s.handleIncomeInput(userID, text)
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
	case "join_enter_sphere":
		s.handleEnterSphere(cb.From, userID)
	case "join_change_name":
		s.handleChangeName(userID)
	case "join_success":
		s.handleJoinSuccess(userID)
	case "join_reset":
		s.handleJoinReset(userID)
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

	s.dozen = bot_types.Dozen{}
	s.user = user_types.User{}

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

func (s *BotService) handleUnknown(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	message := tgbotapi.NewMessage(chatID, "Команда не распознана. Используйте /start.")
	s.bot.Send(message)
}

func (s *BotService) handleEnterSphere(usr *tgbotapi.User, userID int64) {
	if s.user.FullName == "" {
		s.user.FullName = usr.FirstName + " " + usr.LastName
	}
	s.repo.SetUserState(userID, "join_enter_sphere")
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите сферу: "))
}

func (s *BotService) handleChangeName(userID int64) {
	s.repo.SetUserState(userID, "join_changed_name")
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите ваше имя: "))
}

func (s *BotService) handleChangedNameInput(userID int64, text string) {
	s.bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Имя изменено на %s", text)))
	s.user.FullName = text
	s.repo.SetUserState(userID, "join_enter_sphere")
	s.handleEnterSphere(nil, userID)
}

func (s *BotService) handleSphereInput(userID int64, text string) {
	slog.Info("got sphere", "sphere", text)
	s.repo.SetUserState(userID, "join_enter_income")
	s.user.Niche = text
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите годовой оборот(млн. руб):"))
}

func (s *BotService) handleIncomeInput(userID int64, text string) {
	text = strings.ReplaceAll(text, ",", ".") // Поддержка запятых как десятичного разделителя
	income, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	if err != nil || income <= 0 {
		s.bot.Send(tgbotapi.NewMessage(userID, "Пожалуйста, введите корректное число (например: 12.5 или 7)"))
		return
	}

	slog.Info("got income", "income", income)

	s.user.AnnualIncome = income

	if err := s.repo.SetUserState(userID, "join_approve"); err != nil {
		slog.Error("failed to set state", "user_id", userID, "err", err)
		return
	}

	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Все верно", "join_success"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Внести заново", "join_reset"),
		),
	)

	msgOut := tgbotapi.NewMessage(userID, fmt.Sprintf(`
	Принято! Проверьте правильность данных:
	Имя: %s
	Сфера бизнеса: %s
	Годовой оборот: %.1f млн.руб`, s.user.FullName, s.user.Niche, s.user.AnnualIncome))
	msgOut.ReplyMarkup = btns
	s.bot.Send(msgOut)
	slog.Info("current user state", "user", s.user)
}

func (s *BotService) handleJoinCodeInput(msg *tgbotapi.Message, userID int64, chatID int64, code string) {
	code = strings.ToLower(strings.TrimSpace(code))
	dozen, err := s.repo.GetDozenByCode(code)
	if err != nil {
		slog.Warn("Invalid dozen code", "code", code, "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Код десятки не найден. Проверьте правильность и попробуйте снова"))
		return
	}
	slog.Info("got dozen: ", "dozen", dozen)
	s.dozen = dozen

	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Все верно", "join_enter_sphere"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Изменить имя", "join_change_name"),
		),
	)

	text := fmt.Sprintf("Код принят! Ваше имя: %s", msg.From.FirstName+" "+msg.From.LastName)
	msgOut := tgbotapi.NewMessage(chatID, text)
	msgOut.ReplyMarkup = btns

	s.bot.Send(msgOut)
	s.user.ID = userID

	if err := s.repo.SetUserState(userID, "join_enter_name"); err != nil {
		slog.Error("failed to set user state", "user_id", userID, "err", err)
	}

}

func (s *BotService) handleJoinSuccess(userID int64) {
	if err := s.repo.CreateUser(context.Background(), &s.user); err != nil {
		slog.Error("failed to create user", "error", err, "user", s.user)
		s.bot.Send(tgbotapi.NewMessage(userID, "Не удалось создать пользователя. Попробуйте позже."))
		return
	}

	if err := s.repo.AddUserToDozen(userID, s.dozen.ID); err != nil {
		slog.Error("failed to join dozen user", "error", err, "user", s.user, "dozen", s.dozen)
		s.bot.Send(tgbotapi.NewMessage(userID, "Не удалось присоединиться к десятке, попробуйте позже."))
		return
	}

	if err := s.repo.DeleteUserState(userID); err != nil {
		slog.Error("failed to remove current user state", "userID", userID)
	}

	text := fmt.Sprintf(`Вы успешно присоединились к десятке "%s"`, s.dozen.Name)
	s.bot.Send(tgbotapi.NewMessage(userID, text))
}

func (s *BotService) handleJoinReset(userID int64) {
	s.repo.SetUserState(userID, "join_enter_name")

	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Все верно", "join_enter_sphere"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Изменить имя", "join_change_name"),
		),
	)

	text := fmt.Sprintf("Начнем сначала. Ваше имя: %s", s.user.FullName)
	msgOut := tgbotapi.NewMessage(userID, text)
	msgOut.ReplyMarkup = btns

	s.bot.Send(msgOut)
	s.user.ID = userID

	if err := s.repo.SetUserState(userID, "join_enter_name"); err != nil {
		slog.Error("failed to set user state", "user_id", userID, "err", err)
	}

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
