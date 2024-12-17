package biz

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
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
}

type NotificationUseCase struct {
	notificationRepo NotificationRepo
}

const (
	NotificationTypeSystem int32 = iota + 1
	NotificationTypeInvitation
)

func NewNotificationUseCase(notificationRepo NotificationRepo) *NotificationUseCase {
	return &NotificationUseCase{notificationRepo: notificationRepo}
}

func (uc *NotificationUseCase) ListInbox(ctx context.Context, userId int64, p *common.Pagination) ([]*Notification, error) {
	return uc.notificationRepo.ListInbox(ctx, userId, p)
}
