package repo

import (
	"context"
	"encoding/json"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) Store(ctx context.Context, user entity.User) error {
	return u.Db.Create(&user).Error
}

func (u *UserRepo) GetByUserID(ctx context.Context, userID int) (entity.User, error) {
	var user entity.User
	err := u.Db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	return user, err
}

func (u *UserRepo) GetUserSettingByUserID(ctx context.Context, userID int) (entity.UserSetting, error) {
	var setting entity.UserSetting
	err := u.Db.WithContext(ctx).Where("user_id = ?", userID).First(&setting).Error
	return setting, err
}

func (u *UserRepo) UpdateUserSetting(ctx context.Context, userID int, setting entity.UserSetting) error {
	var user entity.User
	err := u.Db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}
	data, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	user.Ext = string(data)

	if err := u.Db.WithContext(ctx).Where("id = ?", userID).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetByUserName(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := u.Db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return user, err
}
