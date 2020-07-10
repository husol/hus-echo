package user

import (
	"context"
	"crm-service/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type User struct {
	getDB func(ctx context.Context) *gorm.DB
}

func New(getDB func(ctx context.Context) *gorm.DB) IUser {
	return &User{getDB}
}

func (p User) GetTableName(ctx context.Context) string {
	return p.getDB(ctx).NewScope(model.User{}).GetModelStruct().TableName(p.getDB(ctx))
}

func (p User) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var obj model.User
	err := p.getDB(ctx).
		Where("email = ?", email).
		First(&obj).Error

	return &obj, errors.Wrap(err, "Fail to GetUserByEmail")
}

func (p User) GetAll(ctx context.Context) ([]model.User, error) {
	var obj []model.User
	err := p.getDB(ctx).Find(&obj).Error

	return obj, errors.Wrap(err, "Fail to GetAll")
}
