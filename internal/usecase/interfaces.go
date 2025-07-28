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

	// ShoesRepo -.
	ShoesRepoI interface {
		Create(ctx context.Context, req *entity.ShoesCreate) error
		GetById(ctx context.Context, req *entity.ById) (*entity.ShoesRes, error)
		GetAll(ctx context.Context, req *entity.Filter) (*entity.ShoesGetAllRes, error)
		Update(ctx context.Context, req *entity.ShoesUpdate) error
		Delete(ctx context.Context, req *entity.ById) error
	}

	// UserRepo -.
	UserRepoI interface {
		Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error)
		Create(ctx context.Context, req *entity.CreateUser) error
		GetById(ctx context.Context, req *entity.ById) (*entity.UserInfo, error)
		GetAll(ctx context.Context, req *entity.Filter) (*entity.UserList, error)
		Update(ctx context.Context, req *entity.UpdateUser) error
		Delete(ctx context.Context, req *entity.ById) error
	}
)
