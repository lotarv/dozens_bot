package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	document_types "github.com/lotarv/dozens_bot/internal/domains/documents/types"
	member_types "github.com/lotarv/dozens_bot/internal/domains/members/types"
	"github.com/lotarv/dozens_bot/internal/utils"
)

type repository interface {
	SyncMembersWithNotion([]member_types.Member) error
	SyncDeclarationsWithNotion([]document_types.Declaration) error
	SyncReportsWithNotion([]document_types.Report) error
	SyncDocumentsWithNotion([]document_types.Document) error
}

type NotionSyncService struct {
	repo repository
}

func New(repo repository) *NotionSyncService {
	return &NotionSyncService{
		repo: repo,
	}
}

func (s *NotionSyncService) SyncMembersWithNotion() error {

	var members []member_types.Member
	members, err := fetchMembersFromNotion()
	if err != nil {
		return err
	}

	err = s.repo.SyncMembersWithNotion(members)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotionSyncService) SyncDeclarationsWithNotion() error {
	declarations, err := fetchDeclarationsFromNotion()
	if err != nil {
		return err
	}

	err = s.repo.SyncDeclarationsWithNotion(declarations)
	if err != nil {
		return err
	}
	return nil

}

func (s *NotionSyncService) SyncReportsWithNotion() error {
	reports, err := fetchReportsFromNotion()
	if err != nil {
		return err
	}

	err = s.repo.SyncReportsWithNotion(reports)
	if err != nil {
		return err
	}

	return nil
}

func (s *NotionSyncService) SyncDocumentsWithNotion() error {
	documents, err := fetchDocumentsFromNotion()
	if err != nil {
		return err
	}

	err = s.repo.SyncDocumentsWithNotion(documents)
	if err != nil {
		return err
	}
	return nil
}

