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
	"hey-go-zero/service/user/api/internal/logic"
	"hey-go-zero/service/user/api/internal/svc"
	"hey-go-zero/service/user/api/internal/types"
	"hey-go-zero/service/user/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type UserInfoEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoEditLogic {
	return UserInfoEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoEditLogic) UserInfoEdit(id int64, req types.UserInfoReq) error {
	// 全量更新，允许字段为空
	resp, err := l.svcCtx.UserModel.FindOne(id)
	switch err {
	case nil:
		resp.Name = req.Name
		switch req.Gender {
		case "男":
			resp.Gender = 1
		case "女":
			resp.Gender = 2
		default:
			return errorx.NewInvalidParameterError("gender")
		}
		return l.svcCtx.UserModel.Update(*resp)
	case model.ErrNotFound:
		return logic.ErrUserNotFound
	default:
		return err
	}
}
