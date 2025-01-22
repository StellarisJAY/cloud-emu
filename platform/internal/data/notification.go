package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"time"
)

type Notification struct {
	NotificationId int64
	Type           int32
	SenderId       int64
	Content        string
	AddTime        time.Time
	ReceiverId     int64
}

type NotificationRepo struct {
	data *Data
}

const NotificationTableName = "notification"

func NewNotificationRepo(data *Data) biz.NotificationRepo {
	return &NotificationRepo{data: data}
}

func (nr *NotificationRepo) Create(ctx context.Context, n *biz.Notification) error {
	return nr.data.DB(ctx).Table(NotificationTableName).Create(convertNotificationBizToEntity(n)).WithContext(ctx).Error
}

func (nr *NotificationRepo) ListInbox(ctx context.Context, userId int64, p *common.Pagination) ([]*biz.Notification, error) {
	var result []*biz.Notification
	err := nr.data.DB(ctx).Table(NotificationTableName+" n").
		Joins("LEFT JOIN sys_user su ON su.user_id = n.sender_id").
		Select("n.*, su.user_name AS 'sender_user_name', su.nick_name AS 'sender_nick_name'").
		Where("receiver_id = ?", userId).
		Scopes(common.WithPagination(p)).
		WithContext(ctx).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (nr *NotificationRepo) Delete(ctx context.Context, notificationIds []int64) error {
	return nr.data.DB(ctx).Table(NotificationTableName).
		WithContext(ctx).
		Where("notification_id IN ?", notificationIds).
		Delete(&Notification{}).
		Error
}

func (nr *NotificationRepo) ClearInbox(ctx context.Context, userId int64) error {
	return nr.data.DB(ctx).Table(NotificationTableName).
		WithContext(ctx).
		Where("receiver_id = ?", userId).
		Delete(&Notification{}).
		Error
}

func convertNotificationBizToEntity(notification *biz.Notification) *Notification {
	return &Notification{
		NotificationId: notification.NotificationId,
		Type:           notification.Type,
		SenderId:       notification.SenderId,
		Content:        notification.Content,
		AddTime:        notification.AddTime,
		ReceiverId:     notification.ReceiverId,
	}
}