func fetchDeclarationsFromNotion() ([]document_types.Declaration, error) {
	client := utils.GetHTTPClient()
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", os.Getenv("NOTION_DECLARATIONS_DATABASE_ID"))
	// Простой запрос без фильтров
	reqBody := map[string]interface{}{}
	body, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal request", "error", err)
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_TOKEN"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	resp, err := client.Do(httpReq)
	if err != nil {
		slog.Error("failed to execute request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected status code", "code", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Парсим JSON деклараций
	var declResult struct {
		Results []struct {
			ID         string `json:"id"`
			Properties struct {
				CreationDate struct {
					Date struct {
						Start string `json:"start"`
					} `json:"date"`
				} `json:"Дата создания"`
				EndDate struct {
					Date struct {
						Start string `json:"start"`
					} `json:"date"`
				} `json:"Дата окончания"`
				Author struct {
					Relation []struct {
						ID string `json:"id"`
					} `json:"relation"`
				} `json:"Автор"`
				ID struct {
					Relation []struct {
						ID string `json:"id"`
					} `json:"relation"`
				} `json:"ID"`
				Status struct {
					Status struct {
						Name string `json:"name"`
					} `json:"status"`
				} `json:"Статус"`
			} `json:"properties"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&declResult); err != nil {
		slog.Error("failed to decode response", "error", err)
		return nil, err
	}

	// Извлекаем декларации
	declarations := make([]document_types.Declaration, 0, len(declResult.Results))
	for i, page := range declResult.Results {
		declaration := document_types.Declaration{
			CreationDate: page.Properties.CreationDate.Date.Start,
			EndDate:      page.Properties.EndDate.Date.Start,
		}

		// Извлекаем ID (UUID как строка)
		declaration.ID = page.ID

		if len(page.Properties.ID.Relation) > 0 {
			declaration.DocumentID = page.Properties.ID.Relation[0].ID
		} else {
			slog.Warn("ID relation not found", "index", i)
			declaration.DocumentID = ""
		}

		// Извлекаем AuthorNotionID
		if len(page.Properties.Author.Relation) > 0 {
			declaration.AuthorNotionID = page.Properties.Author.Relation[0].ID
		} else {
			slog.Warn("Author relation not found", "index", i)
			declaration.AuthorNotionID = ""
		}

		//Извлекаем статус
		if page.Properties.Status.Status.Name != "" {
			declaration.Status = page.Properties.Status.Status.Name
		}

		declarations = append(declarations, declaration)
	}

	// Логируем успешное получение
	slog.Info("Successfully fetched declarations from Notion", "length", len(declarations))
	slog.Info("what we got: ", "declarations", declarations)
	return declarations, nil
}

func fetchMembersFromNotion() ([]member_types.Member, error) {
	client := utils.GetHTTPClient()
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", os.Getenv("NOTION_USERS_DATABASE_ID"))
	// Простой запрос без фильтров
	reqBody := map[string]interface{}{}
	body, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal request", "error", err)
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_TOKEN"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	resp, err := client.Do(httpReq)
	if err != nil {
		slog.Error("failed to execute request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected status code", "code", resp.StatusCode)
		return nil, err
	}

	// Парсим сырой JSON
	var result struct {
		Results []struct {
			ID         string `json:"id"`
			Properties struct {
				FIO struct {
					Title []struct {
						PlainText string `json:"plain_text"`
					} `json:"title"`
				} `json:"ФИО"`
				AvatarUrl struct {
					Files []struct {
						File struct {
							URL string `json:"url"`
						} `json:"file"`
					} `json:"files"`
				} `json:"Фото"`
				Niche struct {
					RichText []struct {
						PlainText string `json:"plain_text"`
					} `json:"rich_text"`
				} `json:"Ниша бизнеса"`
				AnnualIncome struct {
					Number int64 `json:"number"`
				} `json:"Годовой оборот"`
				Username struct {
					URL string `json:"url"`
				} `json:"Telegram"`
			} `json:"properties"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("failed to decode response", "error", err)
		return nil, err
	}

	// Преобразуем данные в упрощенный формат
	members := make([]member_types.Member, 0, len(result.Results))
	for _, page := range result.Results {
		member := member_types.Member{
			Notion_database_id: page.ID,
			AnnualIncome:       page.Properties.AnnualIncome.Number,
			Username:           page.Properties.Username.URL,
		}

		// Извлекаем ФИО (берем первый элемент title)
		if len(page.Properties.FIO.Title) > 0 {
			member.FIO = page.Properties.FIO.Title[0].PlainText
		}

		// Извлекаем Нишу Бизнеса (берем первый элемент rich_text)
		if len(page.Properties.Niche.RichText) > 0 {
			member.Niche = page.Properties.Niche.RichText[0].PlainText
		}

		// Извлекаем Фото (берем URL первого файла, если есть)
		if len(page.Properties.AvatarUrl.Files) > 0 {
			notionURL := page.Properties.AvatarUrl.Files[0].File.URL

			if page.Properties.Username.URL != "" {
				username := member.Username
				filename := fmt.Sprintf("member_%s.jpg", username)
				localPath := "/static/members/" + filename
				publicPath := "/static/members/" + filename

				err := utils.DownloadImage(notionURL, localPath)
				if err != nil {
					slog.Error("failed to download image", "url", notionURL, "error", err)
				} else {
					member.AvatarUrl = publicPath
				}
			}
		}

		members = append(members, member)
	}

	return members, nil
}

func fetchReportsFromNotion() ([]document_types.Report, error) {
	client := utils.GetHTTPClient()
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", os.Getenv("NOTION_REPORTS_DATABASE_ID"))
	// Простой запрос без фильтров
	reqBody := map[string]interface{}{}
	body, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal request", "error", err)
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_TOKEN"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	resp, err := client.Do(httpReq)
	if err != nil {
		slog.Error("failed to execute request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected status code", "code", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Парсим сырой JSON
	var result struct {
		Results []struct {
			ID         string `json:"id"`
			Properties struct {
				Author struct {
					Relation []struct {
						ID string `json:"id"`
					} `json:"relation"`
				} `json:"Автор"`
				ReportID struct {
					Relation []struct {
						ID string `json:"id"`
					} `json:"relation"`
				} `json:"ID"`
				CreationDate struct {
					Title []struct {
						PlainText string `json:"plain_text"`
					} `json:"title"`
				} `json:"Дата создания"`
			} `json:"properties"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("не удалось декодировать ответ", "ошибка", err)
		return nil, err
	}

	// Преобразуем данные в формат Report
	reports := make([]document_types.Report, 0, len(result.Results))
	for _, page := range result.Results {
		report := document_types.Report{
			ID: page.ID,
		}

		// Получаем relation на документ
		if len(page.Properties.ReportID.Relation) > 0 {
			report.DocumentID = page.Properties.ReportID.Relation[0].ID
		}
		// Извлекаем ID автора (берем первый элемент relation)
		if len(page.Properties.Author.Relation) > 0 {
			report.AuthorNotionID = page.Properties.Author.Relation[0].ID
		}

		// Извлекаем и форматируем дату создания
		if len(page.Properties.CreationDate.Title) > 0 {
			dateStr := page.Properties.CreationDate.Title[0].PlainText
			// Парсим дату из формата DD/MM/YYYY
			parsedDate, err := time.Parse("02/01/2006", dateStr)
			if err != nil {
				slog.Error("не удалось распарсить дату", "report_id", page.ID, "дата", dateStr, "ошибка", err)
				return nil, fmt.Errorf("не удалось распарсить дату %s для отчета %s: %w", dateStr, page.ID, err)
			}
			// Форматируем в YYYY-MM-DD для PostgreSQL
			report.CreationDate = parsedDate.Format("2006-01-02")
		}

		reports = append(reports, report)
	}
	slog.Info("reports", "reports", reports)
	return reports, nil
}

func fetchDocumentsFromNotion() ([]document_types.Document, error) {
	client := utils.GetHTTPClient()
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", os.Getenv("NOTION_DOCUMENTS_DATABASE_ID"))
	reqBody := map[string]interface{}{}
	body, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("не удалось сериализовать запрос", "ошибка", err)
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("не удалось создать запрос", "ошибка", err)
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_TOKEN"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	resp, err := client.Do(httpReq)
	if err != nil {
		slog.Error("не удалось выполнить запрос", "ошибка", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("неожиданный код состояния", "код", resp.StatusCode)
		return nil, fmt.Errorf("неожиданный код состояния: %d", resp.StatusCode)
	}

	// Парсим JSON-ответ
	var result struct {
		Results []struct {
			ID         string `json:"id"`
			Properties struct {
				ID struct {
					Title []struct {
						PlainText string `json:"plain_text"`
					} `json:"title"`
				} `json:"ID"`
				Text struct {
					RichText []struct {
						PlainText string `json:"plain_text"`
					} `json:"rich_text"`
				} `json:"Текст"`
			} `json:"properties"`
			CreatedTime string `json:"created_time"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("не удалось декодировать ответ", "ошибка", err)
		return nil, err
	}

	// Преобразуем данные в []document_types.Document
	documents := make([]document_types.Document, 0, len(result.Results))
	for _, page := range result.Results {
		doc := document_types.Document{
			DocumentNotionID: page.ID,
		}

		// Извлекаем ID (должно быть числом)
		if len(page.Properties.ID.Title) > 0 {
			idStr := page.Properties.ID.Title[0].PlainText
			doc.ID = idStr
		} else {
			slog.Error("поле ID пустое", "document_notion_id", page.ID)
			return nil, fmt.Errorf("поле ID пустое для документа %s", page.ID)
		}

		// Извлекаем текст (объединяем все rich_text.plain_text)
		if len(page.Properties.Text.RichText) > 0 {
			var textParts []string
			for _, richText := range page.Properties.Text.RichText {
				textParts = append(textParts, richText.PlainText)
			}
			doc.Text = strings.Join(textParts, "")
		}

		documents = append(documents, doc)
	}

	return documents, nil
}
