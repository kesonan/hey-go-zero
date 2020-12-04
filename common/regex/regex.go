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

package regex

import "regexp"

const (
	Username = `(?m)[a-zA-Z_0-9]{6,20}`
	Password = `(?m)[a-zA-Z_0-9.-]{6,18}`
)

func Match(s, reg string) bool {
	r := regexp.MustCompile(reg)
	ret := r.FindString(s)
	return ret == s && r.MatchString(s)
}
