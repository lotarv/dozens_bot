package external

// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.com/PuerkitoBio/goquery"
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"regexp"
// 	"strings"

// 	"github.com/lotarv/dozens_bot/internal/domains/members"
// 	"github.com/lotarv/dozens_bot/internal/utils"
// )

// func ParsePublicGoogleDoc(url string) (string, error) {
// 	// Отправляем HTTP-запрос
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", fmt.Errorf("ошибка загрузки документа: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("не удалось загрузить документ, статус: %s", resp.Status)
// 	}

// 	// Парсим HTML
// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("ошибка парсинга HTML: %v", err)
// 	}

// 	// Находим контейнер с содержимым (обычно <div id="contents">)
// 	content := doc.Find("#contents")
// 	if content.Length() == 0 {
// 		return "", fmt.Errorf("не удалось найти содержимое документа")
// 	}

// 	// Извлекаем текст
// 	var text strings.Builder
// 	content.Each(func(i int, s *goquery.Selection) {
// 		// Получаем текст, заменяя HTML-теги переносами строк
// 		text.WriteString(s.Text())
// 		text.WriteString("\n")
// 	})

// 	// Удаляем лишние пробелы и переносы строк
// 	result := strings.TrimSpace(text.String())
// 	if result == "" {
// 		return "", fmt.Errorf("документ пуст")
// 	}

// 	return result, nil
// }

// func extractPageID(url string) (string, error) {
// 	// Формат URL: https://www.notion.so/incetro/1e027b0504048016b135f7c62fc2c749
// 	// pageId — 32-символьная строка из букв a-f и цифр 0-9
// 	re := regexp.MustCompile(`[0-9a-f]{32}`)
// 	match := re.FindString(url)
// 	if match == "" {
// 		return "", fmt.Errorf("pageId не найден в URL: %s", url)
// 	}

// 	return match, nil
// }

// func GetMemberData(member_url string) (members.Member, error) {
// 	var member members.Member

// 	client := utils.GetHTTPClient()

// 	pageID, err := extractPageID(member_url)
// 	if err != nil {
// 		return member, fmt.Errorf("failed to extract pageId: %v", err)
// 	}

// 	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", pageID)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return member, fmt.Errorf("error occured while member request creation: %v", err)
// 	}
// 	req.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_TOKEN"))
// 	req.Header.Set("Notion-Version", "2022-06-28")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return member, fmt.Errorf("error executing request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return member, fmt.Errorf("API error, status: %s", resp.Status)
// 	}

// 	var page struct {
// 		ID          string                 `json:"id"`
// 		Propertions map[string]interface{} `json:"properties"`
// 	}

// 	err = json.NewDecoder(resp.Body).Decode(&page)
// 	if err != nil {
// 		return member, fmt.Errorf("error while decoding response: %v", err)
// 	}

// 	slog.Info("Got data from Notion: ", "page", page)
// 	return member, nil

// }
