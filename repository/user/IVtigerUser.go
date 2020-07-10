package user

import (
	"context"
	"crm-service/model"
)

type IVtigerUser interface {
	GetUserByEmail(ctx context.Context, email string) (*model.VtigerUser, error)
	GetAll(ctx context.Context) ([]model.VtigerUser, error)
}
