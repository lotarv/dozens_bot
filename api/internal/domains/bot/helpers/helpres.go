package helpers

import (
	"fmt"
	"io"
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
