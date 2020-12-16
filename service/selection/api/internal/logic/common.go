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
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"hey-go-zero/common/errorx"
	"hey-go-zero/service/selection/api/internal/types"
)

func checkCourseSelection(in types.CreateSelectionReq) error {
	if len(strings.TrimSpace(in.Name)) == 0 {
		return errorx.NewInvalidParameterError("name")
	}

	if utf8.RuneCountInString(in.Name) > 20 {
		return lengthAlert("课程名称", 20)
	}

	now := time.Now()
	startTime := now.Add(2 * time.Hour)
	endTime := startTime.AddDate(0, 0, 5)
	if in.StartTime < startTime.Unix() {
		return errorx.NewDescriptionError(fmt.Sprintf("选课开始时间不能早于%s", startTime.Format("2006年01月02日 03时04分05秒")))
	}

	if in.EndTime > endTime.Unix() {
		return errorx.NewDescriptionError(fmt.Sprintf("选课结束时间不能晚于%s", startTime.Format("2006年01月02日 03时04分05秒")))
	}

	if utf8.RuneCountInString(in.Notification) > 500 {
		return lengthAlert("选课通知", 500)
	}

	return nil
}

func lengthAlert(hint string, length int) error {
	return errorx.NewDescriptionError(fmt.Sprintf("%s不能超过%d个字符", hint, length))
}
