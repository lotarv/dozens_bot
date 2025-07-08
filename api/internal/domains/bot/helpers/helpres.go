package helpers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
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
