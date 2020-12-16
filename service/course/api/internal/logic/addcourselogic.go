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
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"hey-go-zero/common/errorx"
	"hey-go-zero/service/course/api/internal/svc"
	"hey-go-zero/service/course/api/internal/types"
	"hey-go-zero/service/course/model"

	"github.com/tal-tech/go-zero/core/logx"
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

func (l *AddCourseLogic) AddCourse(req types.AddCourseReq) error {
	if err := l.parametersCheck(req); err != nil {
		return err
	}

	// 如果数量小于等于0则为不限
	if req.MemberLimit < 0 {
		req.MemberLimit = 0
	}

	_, err := l.svcCtx.CourseModel.FindOneByName(req.Name)
	switch err {
	case nil:
		return errorx.NewDescriptionError("课程已存在")
	case model.ErrNotFound:
		_, err = l.svcCtx.CourseModel.Insert(model.Course{
			Name:        req.Name,
			Description: req.Description,
			Classify:    req.Classify,
			GenderLimit: int64(req.GenderLimit),
			MemberLimit: int64(req.MemberLimit),
			StartTime:   req.StartTime,
			Credit:      int64(req.Credit),
		})
		return err
	default:
		return err
	}
}

func (l *AddCourseLogic) parametersCheck(req types.AddCourseReq) error {
	wordLimitErr := func(key string, limit int) error {
		return errorx.NewDescriptionError(fmt.Sprintf("%s不能超过%d个字符", key, limit))
	}

	if len(strings.TrimSpace(req.Name)) == 0 {
		return errorx.NewInvalidParameterError("name")
	}

	if utf8.RuneCountInString(req.Name) > 20 {
		return wordLimitErr("课程名称", 20)
	}

	if utf8.RuneCountInString(req.Description) > 500 {
		return wordLimitErr("课程描述", 500)
	}

	now := time.Now().AddDate(0, 0, 1)
	validEarliestStartTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
	if req.StartTime < validEarliestStartTime.Unix() {
		return errorx.NewDescriptionError(fmt.Sprintf("开课时间不能早于%s", validEarliestStartTime.Format("2006年01月02日 03时04分05秒")))
	}

	return nil
}
