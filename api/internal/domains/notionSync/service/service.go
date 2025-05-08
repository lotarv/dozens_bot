package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	document_types "github.com/lotarv/dozens_bot/internal/domains/documents/types"
	member_types "github.com/lotarv/dozens_bot/internal/domains/members/types"
	"github.com/lotarv/dozens_bot/internal/utils"
)

type repository interface {
	SyncMembersWithNotion([]member_types.Member) error
	SyncDeclarationsWithNotion([]document_types.Declaration) error
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
		if len(page.Properties.ID.Relation) > 0 {
			declaration.ID = page.Properties.ID.Relation[0].ID
		} else {
			slog.Warn("ID relation not found", "index", i)
			declaration.ID = ""
		}

		// Извлекаем AuthorNotionID
		if len(page.Properties.Author.Relation) > 0 {
			declaration.AuthorNotionID = page.Properties.Author.Relation[0].ID
		} else {
			slog.Warn("Author relation not found", "index", i)
			declaration.AuthorNotionID = ""
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
			member.AvatarUrl = page.Properties.AvatarUrl.Files[0].File.URL
		}

		members = append(members, member)
	}

	return members, nil
}
