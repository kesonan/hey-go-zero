package logic

import (
	"context"

	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"

	"github.com/tal-tech/go-zero/core/fx"
	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteSelectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSelectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteSelectionLogic {
	return DeleteSelectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSelectionLogic) DeleteSelection(req types.SelectionIdReq) error {
	err := l.svcCtx.SelectionModel.Delete(req.Id)
	switch err {
	case nil:
		selectCourseList, err := l.svcCtx.SelectionCourseModel.FindBySelectionId(req.Id)
		if err != nil {
			return err
		}

		fx.From(func(source chan<- interface{}) {
			for _, item := range selectCourseList {
				source <- item
			}
		}).Walk(func(item interface{}, pipe chan<- interface{}) {
			data := item.(*model.SelectionCourse)
			list, err := l.svcCtx.SelectionStudentModel.FindBySelectionCourseId(data.Id)
			if err != nil {
				return
			}

			for _, each := range list {
				_ = l.svcCtx.SelectionStudentModel.Delete(each.Id)
			}
		}).Done()

		_ = l.svcCtx.SelectionCourseModel.DeleteBySelectionId(req.Id)

		return nil
	case model.ErrNotFound:
		return errSelectionNotFound
	default:
		return err
	}
}
