package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetTeachingCoursesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTeachingCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTeachingCoursesLogic {
	return GetTeachingCoursesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTeachingCoursesLogic) GetTeachingCourses() (*types.MineCourseReply, error) {
	// todo: add your logic here and delete this line

	return &types.MineCourseReply{}, nil
}
