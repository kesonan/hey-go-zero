package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type MineSelectionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMineSelectionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) MineSelectionsLogic {
	return MineSelectionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MineSelectionsLogic) MineSelections() (*types.MineCourseReply, error) {
	// todo: add your logic here and delete this line

	return &types.MineCourseReply{}, nil
}
