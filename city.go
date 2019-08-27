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
	index       uint8  // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	flexibleKey uint32 // 下一级最左最小树所对应真实key
	lily        *lily
	malls       []*mall
}

func (c *city) put(originalKey Key, key uint32, value interface{}) error {
	index := uint8(key / mallDistance)
	c.flexibleKey = key - uint32(index)*mallDistance
	//log.Self.Debug("city", log.Uint32("key", key), log.Uint32("index", index))
	data := c.createChild(uint8(index))
	return data.put(originalKey, key, value)
}

func (c *city) get(originalKey Key, key uint32) (interface{}, error) {
	index := uint8(key / mallDistance)
	c.flexibleKey = key - uint32(index)*mallDistance
	if realIndex, err := binaryMatchData(uint8(index), c); nil == err {
		return c.malls[realIndex].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"city key", string(originalKey), "is nil"}, " "))
	}
}

func (c *city) existChild(index uint8) bool {
	return matchableData(index, c)
}

func (c *city) createChild(index uint8) database {
	if realIndex, err := binaryMatchData(index, c); nil != err {
		m := &mall{
			index:    index,
			city:     c,
			trolleys: []*trolley{},
		}
		lenCity := len(c.malls)
		if lenCity == 0 {
			c.malls = append(c.malls, m)
			return m
		}
		c.malls = append(c.malls, nil)
		for i := len(c.malls) - 2; i >= 0; i-- {
			if c.malls[i].index < index {
				c.malls[i+1] = m
				break
			} else if c.malls[i].index > index {
				c.malls[i+1] = c.malls[i]
				c.malls[i] = m
			} else {
				return c.malls[i]
			}
		}
		return m
	} else {
		return c.malls[realIndex]
	}
}

func (c *city) childCount() int {
	return len(c.malls)
}

func (c *city) child(index int) nodeIndex {
	return c.malls[index]
}

func (c *city) getIndex() uint8 {
	return c.index
}
