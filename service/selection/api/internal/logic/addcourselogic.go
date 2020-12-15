package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddCourseLogic {
	return AddCourseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCourseLogic) AddCourse(req types.SelectionAddCourseReq) error {
	// todo: add your logic here and delete this line

	return nil
}
