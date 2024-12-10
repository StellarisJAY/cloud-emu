package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type NotificationRepo struct {
	data *Data
}

const NotificationTableName = "notification"

func NewNotificationRepo(data *Data) biz.NotificationRepo {
	return &NotificationRepo{data: data}
}

func (nr *NotificationRepo) Create(ctx context.Context, n *biz.Notification) error {
	return nr.data.db.Table(NotificationTableName).Create(n).WithContext(ctx).Error
}

func (nr *NotificationRepo) ListInbox(ctx context.Context, userId int64, p *common.Pagination) ([]*biz.Notification, error) {
	var result []*biz.Notification
	err := nr.data.db.Table(NotificationTableName+" n").
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
