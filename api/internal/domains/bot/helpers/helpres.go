package helpers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

func TriggerSyncDocuments() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sync-documents", os.Getenv("BASE_URL")), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", os.Getenv("AUTHORIZATION_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sync failed: %s", string(body))
	}

	return nil
}

func TriggerSyncReports() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sync-reports", os.Getenv("BASE_URL")), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", os.Getenv("AUTHORIZATION_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sync failed: %s", string(body))
	}

	return nil
}

var (
	vowels     = []rune{'a', 'e', 'i', 'o', 'u'}
	consonants = []rune{
		'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k',
		'l', 'm', 'n', 'p', 'q', 'r', 's', 't',
		'v', 'w', 'y', 'z',
	}
)

func GenerateRandomDozenCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]rune, 8)
	for i := range 8 {
		if i%2 == 0 {
			code[i] = consonants[r.Intn(len(consonants))]
		} else {
			code[i] = vowels[r.Intn(len(vowels))]
		}
	}
	return string(code)
}

var knownMembers = []string{
	"Анна Берген",
	"Дмитрий Кокорев",
	"Михаил Степанов",
	"Александр Капустянов",
	"Анастасия Яновская",
	"Максим Борцов",
	"Юля Донцова",
	"Полина HAPPYDAYS",
	"Юрий Терентьев",
	"Михаил Старостин",
	"Илья Новоселов",
	"Андрей Грин",
}

// подготовка всех вариантов
func generateNamePatterns() []string {
	var patterns []string
	for _, fullName := range knownMembers {
		parts := strings.Fields(strings.ToLower(strings.ReplaceAll(fullName, "ё", "е")))
		if len(parts) == 2 {
			first, last := parts[0], parts[1]
			patterns = append(patterns,
				first+last,
				last+first,
				last,
			)
		} else if len(parts) == 1 {
			patterns = append(patterns, parts[0])
		}
	}
	return patterns
}

func IsLikelyReport(text string) bool {
	text = strings.ToLower(strings.ReplaceAll(text, "ё", "е"))
	patterns := generateNamePatterns()

	// 1. Явно указанный тег #отчет
	if strings.Contains(text, "#отчет") {
		return true
	}

	// 2. Ищем все слова, начинающиеся с #
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			if strings.HasPrefix(word, "#") {
				tag := strings.TrimPrefix(word, "#")

				if slices.Contains(patterns, tag) {
					return true
				}
			}
		}
	}

	return false
}

func ExtractReportBody(original string) string {
	normalized := strings.ToLower(strings.ReplaceAll(original, "ё", "е"))
	lines := strings.Split(normalized, "\n")

	foundIdx := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			foundIdx = i
			break
		}
	}

	if foundIdx != -1 && foundIdx+1 < len(lines) {
		// Используем оригинал, а не нормализованный текст:
		origLines := strings.Split(original, "\n")
		return strings.Join(origLines[foundIdx+1:], "\n")
	}

	// fallback: если хештег не найден — всё после первой строки
	origLines := strings.Split(original, "\n")
	if len(origLines) > 1 {
		return strings.Join(origLines[1:], "\n")
	}
	return ""
}

func ResolveUsername(msg *tgbotapi.Message) string {
	if msg.ForwardFrom != nil && msg.ForwardFrom.UserName != "" {
		return msg.ForwardFrom.UserName
	}
	if msg.From != nil && msg.From.UserName != "" {
		return msg.From.UserName
	}
	return "unknown_user"
}
