package repository

import (
	"context"
	"crm-service/repository/user"
	"github.com/jinzhu/gorm"
)

type Repository struct {
	UserRepo user.IVtigerUser
}

func New(getClient func(ctx context.Context) *gorm.DB) *Repository {
	return &Repository{
		UserRepo: user.New(getClient),
	}
}
