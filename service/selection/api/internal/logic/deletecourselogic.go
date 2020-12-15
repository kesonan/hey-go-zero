package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteCourseLogic {
	return DeleteCourseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCourseLogic) DeleteCourse(req types.SelectionAddCourseReq) error {
	// todo: add your logic here and delete this line

	return nil
}
