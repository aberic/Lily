/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package Lily

import (
	"errors"
	"strings"
)

// city 城市
//
// 这里面能存放很多商城
//
// malls 商城集合
//
// b+tree 模型 degree=128;level=4;nodes=[d^l]/(d-1)=2113665;
//
// node 内范围控制数量 key=127
//
// tree 内范围控制数量 n*k=268435455
type city struct {
	key         uint32
	realKey     uint8
	flexibleKey uint32
	lily        *lily
	indexes     []uint8
	malls       map[uint8]*mall
}

func (c *city) put(originalKey Key, key uint32, value interface{}) error {
	c.realKey = uint8(key / mallDistance)
	c.flexibleKey = key - uint32(c.realKey)*mallDistance
	//log.Self.Debug("city", log.Uint32("key", key), log.Uint32("realKey", realKey))
	c.createChild(c.realKey)
	return c.malls[c.realKey].put(originalKey, key, value)
}

func (c *city) get(originalKey Key, key uint32) (interface{}, error) {
	c.realKey = uint8(key / mallDistance)
	c.flexibleKey = key - uint32(c.realKey)*mallDistance
	if c.existChild(c.realKey) {
		return c.malls[c.realKey].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"city key", string(originalKey), "is nil"}, " "))
	}
}

func (c *city) existChild(index uint8) bool {
	return nil != c.malls[index]
}

func (c *city) createChild(index uint8) {
	if !c.existChild(index) {
		c.indexes = append(c.indexes, index)
		c.malls[index] = &mall{
			key:      uint32(index+1) * mallDistance,
			city:     c,
			indexes:  []uint8{},
			trolleys: map[uint8]*trolley{},
		}
	}
}
