package service

import (
	"log/slog"

	"github.com/lotarv/dozens_bot/internal/domains/documents/types"
)

type repository interface {
	GetRules() (types.Document, error)
	GetRawReports(username string) ([]types.Report, error)
	GetUserAvatarUrl(username string) (string, error)
	GetReportDocuments(reportsRaw []types.Report) (map[string]types.Document, error)
	GetDeclarations(username string) ([]types.DeclarationDocument, error)
}

type DocumentsService struct {
	repo repository
}

func New(repo repository) *DocumentsService {
	return &DocumentsService{
		repo: repo,
	}
}

func (s *DocumentsService) GetRules() (types.Document, error) {
	return s.repo.GetRules()
}

func (s *DocumentsService) GetReports(username string) (*types.ReportsResponse, error) {
	reportsRaw, err := s.repo.GetRawReports(username)
	if err != nil {
		return nil, err
	}

	userAvatarUrl, err := s.repo.GetUserAvatarUrl(username)
	if err != nil {
		return nil, err
	}

	documentsMap, err := s.repo.GetReportDocuments(reportsRaw)
	if err != nil {
		return nil, err
	}
	slog.Info("show map", "documentsMap", documentsMap)

	var reports []types.ReportItem
	for _, r := range reportsRaw {
		doc, ok := documentsMap[r.DocumentID]
		if !ok {
			continue
		}
		reports = append(reports, types.ReportItem{
			CreationDate: r.CreationDate,
			Text:         doc.Text,
		})
	}

	response := &types.ReportsResponse{
		Username:  username,
		AvatarUrl: userAvatarUrl,
		Reports:   reports,
	}

	return response, nil
}

func (s *DocumentsService) GetDeclarations(username string) ([]types.DeclarationDocument, error) {
	return s.repo.GetDeclarations(username)
}
