package bot

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/jomei/notionapi"
	"github.com/lotarv/dozens_bot/internal/domains/bot/helpers"
	"github.com/lotarv/dozens_bot/internal/domains/bot/repository"
	"github.com/lotarv/dozens_bot/internal/storage"
	"github.com/lotarv/dozens_bot/internal/utils"
)

type BotController struct {
	repo          *repository.BotRepository
	bot           *tgbotapi.BotAPI
	notionClient  *notionapi.Client
	documentsDBID string
	reportsDBID   string
}

func NewBotController(storage *storage.Storage) *BotController {

	repo := repository.New(storage)

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

	documentsDBID := os.Getenv("NOTION_DOCUMENTS_DATABASE_ID")
	if documentsDBID == "" {
		slog.Error("NOTION_DOCUMENTS_DATABASE_ID is not set")
		panic("NOTION_DOCUMENTS_DATABASE_ID is not set")
	}

	reportsDBID := os.Getenv("NOTION_REPORTS_DATABASE_ID")
	if reportsDBID == "" {
		slog.Error("NOTION_REPORTS_DATABASE_ID is not set")
		panic("NOTION_REPORTS_DATABASE_ID is not set")
	}

	http.DefaultClient = utils.GetHTTPClient()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		slog.Error("Failed to initialize bot", "error", err)
		panic("Failed to initialize bot")
	}

	bot.Debug = true
	slog.Info("Authorized on account", "username", bot.Self.UserName)

	notionClient := notionapi.NewClient(notionapi.Token(notionToken))
	return &BotController{
		repo:          repo,
		bot:           bot,
		notionClient:  notionClient,
		documentsDBID: documentsDBID,
		reportsDBID:   reportsDBID,
	}

}

func (c *BotController) Build() {

}

func (c *BotController) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		c.handleIncomingMessage(update.Message)
	}
}

func (c *BotController) handleIncomingMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	user := message.From
	if user == nil {
		slog.Warn("Got a message without user", "chat_id", chatID, "message", message)
		return
	}

	text := message.Text
	lowerText := strings.ToLower(text)

	if strings.Contains(lowerText, "#отчет") {
		reportParts := strings.SplitN(text, "#отчет", 2)
		reportText := strings.TrimSpace(reportParts[1])

		username := user.UserName
		if username == "" {
			slog.Error("message without username", "reportText", reportText)
			return
		}

		reportTime := time.Unix(int64(message.Date), 0).Format("02/01/2006")
		uuidStr := uuid.New().String()

		slog.Info("Creating new report", "username", username, "uuid", uuidStr, "reportTime", reportTime)
		slog.Info("Using Notion API token prefix", "prefix", os.Getenv("NOTION_API_TOKEN")[:6])
		slog.Info("Using Notion database ID", "documentsDBID", c.documentsDBID, "reportsDBID", c.reportsDBID)

		// Создание документа в Notion
		err := c.createDocumentInNotion(context.Background(), uuidStr, reportText)
		if err != nil {
			slog.Error("failed to save document in notion", "uuid", uuidStr, "chat_id", chatID, "error", err)
			return
		} else {
			slog.Info("new document was saved in notion", "uuid", uuidStr)
		}

		// Синхронизация
		if err := helpers.TriggerSyncDocuments(); err != nil {
			slog.Error("failed to synchronize documents after insertion", "error", err)
			return
		}
		slog.Info("documents sync triggered")

		// Получение notion_id документа
		slog.Info("Looking up document notion ID after sync", "uuid", uuidStr)
		document_notion_id, err := c.repo.GetDocumentNotionId(uuidStr)
		if err != nil {
			slog.Error("failed to get document notion id", "uuid", uuidStr, "error", err)
			return
		}
		slog.Info("got document notion id", "notion_id", document_notion_id)

		// Получение ID автора
		author_notion_id, err := c.repo.GetMemberNotionId("incetro")
		if err != nil {
			slog.Error("failed to get member notion id", "error", err)
			return
		}
		slog.Info("got author notion id", "author_notion_id", author_notion_id)

		// Создание отчета
		slog.Info("Creating report in Notion", "documentNotionID", document_notion_id, "authorNotionID", author_notion_id, "date", reportTime)
		if err := c.createReportInNotion(context.Background(), document_notion_id, author_notion_id, reportTime); err != nil {
			slog.Error("failed to save report in notion", "error", err)
			return
		}
		slog.Info("successfully saved new report in notion", "documentID", uuidStr)

		// Синхронизация отчётов
		if err := helpers.TriggerSyncReports(); err != nil {
			slog.Error("failed to synchronize reports after insertion", "error", err)
			return
		}
		slog.Info("successfully synchronized new report in postgres", "documentID", uuidStr)
	}
}

func (c *BotController) createDocumentInNotion(ctx context.Context, id string, text string) error {
	page := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(c.documentsDBID),
		},
		Properties: notionapi.Properties{
			"ID": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: id,
						},
					},
				},
			},
			"Текст": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: text,
						},
					},
				},
			},
		},
	}

	_, err := c.notionClient.Page.Create(ctx, page)
	return err
}

func (c *BotController) createReportInNotion(ctx context.Context, documentNotionID, authorNotionID, date string) error {
	page := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(c.reportsDBID),
		},
		Properties: notionapi.Properties{
			"ID": notionapi.RelationProperty{
				Relation: []notionapi.Relation{
					{ID: notionapi.PageID(documentNotionID)},
				},
			},
			"Автор": notionapi.RelationProperty{
				Relation: []notionapi.Relation{
					{ID: notionapi.PageID(authorNotionID)},
				},
			},
			"Дата создания": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: date,
						},
					},
				},
			},
		},
	}

	_, err := c.notionClient.Page.Create(ctx, page)
	return err
}
