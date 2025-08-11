package usecase

import (
	"shop/config"
	"shop/internal/usecase/repo"
	"shop/pkg/logger"
	"shop/pkg/postgres"
)

type UseCase struct {
	CategoryRepo CategoryRepoI
	ProductRepo  ProductRepoI
	UserRepo     UserRepoI
	BucketRepo   BucketRepoI
	OrderRepo    OrderRepoI
	DebtLogsRepo DebtLogsRepoI
}

func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		CategoryRepo: repo.NewCategoryRepo(pg, config, logger),
		ProductRepo:  repo.NewProductRepo(pg, config, logger),
		UserRepo:     repo.NewUserRepo(pg, config, logger),
		BucketRepo:   repo.NewBucketRepo(pg, config, logger),
		OrderRepo:    repo.NewOrderRepo(pg, config, logger),
		DebtLogsRepo: repo.NewDebtLogsRepo(pg, config, logger),
	}
}
