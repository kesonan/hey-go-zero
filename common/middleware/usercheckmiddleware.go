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

package middleware

import (
	"fmt"
	"net/http"

	"hey-go-zero/common/errorx"
	"hey-go-zero/common/jwtx"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UserCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(jwtx.JwtWithUserKey)
		xUserId := r.Header.Get("x-user-id")
		if len(xUserId) == 0 {
			httpx.Error(w, errorx.NewDescriptionError("x-user-id不能为空"))
			return
		}

		if xUserId != fmt.Sprintf("%v", v) {
			httpx.Error(w, errorx.NewDescriptionError("用户信息不一致"))
			return
		}
		next(w, r)
	}
}
