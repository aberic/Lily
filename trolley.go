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

// trolley 购物车
//
// 这里面能存放很多手提袋
//
// purses 手提袋集合
//
// b+tree 模型 degree=128;level=4;nodes=[d^l]/(d-1)=2113665;
//
// node 内范围控制数量 key=127
//
// tree 内范围控制数量 n*k=268435455
type trolley struct {
	index       uint8  // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	flexibleKey uint32 // 下一级最左最小树所对应真实key
	mall        *mall
	purses      []*purse
}

func (t *trolley) put(originalKey Key, key uint32, value interface{}) error {
	index := uint8(t.mall.flexibleKey / purseDistance)
	t.flexibleKey = t.mall.flexibleKey - uint32(index)*purseDistance
	//log.Self.Debug("trolley", log.Uint32("key", key), log.Uint32("index", index))
	data := t.createChild(uint8(index))
	return data.put(originalKey, key, value)
}

func (t *trolley) get(originalKey Key, key uint32) (interface{}, error) {
	index := uint8(t.mall.flexibleKey / purseDistance)
	t.flexibleKey = t.mall.flexibleKey - uint32(index)*purseDistance
	if realIndex, err := binaryMatchData(uint8(index), t); nil == err {
		return t.purses[realIndex].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"trolley key", string(originalKey), "is nil"}, " "))
	}
}

func (t *trolley) existChild(index uint8) bool {
	return matchableData(index, t)
}

func (t *trolley) createChild(index uint8) database {
	if realIndex, err := binaryMatchData(index, t); nil != err {
		p := &purse{
			index:   index,
			trolley: t,
			box:     []*box{},
		}
		lenCity := len(t.purses)
		if lenCity == 0 {
			t.purses = append(t.purses, p)
			return p
		}
		t.purses = append(t.purses, nil)
		for i := len(t.purses) - 2; i >= 0; i-- {
			if t.purses[i].index < index {
				t.purses[i+1] = p
				break
			} else if t.purses[i].index > index {
				t.purses[i+1] = t.purses[i]
				t.purses[i] = p
			} else {
				return t.purses[i]
			}
		}
		return p
	} else {
		return t.purses[realIndex]
	}
}

func (t *trolley) childCount() int {
	return len(t.purses)
}

func (t *trolley) child(index int) nodeIndex {
	return t.purses[index]
}

func (t *trolley) getIndex() uint8 {
	return t.index
}
