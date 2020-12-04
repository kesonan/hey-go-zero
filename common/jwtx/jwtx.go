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

package jwtx

import (
	"encoding/json"
	"net/http"

	"hey-go-zero/common/errorx"

	"github.com/tal-tech/go-zero/rest/httpx"
)

const JwtWithUserKey = "id"

func GetUserId(w http.ResponseWriter, r *http.Request) (int64, bool) {
	v := r.Context().Value(JwtWithUserKey)
	jn, ok := v.(json.Number)
	if !ok {
		httpx.Error(w, errorx.NewDescriptionError("用户信息获取失败"))
		return 0, false
	}
	vInt, err := jn.Int64()
	if err != nil {
		httpx.Error(w, errorx.NewDescriptionError(err.Error()))
		return 0, false
	}
	return vInt, true
}
