// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"shop/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (

	// CategoryRepo -.
	CategoryRepoI interface {
		Create(ctx context.Context, req *entity.CategoryCreate) error
		GetById(ctx context.Context, req *entity.ById) (*entity.CategoryRes, error)
		GetAll(ctx context.Context, req *entity.Filter) (*entity.CategoryGetAllRes, error)
		Update(ctx context.Context, req *entity.CategoryUpdate) error
		Delete(ctx context.Context, req *entity.ById) error
	}

	// ProductRepo -.
	ProductRepoI interface {
		Create(ctx context.Context, req *entity.ProductCreate) error
		GetById(ctx context.Context, req *entity.ById) (*entity.ProductRes, error)
		GetAll(ctx context.Context, req *entity.ProductGetAllReq) (*entity.ProductGetAllRes, error)
		Update(ctx context.Context, req *entity.ProductUpdate) error
		Delete(ctx context.Context, req *entity.ById) error
	}

	// UserRepo -.
	UserRepoI interface {
		Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error)
		Create(ctx context.Context, req *entity.CreateUser) error
		GetById(ctx context.Context, req *entity.ById) (*entity.UserInfo, error)
		GetAll(ctx context.Context, req *entity.Filter, name string) (*entity.UserList, error)
		Update(ctx context.Context, req *entity.UpdateUser) error
		Delete(ctx context.Context, req *entity.ById) error
	}

	// BucketRepo -.
	BucketRepoI interface {
		Create(ctx context.Context, req *entity.BucketItemCreate) error
		GetBucket(ctx context.Context, user_id string) (*entity.BucketRes, error)
		Update(ctx context.Context, req *entity.BucketItemUpdate) error
		Delete(ctx context.Context, id string) error
	}

	// OrderRepo -.
	OrderRepoI interface {
		Create(ctx context.Context, req *entity.OrderCreate) error
		GetOrders(ctx context.Context, status string, user_id string) (*[]entity.OrderRes, error)
		Update(ctx context.Context, status, id string) error
	}

	// DebtLogsRepo -.
	DebtLogsRepoI interface {
		Create(ctx context.Context, req *entity.DebtLogCreate) error
		GetDebtLogs(ctx context.Context, user_id string, status string) (*entity.DebtLogGetAllRes, error)
		Update(ctx context.Context, req entity.DebtLogUpdate) error
	}
)
