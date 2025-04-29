package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jomei/notionapi"
)

type BotController struct {
	bot              *tgbotapi.BotAPI
	notionClient     *notionapi.Client
	groupsDBID       string
	usersDBID        string
	pendingGroups    sync.Map //Хранит chatID -> ожидаемое название
	registeredGroups sync.Map
	registeredUsers  sync.Map
}

func NewBotController() *BotController {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		slog.Error("BOT_TOKEN is not set")
		panic("BOT_TOKEN is not set")
	}

	notionToken := os.Getenv("NOTION_API_TOKEN")
	if notionToken == "" {
		slog.Error("NOTION_API_TOKEN is not set")
		panic("NOTION_API_TOKEN is not set")
	}

	groupsDBID := os.Getenv("NOTION_GROUPS_DATABASE_ID")
	if groupsDBID == "" {
		slog.Error("NOTION_GROUPS_DATABASE_ID is not set")
		panic("NOTION_GROUPS_DATABASE_ID is not set")
	}

	usersDBID := os.Getenv("NOTION_USERS_DATABASE_ID")
	if usersDBID == "" {
		slog.Error("NOTION_USERS_DATABASE_ID is not set")
		panic("NOTION_USERS_DATABASE_ID is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		slog.Error("Failed to initialize bot", "error", err)
		panic("Failed to initialize bot")
	}

	bot.Debug = true
	slog.Info("Authorized on account", "username", bot.Self.UserName)

	notionClient := notionapi.NewClient(notionapi.Token(notionToken))
	return &BotController{
		bot:              bot,
		notionClient:     notionClient,
		groupsDBID:       groupsDBID,
		usersDBID:        usersDBID,
		pendingGroups:    sync.Map{},
		registeredGroups: sync.Map{},
		registeredUsers:  sync.Map{},
	}

}

func (c *BotController) Build() {

}

func (c *BotController) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)
	ctx := context.Background()
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		if update.Message.NewChatMembers != nil {
			for _, member := range update.Message.NewChatMembers {
				if member.ID == c.bot.Self.ID {
					c.handleBotAddedToGroup(update.Message.Chat)
					continue
				}
			}
		}
		//Обработка ответа с названием сообщества
		if title, ok := c.pendingGroups.Load(chatID); ok && title == true {
			c.handleGroupTitleResponse(ctx, update.Message)
		}

		//Обработка сообщений в зарегистрированных группах
		if _, ok := c.registeredGroups.Load(chatID); ok {
			c.handleUserMessage(ctx, update.Message)
		}
	}
}

func (c *BotController) handleBotAddedToGroup(chat *tgbotapi.Chat) {
	msg := tgbotapi.NewMessage(chat.ID, "Пожалуйста, укажите название сообщества \"Десятка\"")
	_, err := c.bot.Send(msg)
	if err != nil {
		slog.Error("failed to send message", "chat_id", chat.ID, "error", err)
		return
	}
	//Отмечаем, что ждем название
	c.pendingGroups.Store(chat.ID, true)
	slog.Info("Requested group title", "chat_id", chat.ID, "title", chat.Title)

	//Собираем пользователей группы
	// c.collectGroupUsers(ctx, chat)
}

func (c *BotController) handleGroupTitleResponse(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	title := strings.TrimSpace(message.Text)
	if title == "" {
		msg := tgbotapi.NewMessage(chatID, "Название не может быть пустым. Пожалуйста, укажите название сообщества \"Десятка\"")
		_, err := c.bot.Send(msg)
		if err != nil {
			slog.Error("Failed to send message", "chat_id", chatID, "error", err)
		}
		return
	}

	err := c.createGroupInNotion(ctx, chatID, title)
	if err != nil {
		slog.Error("failed to create group in Notion", "chat_id", chatID, "error", err)
		msg := tgbotapi.NewMessage(chatID, "Не удалось сохранить сообщество. Попробуйте позже")
		c.bot.Send(msg)
		return
	}
	c.pendingGroups.Delete(chatID)
	c.registeredGroups.Store(chatID, true)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Сообщество \"%s\" успешно создано!", title))
	_, err = c.bot.Send(msg)
	if err != nil {
		slog.Error("Failed to send message", "chat_id", chatID, "error", err)
	}
	slog.Info("Group created", "chat_id", chatID, "title", title)
}

func (c *BotController) handleUserMessage(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	user := message.From
	if user == nil {
		slog.Warn("Message without user", "chat_id", chatID)
		return
	}

	userKey := fmt.Sprintf("%d:%d", chatID, user.ID)
	if _, ok := c.registeredUsers.Load(userKey); ok {
		return //Пользователь уже зарегистрирован
	}

	slog.Info("User messsage",
		"chat_id", chatID,
		"user_id", user.ID,
		"first_name", user.FirstName,
		"last_name", user.LastName,
		"username", user.UserName,
	)

	//Получаем фотографии пользователя
	photos, err := c.bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: user.ID,
		Limit:  1,
	})

	if err != nil {
		slog.Error("failed to get user profile photos", "chat_id", chatID, "user_id", user.ID, "eror", err)
	} else if photos.TotalCount > 0 && len(photos.Photos) > 0 {
		photo := photos.Photos[0][0]
		slog.Info("User profile photo",
			"chat_id", chatID,
			"user_id", user.ID,
			"file_id", photo.FileID,
			"file_size", photo.FileSize,
			"width", photo.Width,
			"height", photo.Height,
		)

	} else {
		slog.Info("No profile photo", "chat_id", chatID, "user_id", user.ID)
	}

	// Отправляем сообщение о регистрации
	name := user.FirstName
	if name == "" {
		name = user.UserName
	}
	if name == "" {
		name = fmt.Sprintf("ID%d", user.ID)
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Пользователь %s зарегистрирован в сообществе", name))
	_, err = c.bot.Send(msg)
	if err != nil {
		slog.Error("Failed to send registration message", "chat_id", chatID, "user_id", user.ID, "error", err)
		return
	}

	// Отмечаем пользователя как зарегистрированного
	c.registeredUsers.Store(userKey, true)

}

func (c *BotController) createGroupInNotion(ctx context.Context, chatID int64, title string) error {
	page := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(c.groupsDBID),
		},
		Properties: notionapi.Properties{
			"Title": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: title}},
				},
			},
			"ChatID": notionapi.NumberProperty{
				Number: float64(chatID),
			},
		},
	}

	_, err := c.notionClient.Page.Create(ctx, page)
	if err != nil {
		return err
	}
	return nil
}
