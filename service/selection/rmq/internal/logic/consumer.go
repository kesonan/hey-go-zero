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

	"hey-go-zero/service/selection/rmq/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type ConsumerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConsumerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumerLogic {
	return &ConsumerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (c *ConsumerLogic) Start() {
	c.svcCtx.Consumer.Consume(c.Consume)
}

func (c *ConsumerLogic) Stop() {
}

func (c *ConsumerLogic) Consume(body []byte) {
	msg := string(body)
	fmt.Println(msg)
}