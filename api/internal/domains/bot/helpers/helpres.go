package helpers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
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
		parts := strings.Fields(strings.ToLower(fullName)) // разбиваем по пробелу
		if len(parts) == 2 {
			first, last := parts[0], parts[1]
			patterns = append(patterns,
				first+" "+last,
				last+" "+first,
				last, // только фамилия
			)
		} else if len(parts) == 1 {
			patterns = append(patterns, parts[0])
		}
	}
	return patterns
}

func IsLikelyReport(text string) bool {
	text = strings.ToLower(strings.ReplaceAll(text, "ё", "е"))

	if strings.Contains(text, "#отчет") {
		return true
	}

	if !strings.HasPrefix(text, "#") {
		return false
	}

	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return false
	}

	firstLine := strings.TrimPrefix(strings.TrimSpace(lines[0]), "#")
	firstLine = strings.ToLower(firstLine)

	for _, pattern := range generateNamePatterns() {
		if strings.HasPrefix(firstLine, pattern) {
			return true
		}
	}

	return false
}
