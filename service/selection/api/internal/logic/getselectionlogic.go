package logic

import (
	"context"
	"sort"

	"hey-go-zero/service/course/rpc/courseservice"
	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"
	"hey-go-zero/service/user/rpc/user"

	"github.com/tal-tech/go-zero/core/fx"
	"github.com/tal-tech/go-zero/core/logx"
)

type GetSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSelectionLogic {
	return GetSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelectionLogic) GetSelection(req types.SelectionIdReq) (*types.SelectionReply, error) {
	list, err := l.svcCtx.SelectionCourseModel.FindBySelectionId(req.Id)
	if err != nil {
		return nil, err
	}

	var courseList []*types.Course
	fx.From(func(source chan<- interface{}) {
		for _, item := range list {
			source <- item
		}
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		data := item.(*model.SelectionCourse)
		courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{
			Id: data.CourseId,
		})
		if err != nil {
			logx.Error(err)
			return
		}

		var teacherName string
		userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &user.UserReq{Id: data.TeacherId})
		if err != nil {
			logx.Error(err)
		} else {
			teacherName = userInfo.Name
		}

		courseList = append(courseList, &types.Course{
			Id:                courseInfo.Id,
			SelectionCourseId: data.Id,
			Name:              courseInfo.Name,
			Description:       courseInfo.Description,
			Classify:          courseInfo.Classify,
			GenderLimit:       int(courseInfo.GenderLimit),
			MemberLimit:       int(courseInfo.MemberLimit),
			StartTime:         courseInfo.StartTime,
			Credit:            int(courseInfo.Credit),
			TeacherName:       teacherName,
		})
	})

	// sort by id desc
	sort.Slice(courseList, func(i, j int) bool {
		return courseList[i].Id > courseList[j].Id
	})

	data, err := l.svcCtx.SelectionModel.FindOne(req.Id)
	switch err {
	case nil:
		return &types.SelectionReply{
			Id:           data.Id,
			Name:         data.Name,
			MaxCredit:    int(data.MaxCredit),
			StartTime:    data.StartTime,
			EndTime:      data.EndTime,
			Notification: data.Notification,
			CourseList:   nil,
		}, nil
	case model.ErrNotFound:
		return nil, errSelectionNotFound
	default:
		return nil, err
	}
}
