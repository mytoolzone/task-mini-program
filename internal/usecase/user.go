package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/entity"
)

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(r UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (u UserUseCase) Login(ctx context.Context, username, password string) (entity.User, error) {
	user, err := u.repo.GetByUserName(ctx, username)
	if err != nil {
		return entity.User{}, app_code.New(app_code.ErrorUserNotFound, "用户不存在")
	}

	// 将 password md5 2次
	password = md5V(md5V(password))

	if user.Password != password {
		return entity.User{}, app_code.New(app_code.ErrorUserPassword, "密码错误")
	}

	return user, nil
}

func (u UserUseCase) Register(ctx context.Context, user entity.User) (entity.User, error) {
	user, err := u.repo.GetByUserName(ctx, user.Username)
	if err != nil {
		return entity.User{}, app_code.New(app_code.ErrorUserExist, "用户已存在")
	}

	// 将 password md5 2次
	user.Password = md5V(md5V(user.Password))

	err = u.repo.Store(ctx, user)
	return user, err
}

func (u UserUseCase) UpdateSetting(ctx context.Context, userId int, setting entity.UserSetting) error {
	err := u.repo.UpdateUserSetting(ctx, userId, setting)
	if err != nil {
		return err
	}
	return nil
}

func (u UserUseCase) GetSettingByUserID(ctx context.Context, userID int) (entity.UserSetting, error) {
	setting, err := u.repo.GetUserSettingByUserID(ctx, userID)
	if err != nil {
		return entity.UserSetting{}, err
	}
	return setting, nil
}
