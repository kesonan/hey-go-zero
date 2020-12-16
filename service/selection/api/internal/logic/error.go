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

import "hey-go-zero/common/errorx"

var (
	errCourseNotFound          = errorx.NewDescriptionError("课程不存在")
	errSelectionNotFound       = errorx.NewDescriptionError("选课任务不存在")
	errSelectionIsExists       = errorx.NewDescriptionError("选课任务已存在")
	errNotInSelectionTimeRange = errorx.NewDescriptionError("选课未开始")
	errSelectionExpired        = errorx.NewDescriptionError("选课已结束")
	errSelectionCourseNotFound = errorx.NewDescriptionError("不存在该选课课程")
	errSelectionSelected       = errorx.NewDescriptionError("你已选择过该课程，请勿重复选择")
	errSelectionNotSelected    = errorx.NewDescriptionError("你尚未选择该课程，无法取消选课")
)
