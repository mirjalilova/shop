package usecase

import (
	"shop/config"
	"shop/internal/usecase/repo"
	"shop/pkg/logger"
	"shop/pkg/postgres"
)

type UseCase struct {
	CategoryRepo CategoryRepoI
	ShoesRepo    ShoesRepoI
	UserRepo     UserRepoI
}

func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		CategoryRepo: repo.NewCategoryRepo(pg, config, logger),
		ShoesRepo:    repo.NewShoesRepo(pg, config, logger),
		UserRepo:     repo.NewUserRepo(pg, config, logger),
	}
}
