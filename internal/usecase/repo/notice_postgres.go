package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

func NewNoticeRepo(pg *postgres.Postgres) *NoticeRepo {
	return &NoticeRepo{pg}
}

type NoticeRepo struct {
	*postgres.Postgres
}

func (n NoticeRepo) SetNoticeRead(ctx context.Context, noticeID int) error {
	return n.Db.WithContext(ctx).Model(&entity.Notice{}).Where("id = ?", noticeID).Update("status", entity.NotifyStatusRead).Error
}

func (n NoticeRepo) AddNotice(ctx context.Context, notice entity.Notice) error {
	return n.Db.Create(&notice).Error
}

func (n NoticeRepo) GetNoticeListByUser(ctx context.Context, userID int) ([]entity.Notice, error) {
	var notices []entity.Notice
	err := n.Db.WithContext(ctx).Where("user_id = ? and status = ?", userID, entity.NotifyStatusUnread).Find(&notices).Error
	return notices, err
}

func (n NoticeRepo) GetNoticeByNoticeID(ctx context.Context, noticeID int) (entity.Notice, error) {
	var notice entity.Notice
	err := n.Db.WithContext(ctx).Where("id = ?", noticeID).First(&notice).Error
	return notice, err
}
