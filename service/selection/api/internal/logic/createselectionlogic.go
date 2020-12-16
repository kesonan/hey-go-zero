package logic

import (
	"context"
	"time"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type CreateSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateSelectionLogic {
	return CreateSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSelectionLogic) CreateSelection(req types.CreateSelectionReq) error {
	if err := checkCourseSelection(req); err != nil {
		return err
	}

	_, err := l.svcCtx.SelectionModel.FindOneByName(req.Name)
	switch err {
	case nil:
		return errSelectionIsExists
	case model.ErrNotFound:
		_, err := l.svcCtx.SelectionModel.Insert(model.Selection{
			MaxCredit:    int64(req.MaxCredit),
			StartTime:    req.StartTime,
			EndTime:      req.EndTime,
			Notification: req.Notification,
			Name:         req.Name,
		})
		if err != nil {
			return err
		}

		// dq，todo：这里建议用cron-job替代，如果用dq对于这种需要变更时间的逻辑，将导致发送了多个不同时间点的message，本案例仅用于演示dq这么使用。
		_, err = l.svcCtx.Producer.At([]byte(req.Notification), time.Unix(req.StartTime, 0).Add(-2*time.Hour))

		return err
	default:
		return err
	}
}
