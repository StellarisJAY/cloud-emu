package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"time"
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
		e := errors.FromError(err)
		return &v1.ListInboxNotificationResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.NotificationDto, len(notifications))
	for i, notification := range notifications {
		result[i] = &v1.NotificationDto{
			NotificationId: notification.NotificationId,
			SenderId:       notification.SenderId,
			SenderNickName: notification.SenderNickName,
			SenderUserName: notification.SenderUserName,
			AddTime:        notification.AddTime.Format(time.DateTime),
			Content:        notification.Content,
			Type:           notification.Type,
		}
	}
	return &v1.ListInboxNotificationResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   page.Total,
	}, nil
}

func (n *NotificationService) DeleteInboxNotifications(ctx context.Context, request *v1.DeleteInboxNotificationRequest) (*v1.DeleteInboxNotificationResponse, error) {
	err := n.notificationUC.Delete(ctx, request.NotificationIds)
	if err != nil {
		e := errors.FromError(err)
		return &v1.DeleteInboxNotificationResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.DeleteInboxNotificationResponse{
		Code:    200,
		Message: "删除成功",
	}, nil
}

func (n *NotificationService) ClearInbox(ctx context.Context, request *v1.ClearInboxRequest) (*v1.ClearInboxResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := n.notificationUC.ClearInbox(ctx, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ClearInboxResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.ClearInboxResponse{
		Code:    200,
		Message: "清空成功",
	}, nil
}
