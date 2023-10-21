// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/mytoolzone/task-mini-program/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}

	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// User -.
	User interface {
		Login(context.Context, entity.User) (entity.User, error)
		Register(context.Context, entity.User) (entity.User, error)
		UpdateSetting(ctx context.Context, userID int, setting entity.UserSetting) error
		GetSettingByUserID(ctx context.Context, userID int) (entity.UserSetting, error)
	}

	// UserRepo -.
	UserRepo interface {
		Store(context.Context, entity.User) error
		GetByUserID(ctx context.Context, userID int) (entity.User, error)
		GetUserSettingByUserID(ctx context.Context, userID int) (entity.UserSetting, error)
		UpdateUserSetting(ctx context.Context, userID int, setting entity.UserSetting) error
	}
)
