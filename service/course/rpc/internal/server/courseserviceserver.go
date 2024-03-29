// Code generated by goctl. DO NOT EDIT!
// Source: course.proto

package server

import (
	"context"

	"hey-go-zero/service/course/rpc/course"
	"hey-go-zero/service/course/rpc/internal/logic"
	"hey-go-zero/service/course/rpc/internal/svc"
)

type CourseServiceServer struct {
	svcCtx *svc.ServiceContext
}

func NewCourseServiceServer(svcCtx *svc.ServiceContext) *CourseServiceServer {
	return &CourseServiceServer{
		svcCtx: svcCtx,
	}
}

//  查询课程
func (s *CourseServiceServer) FindOne(ctx context.Context, in *course.IdReq) (*course.Course, error) {
	l := logic.NewFindOneLogic(ctx, s.svcCtx)
	return l.FindOne(in)
}

//  批量获取课程
func (s *CourseServiceServer) FindByIds(ctx context.Context, in *course.IdsReq) (*course.CourseListReply, error) {
	l := logic.NewFindByIdsLogic(ctx, s.svcCtx)
	return l.FindByIds(in)
}
