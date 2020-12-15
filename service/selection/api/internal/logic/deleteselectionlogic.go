package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteSelectionLogic {
	return DeleteSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSelectionLogic) DeleteSelection(req types.SelectionIdReq) error {
	// todo: add your logic here and delete this line

	return nil
}
