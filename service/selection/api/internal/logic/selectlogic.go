package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) SelectLogic {
	return SelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectLogic) Select(req types.SelectCourseId) error {
	// todo: add your logic here and delete this line

	return nil
}
