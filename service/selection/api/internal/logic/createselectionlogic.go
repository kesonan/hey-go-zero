package logic

import (
	"context"
	"fmt"
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

		// dq
		msg := fmt.Sprintf("选课【%s】还有两小时就要开始了，请提前做好选课准备。", req.Name)
		_, err = l.svcCtx.Producer.At([]byte(msg), time.Unix(req.StartTime, 0).Add(-2*time.Hour))

		return err
	default:
		return err
	}
}
