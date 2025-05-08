package service

import ()

type repository interface {
	GetAllDeclarations(username string)
	GetDeclarationByID(declaration_id int)
	GetAllReports(username string)
	GetReportByID(report_id int)
	NewDeclaration()
	NewDocument()
}

type DocumentsService struct {
	repo repository
}

func New(repo repository) *DocumentsService {
	return &DocumentsService{
		repo: repo,
	}
}

func (s *DocumentsService) GetAllDeclarations(username string) {

}

func (s *DocumentsService) GetDeclarationByID(declaration_id int) {

}

func (s *DocumentsService) GetAllReports(user_url string) {

}

func (s *DocumentsService) GetReportByID(report_id int) {

}

func (s *DocumentsService) NewDeclaration() {

}

func NewDocument() {

}
