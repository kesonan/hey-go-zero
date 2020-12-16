package logic

import (
	"context"
	"time"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type EditSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditSelectionLogic {
	return EditSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditSelectionLogic) EditSelection(req types.EditSelectionReq) error {
	if err := checkCourseSelection(req.CreateSelectionReq); err != nil {
		return err
	}
	data, err := l.svcCtx.SelectionModel.FindOne(req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return errSelectionNotFound
		}
		return err
	}

	nameData, err := l.svcCtx.SelectionModel.FindOneByName(req.Name)
	if err != nil {
		if err == model.ErrNotFound {
			return err
		}
	} else {
		if nameData.Id != req.Id {
			return errSelectionIsExists
		}
	}

	data.Name = req.Name
	data.MaxCredit = int64(req.MaxCredit)
	data.StartTime = req.StartTime
	data.EndTime = req.EndTime
	data.Notification = req.Notification
	err = l.svcCtx.SelectionModel.Update(*data)
	if err != nil {
		return err
	}

	// dq，todo：这里建议用cron-job替代，如果用dq对于这种需要变更时间的逻辑，将导致发送了多个不同时间点的message，本案例仅用于演示dq这么使用。
	_, err = l.svcCtx.Producer.At([]byte(req.Notification), time.Unix(req.StartTime, 0).Add(-2*time.Hour))
	return err
}
