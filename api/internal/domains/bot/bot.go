package bot

import (
	"log/slog"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jomei/notionapi"
	"github.com/lotarv/dozens_bot/internal/domains/bot/repository"
	"github.com/lotarv/dozens_bot/internal/domains/bot/service"
	bot_types "github.com/lotarv/dozens_bot/internal/domains/bot/types/bot"
	"github.com/lotarv/dozens_bot/internal/storage"
	"github.com/lotarv/dozens_bot/internal/utils"
)

type BotController struct {
	repo    *repository.BotRepository
	service *service.BotService
}

func NewBotController(storage *storage.Storage) *BotController {

	telegram_bot := createTelegramBot()

	notion_client := createNotionClient()

	notion_config := createNotionConfig()

	repo := repository.New(storage)

	service := service.New(repo, telegram_bot, notion_client, notion_config)

	return &BotController{
		repo:    repo,
		service: service,
	}

}

func createTelegramBot() *tgbotapi.BotAPI {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		slog.Error("BOT_TOKEN is not set")
		panic("BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		slog.Error("Failed to initialize bot", "error", err)
		panic("Failed to initialize bot")
	}

	bot.Debug = true
	slog.Info("Authorized on account", "username", bot.Self.UserName)

	return bot
}

func createNotionClient() *notionapi.Client {
	notionToken := os.Getenv("NOTION_API_TOKEN")
	if notionToken == "" {
		slog.Error("NOTION_API_TOKEN is not set")
		panic("NOTION_API_TOKEN is not set")
	}

	http.DefaultClient = utils.GetHTTPClient()

	notionClient := notionapi.NewClient(notionapi.Token(notionToken))

	return notionClient
}

func createNotionConfig() bot_types.NotionConfig {

	documentsNotionDBID := os.Getenv("NOTION_DOCUMENTS_DATABASE_ID")
	if documentsNotionDBID == "" {
		panic("NOTION_DOCUMENTS_DATABASE_ID env variable is not set")
	}

	reportsNotionDBID := os.Getenv("NOTION_REPORTS_DATABASE_ID")
	if reportsNotionDBID == "" {
		panic("NOTION_REPORTS_DATABASE_ID env variable is not set")
	}

	notionCfg := bot_types.NotionConfig{
		DocumentsDBID: documentsNotionDBID,
		ReportsDBID:   reportsNotionDBID,
	}

	return notionCfg
}

func (c *BotController) Build() {

}

func (c *BotController) Run() {
	go c.service.Run()

}
