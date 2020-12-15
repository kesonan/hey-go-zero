package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

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
	// todo: add your logic here and delete this line

	return nil
}
