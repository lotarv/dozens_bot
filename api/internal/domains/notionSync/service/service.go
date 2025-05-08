package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	member_types "github.com/lotarv/dozens_bot/internal/domains/members/types"
	"github.com/lotarv/dozens_bot/internal/utils"
)

type repository interface {
	SyncMembersWithNotion([]member_types.Member) error
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
				} `json:"Ниша Бизнеса"`
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
			AnnualIncome: page.Properties.AnnualIncome.Number,
			Username:     page.Properties.Username.URL,
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
