package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetTeachingStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTeachingStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTeachingStudentsLogic {
	return GetTeachingStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTeachingStudentsLogic) GetTeachingStudents(req types.SelectCourseId) error {
	// todo: add your logic here and delete this line

	return nil
}
