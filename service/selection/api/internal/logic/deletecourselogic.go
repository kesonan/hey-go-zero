package logic

import (
	"context"
	"time"

	"hey-go-zero/common/errorx"
	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"

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

func (l *DeleteCourseLogic) DeleteCourse(req types.DeleteSelectionCourseReq) error {
	selection, err := l.svcCtx.SelectionModel.FindOne(req.SelectionId)
	switch err {
	case nil:
		if time.Now().After(time.Unix(selection.StartTime, 0)) {
			return errorx.NewDescriptionError("该选课已发布，不能编辑课程")
		}

		err = l.svcCtx.SelectionModel.Delete(req.SelectionId)
		if err != nil {
			logx.Error(err)
		}
		for _, each := range req.Ids {
			_ = l.svcCtx.SelectionCourseModel.Delete(each)
		}
		return nil
	case model.ErrNotFound:
		return errSelectionNotFound
	default:
		return err
	}
}
