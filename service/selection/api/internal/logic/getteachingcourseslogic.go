package logic

import (
	"context"

	"hey-go-zero/service/course/rpc/courseservice"
	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"
	"hey-go-zero/service/user/rpc/userservice"

	"github.com/tal-tech/go-zero/core/fx"
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

func (l *GetTeachingCoursesLogic) GetTeachingCourses(userId int64) (*types.MineCourseReply, error) {
	selectCourseList, err := l.svcCtx.SelectionCourseModel.FindByTeacherId(userId)
	if err != nil {
		return nil, err
	}

	var teacherName string
	userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{Id: userId})
	if err != nil {
		logx.Error(err)
	} else {
		teacherName = userInfo.Name
	}

	var resp types.MineCourseReply
	fx.From(func(source chan<- interface{}) {
		for _, item := range selectCourseList {
			source <- item
		}
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		data := item.(*model.SelectionCourse)
		courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{Id: data.CourseId})
		if err != nil {
			logx.Error(err)
			return
		}

		resp.List = append(resp.List, &types.Course{
			Id:                courseInfo.Id,
			SelectionCourseId: data.Id,
			Name:              courseInfo.Name,
			Description:       courseInfo.Description,
			Classify:          courseInfo.Classify,
			GenderLimit:       int(courseInfo.GenderLimit),
			MemberLimit:       int(courseInfo.MemberLimit),
			Credit:            int(courseInfo.Credit),
			TeacherName:       teacherName,
		})
	}).Done()

	return &resp, nil
}
