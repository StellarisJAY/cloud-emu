package biz

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Notification struct {
	NotificationId int64
	Type           int32
	SenderId       int64
	SenderUserName string
	SenderNickName string
	Content        string
	AddTime        time.Time
	ReceiverId     int64
}

type NotificationRepo interface {
	Create(ctx context.Context, n *Notification) error
	ListInbox(ctx context.Context, userId int64, p *common.Pagination) ([]*Notification, error)
	Delete(ctx context.Context, notificationIds []int64) error
	ClearInbox(ctx context.Context, userId int64) error
}

type NotificationUseCase struct {
	notificationRepo NotificationRepo
	logger           *log.Helper
}

const (
	NotificationTypeSystem int32 = iota + 1
	NotificationTypeInvitation
)

func NewNotificationUseCase(notificationRepo NotificationRepo, logger log.Logger) *NotificationUseCase {
	return &NotificationUseCase{notificationRepo: notificationRepo, logger: log.NewHelper(logger)}
}

func (uc *NotificationUseCase) ListInbox(ctx context.Context, userId int64, p *common.Pagination) ([]*Notification, error) {
	return uc.notificationRepo.ListInbox(ctx, userId, p)
}

func (uc *NotificationUseCase) Delete(ctx context.Context, notificationIds []int64) error {
	err := uc.notificationRepo.Delete(ctx, notificationIds)
	if err != nil {
		uc.logger.Error("删除消息错误", err)
		return v1.ErrorServiceError("删除失败")
	}
	return nil
}

func (uc *NotificationUseCase) ClearInbox(ctx context.Context, userId int64) error {
	err := uc.notificationRepo.ClearInbox(ctx, userId)
	if err != nil {
		uc.logger.Error("清空消息错误", err)
		return v1.ErrorServiceError("清空失败")
	}
	return nil
}
