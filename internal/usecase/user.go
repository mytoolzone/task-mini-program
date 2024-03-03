package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/wechat"
	"gorm.io/gorm"
)

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type UserUseCase struct {
	repo  UserRepo
	wxApp wechat.WxApp
}

func NewUserUseCase(r UserRepo, wxApp wechat.WxApp) *UserUseCase {
	return &UserUseCase{
		repo:  r,
		wxApp: wxApp,
	}
}

func (u UserUseCase) MiniProgramLogin(ctx context.Context, code string) (entity.User, error) {
	wxSession, err := u.wxApp.Code2Session(ctx, code)
	glog.Infof("wxSession: %+v - %+v - %+v", code, wxSession, err)
	if err != nil {
		return entity.User{}, app_code.New(app_code.ErrorWxAuthFailed, "微信授权失败 "+err.Error())
	}

	user, exist, err := u.repo.GetByOpenId(ctx, wxSession.OpenID)
	if exist {
		return user, err
	}

	if err != nil {
		return entity.User{}, err
	}

	// 不存在则创建用户
	user = entity.User{
		Openid:   wxSession.OpenID,
		Username: "未命名用户",
		Status:   entity.UserStatusFrozen,
	}

	err = u.repo.Store(ctx, &user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
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
	err = u.repo.Store(ctx, &user)
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.UserSetting{}, nil
	}

	if err != nil {
		return entity.UserSetting{}, err
	}
	return setting, nil
}

// GetUserRole 获取用户角色
func (u UserUseCase) GetUserRole(ctx context.Context, userID int) (entity.UserRole, error) {
	return u.repo.GetUserRole(ctx, userID)
}

// SetUserRole 设置用户角色
func (u UserUseCase) SetUserRole(ctx context.Context, userID int, role string) error {
	return u.repo.SetUserRole(ctx, userID, role)
}

// FindUsersByName 查询用户
func (u UserUseCase) FindUsersByName(ctx context.Context, username string) ([]entity.User, error) {
	return u.repo.FindUsersByName(ctx, username)
}
