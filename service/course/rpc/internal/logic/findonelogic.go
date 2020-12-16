package logic

import (
	"context"

	"hey-go-zero/service/course/model"
	"hey-go-zero/service/course/rpc/course"
	"hey-go-zero/service/course/rpc/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FindOneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneLogic {
	return &FindOneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  查询课程
func (l *FindOneLogic) FindOne(in *course.IdReq) (*course.Course, error) {
	data, err := l.svcCtx.CourseModel.FindOne(in.Id)
	switch err {
	case nil:
		return convertCourseFromDbToPb(data), nil
	case model.ErrNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Unknown, err.Error())
	}
}
