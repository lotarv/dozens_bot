package utils

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
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
