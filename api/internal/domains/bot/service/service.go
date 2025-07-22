package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/jomei/notionapi"
	"github.com/lotarv/dozens_bot/internal/domains/bot/helpers"
	bot_types "github.com/lotarv/dozens_bot/internal/domains/bot/types/bot"
	user_types "github.com/lotarv/dozens_bot/internal/domains/user/types"
	"github.com/lotarv/dozens_bot/internal/utils/crypto"
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
	CreateDozen(dozen bot_types.Dozen) error
	GetUserDozen(userID int64) (bot_types.Dozen, error)

	GetMembers() ([]bot_types.MemberDB, error)

	SaveDocument(id, encryptedText string) error
	SetEncryptedText(uuidStr string, encryptedText string) error

	ChangeBankBalance(ctx context.Context, piggyBankID int, amount int, reason string, username string) error
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
	chatID := msg.Chat.ID
	userID := msg.From.ID
	text := strings.TrimSpace(msg.Text)

	//0. Проверка на активную сессию транзакции

	if session, ok := transactionSessions[userID]; ok && session.Step != Idle {
		s.handleTransactionStep(msg, session)
		return
	}

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
	case "create_enter_name":
		s.handleDozenNameInput(userID, text)
		return
	}

	//2.Общие команды
	if helpers.IsLikelyReport(text) {
		s.handleReport(msg)
		return
	}
	if msg.Chat.IsPrivate() {
		switch {
		case text == "/start":
			s.handleStart(msg)
		default:
			s.handleUnknown(msg)
		}
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
		s.createDozen(cb.From, userID)
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
	case "create_change_name":
		s.createDozen(cb.From, userID)
	case "create_name_approve":
		s.handleDozenNameApprove(userID)
	case "start_deposit":
		transactionSessions[userID] = &TransactionSession{
			Step:      AwaitingMember,
			IsDeposit: true,
		}
		s.askForTransactionMember(userID)
	case "start_withdraw":
		transactionSessions[userID] = &TransactionSession{
			Step:      AwaitingMember,
			IsDeposit: false,
		}
		s.askForTransactionMember(userID)

	default:
		if strings.HasPrefix(cb.Data, "select_member_") {
			username := strings.TrimPrefix(cb.Data, "select_member_")
			if session, ok := transactionSessions[userID]; ok && session.Step == AwaitingMember {
				session.MemberUsername = username
				session.Step = AwaitingAmount

				s.bot.Send(tgbotapi.NewMessage(userID, "Введите сумму:"))
				return
			}
		}
		s.bot.Send(tgbotapi.NewMessage(chatID, "Неизвестное действие"))
	}
}

func (s *BotService) handleStart(msg *tgbotapi.Message) {
	resetSession(msg.From.ID)
	session := getSession(msg.From.ID)
	slog.Info("New session: ", "session", session)
	session.Dozen = bot_types.Dozen{}

	dozen, err := s.repo.GetUserDozen(msg.From.ID)
	if err == nil {
		session.Dozen = dozen
		s.handleStartRegistered(msg)
		return
	}

	session.User = user_types.User{}
	session.User.ID = msg.From.ID
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

func (s *BotService) handleStartRegistered(msg *tgbotapi.Message) {
	session := getSession(msg.From.ID)

	text := fmt.Sprintf("Добро пожаловать, %s! Вы уже являетесь членом десятки", session.User.FullName)
	s.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

func (s *BotService) handleUnknown(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	message := tgbotapi.NewMessage(chatID, "Команда не распознана. Используйте /start.")
	s.bot.Send(message)
}

func (s *BotService) createDozen(usr *tgbotapi.User, userID int64) {
	session := getSession(userID)

	if session.User.FullName != "" {
		s.bot.Send(tgbotapi.NewMessage(userID, "Введите название вашей десятки:"))
		s.repo.SetUserState(userID, "create_enter_name")
		return
	}
	session.IsCreatingDozen = true
	s.bot.Send(tgbotapi.NewMessage(userID, "Чтобы создать десятку, необходимо внести данные о капитане"))

	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Все верно", "join_enter_sphere"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Изменить имя", "join_change_name"),
		),
	)

	text := fmt.Sprintf("Ваше имя: %s", usr.FirstName+" "+usr.LastName)
	msgOut := tgbotapi.NewMessage(userID, text)
	msgOut.ReplyMarkup = btns

	s.bot.Send(msgOut)
}

