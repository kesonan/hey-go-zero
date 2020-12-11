//  Copyright [2020] [hey-go-zero]
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package logic

import (
	"context"

	"hey-go-zero/common/errorx"
	"hey-go-zero/service/course/api/internal/svc"
	"hey-go-zero/service/course/api/internal/types"
	"hey-go-zero/service/course/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetCourseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCourseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCourseInfoLogic {
	return GetCourseInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCourseInfoLogic) GetCourseInfo(req types.CourseInfoReq) (*types.CourseInfoReply, error) {
	if req.Id <= 0 {
		return nil, errorx.NewInvalidParameterError("id")
	}

	data, err := l.svcCtx.CourseModel.FindOne(req.Id)
	switch err {
	case nil:
		return &types.CourseInfoReply{
			Id:     data.Id,
			Course: convertFromDbToLogic(*data),
		}, nil
	case model.ErrNotFound:
		return nil, errCourseNotFound
	default:
		return nil, err
	}
}
