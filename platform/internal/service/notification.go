package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
)

type NotificationService struct {
	v1.UnimplementedNotificationServer
	notificationUC *biz.NotificationUseCase
}

func NewNotificationService(notificationUC *biz.NotificationUseCase) v1.NotificationServer {
	return &NotificationService{notificationUC: notificationUC}
}

func (n *NotificationService) ListInboxNotifications(ctx context.Context, request *v1.ListInboxNotificationRequest) (*v1.ListInboxNotificationResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	page := &common.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	notifications, err := n.notificationUC.ListInbox(ctx, claims.UserId, page)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.NotificationDto, len(notifications))
	for i, notification := range notifications {
		result[i] = &v1.NotificationDto{
			NotificationId: notification.NotificationId,
			SenderId:       notification.SenderId,
			SenderNickName: notification.SenderNickName,
			SenderUserName: notification.SenderUserName,
			AddTime:        notification.AddTime.Format("2006-01-02 15:04:05"),
			Content:        notification.Content,
			Type:           notification.Type,
		}
	}
	return &v1.ListInboxNotificationResponse{
		NotificationList: result,
		Total:            page.Total,
	}, nil
}