func (s *BotService) handleDozenNameInput(userID int64, dozen_name string) {
	session := getSession(userID)

	session.Dozen.Name = dozen_name
	text := fmt.Sprintf("Подтвердите название: %s", dozen_name)
	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Все верно", "create_name_approve"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Внести заново", "create_change_name"),
		),
	)
	msgOut := tgbotapi.NewMessage(userID, text)
	msgOut.ReplyMarkup = btns

	s.bot.Send(msgOut)
}

func (s *BotService) handleDozenNameApprove(userID int64) {
	session := getSession(userID)

	randomCode := helpers.GenerateRandomDozenCode()

	session.Dozen.Code = randomCode
	session.Dozen.Captain = userID

	if err := s.repo.CreateDozen(session.Dozen); err != nil {
		slog.Error("failed to create dozen", "dozen", session.Dozen, "error", err)
		s.bot.Send(tgbotapi.NewMessage(userID, "Не удалось создать десятку. Попробуйте позже!"))
		return
	}

	newDozen, _ := s.repo.GetDozenByCode(session.Dozen.Code)

	if err := s.repo.AddUserToDozen(session.Dozen.Captain, newDozen.ID); err != nil {
		slog.Error("failed to join dozen", "captainID", session.Dozen.Captain, "error", err)
		s.bot.Send(tgbotapi.NewMessage(userID, "Не удалось добавить Вас в десятку. Попробуйте зайти через кнопку присоедениться к десятке"))
	}

	text := fmt.Sprintf("Ваша десятка \"%s\" успешно создана!\nКод десятки: <code>%s</code>", session.Dozen.Name, randomCode)
	s.repo.DeleteUserState(userID)
	resetSession(userID)
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = "HTML"
	s.bot.Send(msg)

}

func (s *BotService) joinDozen(userID int64) {
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите код десятки, к которой хотите присоединиться:"))
	err := s.repo.SetUserState(userID, "join_enter_code")
	if err != nil {
		slog.Error("Failed to set user state", "chat_id", userID, "err", err)
	}
}

func (s *BotService) handleEnterSphere(usr *tgbotapi.User, userID int64) {
	session := getSession(userID)

	if session.User.FullName == "" {
		session.User.FullName = usr.FirstName + " " + usr.LastName
	}
	s.repo.SetUserState(userID, "join_enter_sphere")
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите сферу: "))
}

func (s *BotService) handleChangeName(userID int64) {
	s.repo.SetUserState(userID, "join_changed_name")
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите ваше имя: "))
}

func (s *BotService) handleChangedNameInput(userID int64, text string) {
	session := getSession(userID)

	s.bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Имя изменено на %s", text)))
	session.User.FullName = text
	s.repo.SetUserState(userID, "join_enter_sphere")
	s.handleEnterSphere(nil, userID)
}

func (s *BotService) handleSphereInput(userID int64, text string) {
	session := getSession(userID)

	slog.Info("got sphere", "sphere", text)
	s.repo.SetUserState(userID, "join_enter_income")
	session.User.Niche = text
	s.bot.Send(tgbotapi.NewMessage(userID, "Введите годовой оборот(млн. руб):"))
}

