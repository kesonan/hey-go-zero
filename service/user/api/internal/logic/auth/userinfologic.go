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

	"hey-go-zero/service/user/api/internal/logic"
	"hey-go-zero/service/user/api/internal/svc"
	"hey-go-zero/service/user/api/internal/types"
	"hey-go-zero/service/user/model"

	"github.com/tal-tech/go-zero/core/logx"
)

var genderConvert = map[int64]string{
	0: "未知",
	1: "男",
	2: "女",
}

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoLogic {
	return UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(id int64) (*types.UserInfoReply, error) {
	resp, err := l.svcCtx.UserModel.FindOne(id)
	switch err {
	case nil:
		return &types.UserInfoReply{
			Id:     resp.Id,
			Name:   resp.Name,
			Gender: genderConvert[resp.Gender],
			Role:   resp.Role,
		}, nil
	case model.ErrNotFound:
		return nil, logic.ErrUserNotFound
	default:
		return nil, err
	}
}
