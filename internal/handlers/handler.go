package handlers

import (
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/storage"
)

type Handler struct {
	db            *gorm.DB
	storageClient *storage.Storage
}

func NewHandler(db *gorm.DB, storageClient *storage.Storage) *Handler {
	return &Handler{db: db, storageClient: storageClient}
}
