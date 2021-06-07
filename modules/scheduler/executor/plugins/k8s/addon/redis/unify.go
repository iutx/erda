// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package redis

import (
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/scheduler/executor/plugins/k8s/addon"
)

type UnifiedRedisOperator struct {
	redisoperator *RedisOperator
}

func New(overcommit addon.OvercommitUtil) *UnifiedRedisOperator {
	return &UnifiedRedisOperator{
		redisoperator: NewRedisOperator(overcommit),
	}
}

func (ro *UnifiedRedisOperator) IsSupported() bool {
	if ro.redisoperator.IsSupported() {
		return true
	}
	return false
}

func (ro *UnifiedRedisOperator) Validate(sg *apistructs.ServiceGroup) error {
	return ro.redisoperator.Validate(sg)
}
func (ro *UnifiedRedisOperator) Convert(sg *apistructs.ServiceGroup) interface{} {
	return ro.redisoperator.Convert(sg)
}
func (ro *UnifiedRedisOperator) Create(k8syml interface{}) error {
	return ro.redisoperator.Create(k8syml)
}

func (ro *UnifiedRedisOperator) Inspect(sg *apistructs.ServiceGroup) (*apistructs.ServiceGroup, error) {
	return ro.redisoperator.Inspect(sg)
}
func (ro *UnifiedRedisOperator) Remove(sg *apistructs.ServiceGroup) error {
	return ro.redisoperator.Remove(sg)
}

func (ro *UnifiedRedisOperator) Update(k8syml interface{}) error {
	return ro.redisoperator.Update(k8syml)
}
