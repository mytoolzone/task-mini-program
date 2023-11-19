package wechat

import (
	"context"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

type UserInfo struct {
	OpenID          string `json:"openId"`
	UnionID         string `json:"unionId"`
	NickName        string `json:"nickName"`
	Gender          int    `json:"gender"`
	City            string `json:"city"`
	Province        string `json:"province"`
	Country         string `json:"country"`
	AvatarURL       string `json:"avatarUrl"`
	Language        string `json:"language"`
	PhoneNumber     string `json:"phoneNumber"`
	OpenGID         string `json:"openGId"`
	MsgTicket       string `json:"msgTicket"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

// SessionInfo 登录凭证校验的返回结果
type SessionInfo struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足UnionID下发条件的情况下会返回
}

type WxApp interface {
	// Code2Session 获取用户的 openid
	Code2Session(ctx context.Context, code string) (SessionInfo, error)
	// GetUserInfo 获取用户的信息
	GetUserInfo(encryptedData, iv, sessionKey string) (UserInfo, error)
	// GetPhoneNumber 获取用户的手机号信息
	GetPhoneNumber(ctx context.Context, code string) (string, error)
}

var _ WxApp = (*MiniProgram)(nil)

type MiniProgram struct {
	AppID string
	mini  *miniprogram.MiniProgram
}

func NewMiniProgram(appID, appSecret string) *MiniProgram {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &miniConfig.Config{
		AppID:     appID,
		AppSecret: appSecret,
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)

	return &MiniProgram{
		AppID: appID,
		mini:  mini,
	}
}

func (m MiniProgram) Code2Session(ctx context.Context, jsCode string) (SessionInfo, error) {
	sessionContext, err := m.mini.GetAuth().Code2SessionContext(ctx, jsCode)

	if err != nil {
		return SessionInfo{}, err
	}

	return SessionInfo{
		OpenID:     sessionContext.OpenID,
		SessionKey: sessionContext.SessionKey,
		UnionID:    sessionContext.UnionID,
	}, nil
}

func (m MiniProgram) GetUserInfo(sessionKey, encryptedData, iv string) (UserInfo, error) {
	decrypt, err := m.mini.GetEncryptor().Decrypt(sessionKey, encryptedData, iv)
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{
		OpenID:          decrypt.OpenID,
		UnionID:         decrypt.UnionID,
		NickName:        decrypt.NickName,
		Gender:          decrypt.Gender,
		City:            decrypt.City,
		Province:        decrypt.Province,
		Country:         decrypt.Country,
		AvatarURL:       decrypt.AvatarURL,
		Language:        decrypt.Language,
		PhoneNumber:     decrypt.PhoneNumber,
		OpenGID:         decrypt.OpenGID,
		MsgTicket:       decrypt.MsgTicket,
		PurePhoneNumber: decrypt.PurePhoneNumber,
		CountryCode:     decrypt.CountryCode,
	}, nil
}

func (m MiniProgram) GetPhoneNumber(ctx context.Context, code string) (string, error) {
	phoneInfo, err := m.mini.GetAuth().GetPhoneNumberContext(ctx, code)
	if err != nil {
		return "", err
	}
	return phoneInfo.PhoneInfo.PhoneNumber, nil
}
