package usecase

import (
	"context"
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

func (u UserUseCase) Login(ctx context.Context, user entity.User) (entity.User, error) {
	user, err := u.repo.GetByUserID(ctx, user.ID)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u UserUseCase) Register(ctx context.Context, user entity.User) (entity.User, error) {
	user, err := u.repo.GetByUserID(ctx, user.ID)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
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
