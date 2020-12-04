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
	"time"

	"hey-go-zero/common/errorx"
	"hey-go-zero/common/jwtx"
	"hey-go-zero/common/regex"
	"hey-go-zero/service/user/api/internal/logic"
	"hey-go-zero/service/user/api/internal/svc"
	"hey-go-zero/service/user/api/internal/types"
	"hey-go-zero/service/user/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/tal-tech/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.UserLoginReq) (*types.UserLoginReply, error) {
	if !regex.Match(req.Username, regex.Username) {
		return nil, logic.InvalidUsername
	}

	if !regex.Match(req.Passowrd, regex.Password) {
		return nil, logic.InvalidPassword
	}

	resp, err := l.svcCtx.UserModel.FindOneByUsername(req.Username)
	switch err {
	case nil:
		if resp.Password != req.Passowrd {
			return nil, errorx.NewDescriptionError("密码错误")
		}

		jwtToken, expireAt, err := l.generateJwtToken(resp.Id, time.Now().Unix())
		if err != nil {
			return nil, err
		}

		return &types.UserLoginReply{
			Id:       resp.Id,
			Token:    jwtToken,
			ExpireAt: expireAt,
		}, nil
	case model.ErrNotFound:
		return nil, errorx.NewDescriptionError("用户名未注册")
	default:
		return nil, err
	}
}

func (l *LoginLogic) generateJwtToken(id int64, iat int64) (string, int64, error) {
	claims := make(jwt.MapClaims)
	expireAt := iat + l.svcCtx.Config.Auth.AccessExpire
	claims["exp"] = expireAt
	claims["iat"] = iat
	claims[jwtx.JwtWithUserKey] = id
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	jwtToken, err := token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
	if err != nil {
		return "", 0, err
	}
	return jwtToken, expireAt, nil
}
