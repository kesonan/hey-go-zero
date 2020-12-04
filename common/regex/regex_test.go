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

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatch(t *testing.T) {
	assert.True(t, Match("songmeizi", Username))
	assert.True(t, Match("songmeizi.", Password))
	assert.False(t, Match("song", Username))
	assert.False(t, Match("song", Password))
	assert.False(t, Match("song*", Password))
	assert.False(t, Match("songmeizisongmeizi1", Password))
	assert.False(t, Match("songmeizisongmeizi123", Username))
	assert.False(t, Match("", Username))
	assert.False(t, Match("", Password))
}
