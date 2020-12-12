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
	"hey-go-zero/service/user/model"
	"hey-go-zero/service/user/rpc/user"
)

func convertUserFromDbToPb(in *model.User) *user.UserReply {
	var resp user.UserReply
	resp.Id = in.Id
	resp.Name = in.Name
	resp.Gender = user.Gender(in.Gender)
	resp.Role = in.Role
	resp.CreateTime = in.CreateTime.UnixNano() / 1e6
	resp.UpdateTime = in.UpdateTime.UnixNano() / 1e6
	return &resp
}
