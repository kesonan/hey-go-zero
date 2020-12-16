package logic

import (
	"context"

	"hey-go-zero/service/course/rpc/course"
	"hey-go-zero/service/course/rpc/internal/svc"

	"github.com/tal-tech/go-zero/core/fx"
	"github.com/tal-tech/go-zero/core/logx"
)

type FindByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdsLogic {
	return &FindByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  批量获取课程
func (l *FindByIdsLogic) FindByIds(in *course.IdsReq) (*course.CourseListReply, error) {
	var resp course.CourseListReply
	fx.From(func(source chan<- interface{}) {
		for _, each := range in.Ids {
			source <- each
		}
	}).Split(2000).ForEach(func(item interface{}) {
		chunks, ok := item.([]interface{})
		if !ok {
			return
		}

		var ids []int64
		for _, chunk := range chunks {
			id, ok := chunk.(int64)
			if !ok {
				continue
			}

			ids = append(ids, id)
		}

		if len(ids) == 0 {
			return
		}

		users, err := l.svcCtx.CourseModel.FindByIds(ids)
		if err != nil {
			logx.Error(err)
			return
		}

		for _, each := range users {
			resp.List = append(resp.List, convertCourseFromDbToPb(each))
		}
	})

	return &resp, nil
}
