package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/lotarv/dozens_bot/internal/types"
)

type Credentials struct {
	ID        int64  `json:"id" db:"user_id"`
	Username  string `json:"username" db:"username"`
	IsPremium bool   `json:"is_premium" db:"is_premium"`
	FullName  string `json:"full_name" db:"full_name"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
}

func CheckTelegramAuth(initData string) (Credentials, error) {
	parsedData, _ := url.QueryUnescape(initData)
	chunks := strings.Split(parsedData, "&")
	var dataPairs [][]string
	hash := ""
	user := &struct {
		ID        int64  `json:"id"`
		Username  string `json:"username"`
		IsPremium bool   `json:"is_premium"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		PhotoUrl  string `json:"photo_url"`
	}{}
	// Filter and split the chunks
	for _, chunk := range chunks {
		if strings.HasPrefix(chunk, "user=") {
			parsedData = strings.TrimPrefix(chunk, "user=")
			if err := json.Unmarshal([]byte(parsedData), user); err != nil {
				slog.Error("Failed to unmarshal user data", "error", err)
				return Credentials{}, fmt.Errorf("failed to unmarshal user data: %w", err)
			}
		}
		if strings.HasPrefix(chunk, "hash=") {
			hash = strings.TrimPrefix(chunk, "hash=")
		} else {
			pair := strings.SplitN(chunk, "=", 2)
			dataPairs = append(dataPairs, pair)
		}
	}

	// Sort the data pairs by the key
	sort.Slice(dataPairs, func(i, j int) bool {
		return dataPairs[i][0] < dataPairs[j][0]
	})

	// Join the sorted data pairs into the initData string
	var sortedData []string
	for _, pair := range dataPairs {
		sortedData = append(sortedData, fmt.Sprintf("%s=%s", pair[0], pair[1]))
	}
	initData = strings.Join(sortedData, "\n")
	// Create the secret key using HMAC and the given token
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(os.Getenv("BOT_TOKEN")))
	slog.Info(os.Getenv("BOT_TOKEN"))
	secretKey := h.Sum(nil)

	// Create the data check using the secret key and initData
	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(initData))
	dataCheck := h.Sum(nil)
	// Сравниваем подписи
	if fmt.Sprintf("%x", dataCheck) != hash {
		return Credentials{}, fmt.Errorf("invalid hash: signatures do not match")
	}
	return Credentials{
		ID:        user.ID,
		Username:  user.Username,
		IsPremium: user.IsPremium,
		FullName:  fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		AvatarUrl: user.PhotoUrl,
	}, nil
}

func NewAuthMiddleWare() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initData := r.Header.Get("X-Telegram-Init-Data")
			if initData == "" {
				http.Error(w, "X-Telegram-Init-Data header is required", http.StatusUnauthorized)
				slog.Error("Missing X-Telegram-Init-Data header")
				return
			}
			slog.Info("Trying to authorize: init_data: %v", initData)
			creds, err := CheckTelegramAuth(initData)
			slog.Info("Creds: %v", creds)
			slog.Info("error: %v", err)
			if err != nil {
				slog.Error("error happened due to ")
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				slog.Error("Unauthorized", "error", err)
				return
			}

			ctx := context.WithValue(r.Context(), types.ContextKeyUserID, creds.ID)
			ctx = context.WithValue(ctx, types.ContextKeyCredentials, creds)
			slog.Info("Successful authorize")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
