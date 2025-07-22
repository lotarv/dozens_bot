package bot

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jomei/notionapi"
	"github.com/lotarv/dozens_bot/internal/domains/bot/repository"
	"github.com/lotarv/dozens_bot/internal/domains/bot/service"
	"github.com/lotarv/dozens_bot/internal/domains/bot/transport"
	bot_types "github.com/lotarv/dozens_bot/internal/domains/bot/types/bot"
	bank_repo "github.com/lotarv/dozens_bot/internal/domains/piggy_bank/repository"
	user_repo "github.com/lotarv/dozens_bot/internal/domains/user/repository"
	"github.com/lotarv/dozens_bot/internal/storage"
	"github.com/lotarv/dozens_bot/internal/utils"
)

type BotController struct {
	repo      *repository.BotRepository
	service   *service.BotService
	transport *transport.BotTransport
}

func NewBotController(storage *storage.Storage, userRepo *user_repo.UsersRepository, bankRepo *bank_repo.PiggyBankRepository, router *chi.Mux) *BotController {

	telegram_bot := createTelegramBot()

	notion_client := createNotionClient()

	notion_config := createNotionConfig()

	repo := repository.New(storage, userRepo, bankRepo)

	service := service.New(repo, telegram_bot, notion_client, notion_config)

	transport := transport.New(router, service)

	return &BotController{
		repo:      repo,
		service:   service,
		transport: transport,
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

	bot.Debug = false
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
	c.transport.RegisterRoutes()

}

func (c *BotController) Run() {
	go c.service.Run()

}
