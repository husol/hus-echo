package user

import (
	"context"
	"crm-service/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type VtigerUser struct {
	getDB func(ctx context.Context) *gorm.DB
}

func New(getDB func(ctx context.Context) *gorm.DB) IVtigerUser {
	return &VtigerUser{getDB}
}

func (p VtigerUser) GetTableName(ctx context.Context) string {
	return p.getDB(ctx).NewScope(model.VtigerUser{}).GetModelStruct().TableName(p.getDB(ctx))
}

func (p VtigerUser) GetUserByEmail(ctx context.Context, email string) (*model.VtigerUser, error) {
	var obj model.VtigerUser
	err := p.getDB(ctx).
		Where("email1 = ?", email).
		First(&obj).Error

	return &obj, errors.Wrap(err, "Fail to GetUserByEmail")
}

func (p VtigerUser) GetAll(ctx context.Context) ([]model.VtigerUser, error) {
	var obj []model.VtigerUser
	err := p.getDB(ctx).Find(&obj).Error

	return obj, errors.Wrap(err, "Fail to GetAll")
}
