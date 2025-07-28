package handler

import (
	"shop/config"
	"shop/internal/usecase"
	"shop/pkg/logger"
	"shop/pkg/minio"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	MinIO   *minio.MinIO
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase, mn minio.MinIO) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		MinIO:   &mn,
	}
}
