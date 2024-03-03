package repo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
	"gorm.io/gorm"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) Store(ctx context.Context, user *entity.User) error {
	if user.Ext == "" {
		user.Ext = "{}"
	}
	return u.Db.Create(user).Error
}

func (u *UserRepo) GetByUserID(ctx context.Context, userID int) (entity.User, error) {
	var user entity.User
	err := u.Db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	return user, err
}

func (u *UserRepo) GetUserSettingByUserID(ctx context.Context, userID int) (entity.UserSetting, error) {
	var (
		user    entity.User
		setting entity.UserSetting
	)
	err := u.Db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return entity.UserSetting{}, err
	}

	if user.Ext == "" {
		return entity.UserSetting{}, nil
	}
	err = json.Unmarshal([]byte(user.Ext), &setting)
	if err != nil {
		return entity.UserSetting{}, err
	}

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
	user.Username = setting.Name
	user.Email = setting.Email

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

func (u *UserRepo) GetByOpenId(ctx context.Context, openID string) (entity.User, bool, error) {
	var user entity.User
	err := u.Db.WithContext(ctx).Where("openid = ?", openID).First(&user).Error
	if err == nil {
		return user, true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, false, nil
	}
	return user, false, err
}

func (u *UserRepo) GetUserRole(ctx context.Context, userID int) (entity.UserRole, error) {
	var userRole entity.UserRole
	if err := u.Db.WithContext(ctx).Where("user_id =?", userID).First(&userRole).Error; err != nil {
		return entity.UserRole{}, err
	}
	return userRole, nil
}

// SetUserRole 设置或者更新用户的角色
func (u *UserRepo) SetUserRole(ctx context.Context, userID int, role string) error {
	var userRole entity.UserRole
	err := u.Db.WithContext(ctx).Where("user_id =?", userID).First(&userRole).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userRole = entity.UserRole{
			UserID: userID,
			Role:   role,
		}
		err = u.Db.WithContext(ctx).Create(&userRole).Error
		return err
	}

	if err != nil {
		return err
	}
	userRole.Role = role
	return u.Db.WithContext(ctx).Model(&userRole).Updates(&userRole).Error
}

// FindUsersByName 根据用户名模糊查询用户列表
func (u *UserRepo) FindUsersByName(ctx context.Context, username string) ([]entity.User, error) {
	var users []entity.User
	err := u.Db.WithContext(ctx).Select("id,username,phone,email,status").Where("username like ?", "%"+username+"%").Limit(100).Find(&users).Error
	return users, err
}