func (s *BotService) handleIncomeInput(userID int64, text string) {
	session := getSession(userID)

	text = strings.ReplaceAll(text, ",", ".") // Поддержка запятых как десятичного разделителя
	income, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	if err != nil || income <= 0 {
		s.bot.Send(tgbotapi.NewMessage(userID, "Пожалуйста, введите корректное число (например: 12.5 или 7)"))
		return
	}

	slog.Info("got income", "income", income)

	session.User.AnnualIncome = income

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
	Годовой оборот: %.1f млн.руб`, session.User.FullName, session.User.Niche, session.User.AnnualIncome))
	msgOut.ReplyMarkup = btns
	s.bot.Send(msgOut)
	slog.Info("current user state", "user", session.User)
}

func (s *BotService) handleJoinCodeInput(msg *tgbotapi.Message, userID int64, chatID int64, code string) {
	session := getSession(userID)

	code = strings.ToLower(strings.TrimSpace(code))
	dozen, err := s.repo.GetDozenByCode(code)
	if err != nil {
		slog.Warn("Invalid dozen code", "code", code, "err", err)
		s.bot.Send(tgbotapi.NewMessage(chatID, "Код десятки не найден. Проверьте правильность и попробуйте снова"))
		return
	}
	slog.Info("got dozen: ", "dozen", dozen)
	session.Dozen = dozen

	//Если зарегистрированный пользователь хочет присоединиться к десятке,сразу добавляем его к десятке
	user, err := s.repo.GetUserByID(context.Background(), msg.From.ID)
	if err == nil && user.ID != 0 {
		session.User = *user
		if err := s.repo.AddUserToDozen(userID, dozen.ID); err != nil {
			slog.Error("failed to add existing user to dozen", "user_id", userID, "err", err)
			s.bot.Send(tgbotapi.NewMessage(chatID, "Не удалось присоединиться к десятке. Попробуйте позже."))
			return
		}

		if err := s.repo.DeleteUserState(userID); err != nil {
			slog.Error("failed to clear user state", "user_id", userID, "err", err)
		}

		s.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Вы успешно присоединились к десятке \"%s\"", dozen.Name)))
		resetSession(userID)
		return
	}

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
	session.User.ID = userID

	if err := s.repo.SetUserState(userID, "join_enter_name"); err != nil {
		slog.Error("failed to set user state", "user_id", userID, "err", err)
	}

}

func (s *BotService) handleJoinSuccess(userID int64) {
	session := getSession(userID)

	if err := s.repo.CreateUser(context.Background(), &session.User); err != nil {
		slog.Error("failed to create user", "error", err, "user", session.User)
		s.bot.Send(tgbotapi.NewMessage(userID, "Не удалось создать пользователя. Попробуйте позже."))
		return
	}

	if session.IsCreatingDozen { //Выходим из ветки с присоединением, если пользователь создает десятку без регистрации
		text := fmt.Sprintf(`Вы успешно зарегистрировались,  "%s". Переходим к созданию десятки`, session.User.FullName)
		s.bot.Send(tgbotapi.NewMessage(userID, text))
		s.createDozen(nil, userID)
		return
	}

	if err := s.repo.AddUserToDozen(userID, session.Dozen.ID); err != nil {
		slog.Error("failed to join dozen user", "error", err, "user", session.User, "dozen", session.Dozen)
		s.bot.Send(tgbotapi.NewMessage(userID, "Не удалось присоединиться к десятке, попробуйте позже."))
		return
	}

	if err := s.repo.DeleteUserState(userID); err != nil {
		slog.Error("failed to remove current user state", "userID", userID)
	}

	text := fmt.Sprintf(`Вы успешно присоединились к десятке "%s"`, session.Dozen.Name)
	resetSession(userID)
	s.bot.Send(tgbotapi.NewMessage(userID, text))
}

func (s *BotService) handleJoinReset(userID int64) {
	session := getSession(userID)

	s.repo.SetUserState(userID, "join_enter_name")

	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Все верно", "join_enter_sphere"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Изменить имя", "join_change_name"),
		),
	)

	text := fmt.Sprintf("Начнем сначала. Ваше имя: %s", session.User.FullName)
	msgOut := tgbotapi.NewMessage(userID, text)
	msgOut.ReplyMarkup = btns

	s.bot.Send(msgOut)
	session.User.ID = userID

	if err := s.repo.SetUserState(userID, "join_enter_name"); err != nil {
		slog.Error("failed to set user state", "user_id", userID, "err", err)
	}

}
func (s *BotService) createDocument(id string) error {
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
		username = msg.From.UserName
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

	// Всё, кроме первой строки
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
}
