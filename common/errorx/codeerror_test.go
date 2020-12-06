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
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Handle(t *testing.T) {
	errorCode := 1001
	desc := "Something Wrong"
	ce := NewCodeError(errorCode, desc)
	handler := Handler{}
	fn := handler.Handle()
	statusCode, v := fn(ce)
	assert.Equal(t, http.StatusNotAcceptable, statusCode)
	assert.Equal(t, v, ErrorBody{errorCode, desc})

	statusCode, v = fn(NewDescriptionError(desc))
	assert.Equal(t, http.StatusNotAcceptable, statusCode)
	assert.Equal(t, v, ErrorBody{defaultCode, desc})

	statusCode, v = fn(NewInvalidParameterError("user"))
	assert.Equal(t, http.StatusNotAcceptable, statusCode)

	statusCode, v = fn(errors.New(desc))
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}
