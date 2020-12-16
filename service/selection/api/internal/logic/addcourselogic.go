package logic

import (
	"context"

	"hey-go-zero/common/errorx"
	"hey-go-zero/service/course/rpc/courseservice"
	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"
	"hey-go-zero/service/user/rpc/userservice"

	"github.com/tal-tech/go-zero/core/fx"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (l *AddCourseLogic) AddCourse(req types.SelectionCourseReq) error {
	if req.SelectionId <= 0 {
		return errorx.NewInvalidParameterError("selectionId")
	}

	if len(req.List) == 0 {
		return errorx.NewInvalidParameterError("list")
	}

	for _, item := range req.List {
		if item.TeacherId <= 0 {
			return errorx.NewInvalidParameterError("teacherId")
		}

		if item.CourseId <= 0 {
			return errorx.NewInvalidParameterError("courseId")
		}
	}

	_, err := l.svcCtx.SelectionModel.FindOne(req.SelectionId)
	if err != nil {
		if err == model.ErrNotFound {
			return errSelectionNotFound
		}
		return err
	}

	selectionCourseList, err := l.svcCtx.SelectionCourseModel.FindBySelectionId(req.SelectionId)
	if err != nil {
		return err
	}

	selectionCourseM := make(map[int64]struct{})
	for _, item := range selectionCourseList {
		selectionCourseM[item.CourseId] = struct{}{}
	}

	err = mr.MapReduceVoid(func(source chan<- interface{}) {
		for _, item := range req.List {
			source <- item
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		data := item.(*types.SelectionCourse)
		if _, ok := selectionCourseM[data.CourseId]; ok {
			cancel(errorx.NewDescriptionError("已经添加过该课程，请勿重复添加"))
			return

		}
		_, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{
			Id: data.CourseId,
		})
		if err != nil {
			st := status.Convert(err)
			if st.Code() == codes.NotFound {
				cancel(errCourseNotFound)
				return
			}
			cancel(err)
			return
		}

		_, err = l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{
			Id: data.TeacherId,
		})
		if err != nil {
			st := status.Convert(err)
			if st.Code() == codes.NotFound {
				cancel(errCourseNotFound)
				return
			}
			cancel(err)
			return
		}

		writer.Write(struct{}{})
	}, func(pipe <-chan interface{}, cancel func(error)) {
		for range pipe {
		}
	})

	if err != nil {
		return err
	}

	fx.From(func(source chan<- interface{}) {
		for _, item := range req.List {
			source <- item
		}
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		data := item.(*types.SelectionCourse)
		// 不考虑事务
		_, err := l.svcCtx.SelectionCourseModel.Insert(model.SelectionCourse{
			SelectionId: req.SelectionId,
			CourseId:    data.CourseId,
			TeacherId:   data.TeacherId,
		})
		logx.Error(err)
	}).Done()

	return nil
}
