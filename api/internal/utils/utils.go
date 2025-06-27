package utils

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func GetHTTPClient() *http.Client {
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

func DownloadImage(url string, targetPath string) error {
	if _, err := os.Stat(targetPath); err == nil {
		return nil
	}

	// Создаём директорию, если её нет
	if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file:%w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
