package members

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/lotarv/dozens_bot/internal/domains/members/types"
)

type MembersController struct {
	router *chi.Mux
}

func NewMembersController(router *chi.Mux) *MembersController {
	return &MembersController{router: router}
}

func (c *MembersController) Build() {
	slog.Info("Building members controller...")
	c.router.Get("/api/members", c.handleGetMembers)
}

func (c *MembersController) Run() {

}

func getHTTPClient() *http.Client {
	transport := &http.Transport{}

	proxy := os.Getenv("PROXY_URL")
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			slog.Error("Failed to parse proxy URL", slog.String("Error", err.Error()))
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &http.Client{
		Transport: transport,
	}
}

func (c *MembersController) handleGetMembers(w http.ResponseWriter, r *http.Request) {
	slog.Info("function is called")
	client := getHTTPClient()
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", os.Getenv("NOTION_USERS_DATABASE_ID"))

	// Простой запрос без фильтров
	reqBody := map[string]interface{}{}
	body, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal request", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	httpReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_TOKEN"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	resp, err := client.Do(httpReq)
	if err != nil {
		slog.Error("failed to execute request", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected status code", "code", resp.StatusCode)
		http.Error(w, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), http.StatusInternalServerError)
		return
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Преобразуем данные в упрощенный формат
	members := make([]types.Member, 0, len(result.Results))
	for _, page := range result.Results {
		member := types.Member{
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

	// Отправляем упрощенный ответ
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	if err := json.NewEncoder(w).Encode(members); err != nil {
		slog.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
