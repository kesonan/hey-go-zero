package logic

import (
	"context"
	"sort"

	"hey-go-zero/service/course/rpc/courseservice"
	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"
	"hey-go-zero/service/user/rpc/userservice"

	"github.com/tal-tech/go-zero/core/fx"
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

func (l *MineSelectionsLogic) MineSelections(userId int64) (*types.MineCourseReply, error) {
	studentSelectedCourseList, err := l.svcCtx.SelectionStudentModel.FindByStudentId(userId)
	if err != nil {
		return nil, err
	}

	var resp types.MineCourseReply
	fx.From(func(source chan<- interface{}) {
		for _, each := range studentSelectedCourseList {
			source <- each
		}
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		data := item.(*model.SelectionStudent)
		var teacherName string
		userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{Id: data.StudentId})
		if err != nil {
			logx.Error(errSelectionSelected)
		} else {
			teacherName = userInfo.Name
		}

		selectionCourse, err := l.svcCtx.SelectionCourseModel.FindOne(data.SelectionCourseId)
		if err != nil {
			return
		}

		courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{Id: selectionCourse.CourseId})
		if err != nil {
			return
		}

		resp.List = append(resp.List, &types.Course{
			Id:                courseInfo.Id,
			SelectionCourseId: selectionCourse.Id,
			Name:              courseInfo.Name,
			Description:       courseInfo.Description,
			Classify:          courseInfo.Classify,
			GenderLimit:       int(courseInfo.GenderLimit),
			MemberLimit:       int(courseInfo.MemberLimit),
			StartTime:         courseInfo.StartTime,
			Credit:            int(courseInfo.Credit),
			TeacherName:       teacherName,
		})
	}).Done()

	sort.Slice(resp.List, func(i, j int) bool {
		return resp.List[i].Id < resp.List[j].Id
	})

	return &resp, nil
}
