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

package errorx

import (
	"fmt"
	"net/http"
)

const defaultCode = -1

type Handler struct{}

type ErrorBody struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func (h *Handler) Handle() func(error) (int, interface{}) {
	return func(err error) (int, interface{}) {
		switch v := err.(type) {
		case *CodeError:
			return http.StatusNotAcceptable, ErrorBody{
				Code: v.code,
				Desc: v.desc,
			}
		case *DescriptionError:
			return http.StatusNotAcceptable, ErrorBody{
				Code: defaultCode,
				Desc: v.desc,
			}
		case *InvalidParameterError:
			return http.StatusNotAcceptable, ErrorBody{
				Code: defaultCode,
				Desc: fmt.Sprintf("参数错误: %v", v.parameter),
			}
		default:
			return http.StatusInternalServerError, ErrorBody{
				Code: defaultCode,
				Desc: v.Error(),
			}
		}
	}
}

type CodeError struct {
	code int
	desc string
}

func NewCodeError(code int, desc string) *CodeError {
	return &CodeError{
		code: code,
		desc: desc,
	}
}

func (e *CodeError) Error() string {
	return e.desc
}

type DescriptionError struct {
	desc string
}

func NewDescriptionError(desc string) *DescriptionError {
	return &DescriptionError{
		desc: desc,
	}
}

func (e *DescriptionError) Error() string {
	return e.desc
}

type InvalidParameterError struct {
	parameter string
}

func NewInvalidParameterError(parameter string) *InvalidParameterError {
	return &InvalidParameterError{
		parameter: parameter,
	}
}

func (e *InvalidParameterError) Error() string {
	return e.parameter
}
