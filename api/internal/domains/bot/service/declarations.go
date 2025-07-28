package service

import (
	"log/slog"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/lotarv/dozens_bot/internal/domains/bot/helpers"
	document_types "github.com/lotarv/dozens_bot/internal/domains/documents/types"
	"github.com/lotarv/dozens_bot/internal/utils/crypto"
)

func (s *BotService) handleDeclaration(msg *tgbotapi.Message) {

	chatID := msg.Chat.ID
	// Определяем автора
	var username string
	if msg.ForwardFrom != nil && msg.ForwardFrom.UserName != "" {
		username = msg.ForwardFrom.UserName
	} else if msg.From != nil && msg.From.UserName != "" {
		username = "incetro"
	} else {
		slog.Warn("Declaration without username", "chat_id", chatID)
		s.replyTo(msg, "Ошибка при определении автора декларации")
		return
	}

	var declarationDB document_types.DeclarationDB
	// 1. Получаем Notion ID автора
	authorNotionID, err := s.repo.GetMemberNotionId(username)
	if err != nil {
		slog.Error("Failed to get author notion ID", "err", err)
		s.replyTo(msg, "Ошибка при определении автора декларации")
		return
	}

	// Нормализуем текст
	text := msg.Text
	lines := strings.Split(text, "\n")
	if len(lines) < 1 {
		s.replyTo(msg, "Неверный формат декларации: отсутствует содержимое")
		return
	}

	// Достаем тело отчета: все, кроме строки с хештегом
	declarationText := helpers.ExtractReportBody(msg.Text)

	var declarationTime string
	if msg.ForwardDate != 0 {
		declarationTime = time.Unix(int64(msg.ForwardDate), 0).Format("2006-01-02")
	} else {
		declarationTime = time.Unix(int64(msg.Date), 0).Format("2006-01-02")
	}

	//2 Шифруем данные
	//TODO: брать код десятки текущего пользователя
	pepper := os.Getenv("ENCRYPTION_PEPPER")
	dozenCode := os.Getenv("ONLY_DOZEN_CODE")
	passphrase := dozenCode + pepper

	encryptedText, err := crypto.Encrypt(declarationText, passphrase)
	if err != nil {
		slog.Error("failed to encrypt reportText", "err", err)
		s.replyTo(msg, "Не удалось выполнить шифрование отчета. Попробуйте еще раз")
		return
	}

	//3 Генерим uuid
	uuidStr := uuid.New().String()

	declarationDB.AuthorNotionID = authorNotionID
	declarationDB.ID = uuidStr
	declarationDB.Status = "In progress"
	declarationDB.CreationDate = declarationTime
	declarationDB.Text = encryptedText

	if err := s.repo.CreateDeclaration(declarationDB); err != nil {
		slog.Error("failed to create declaration", "err", err)
		s.replyTo(msg, "Не удалось сохранить декларацию. Попробуйте еще раз")
		return
	}

	s.replyTo(msg, "Декларация успешно принята ✅")
}
