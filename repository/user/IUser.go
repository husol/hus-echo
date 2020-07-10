package user

import (
	"context"
	"crm-service/model"
)

type IUser interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
}
