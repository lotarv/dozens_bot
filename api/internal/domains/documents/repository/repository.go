package repository

import (
	"github.com/go-chi/chi/v5"
)

type DocumentsRepository struct {
}

func New(router *chi.Mux) *DocumentsRepository {
	return &DocumentsRepository{}
}

func (r *DocumentsRepository) GetAllDeclarations(username string) {

}

func (r *DocumentsRepository) GetDeclarationByID(declaration_id int) {

}

func (r *DocumentsRepository) GetAllReports(username string) {

}

func (r *DocumentsRepository) GetReportByID(report_id int) {

}

func (r *DocumentsRepository) NewDeclaration() {

}

func (r *DocumentsRepository) NewDocument() {

}
