package usecase

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
)

type NoticeUseCase struct {
	noticeRepo NoticeRepo
}

func NewNoticeUseCase(r NoticeRepo) *NoticeUseCase {
	return &NoticeUseCase{
		noticeRepo: r,
	}
}

func (n NoticeUseCase) AddNotice(ctx context.Context, notice entity.Notice) error {
	return n.noticeRepo.AddNotice(ctx, notice)
}

func (n NoticeUseCase) GetNoticeListByUser(ctx context.Context, userID int) ([]entity.Notice, error) {
	return n.noticeRepo.GetNoticeListByUser(ctx, userID)
}

func (n NoticeUseCase) GetNoticeByNoticeID(ctx context.Context, noticeID int) (entity.Notice, error) {
	return n.noticeRepo.GetNoticeByNoticeID(ctx, noticeID)
}
