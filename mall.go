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

// mall 商城
//
// 这里面能存放很多购物车，是商场提供给shopper的最大载物设备
//
// trolleys 购物车集合
//
// b+tree 模型 degree=128;level=4;nodes=[d^l]/(d-1)=2113665;
//
// node 内范围控制数量 key=127
//
// tree 内范围控制数量 n*k=268435455
type mall struct {
	index       uint8  // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	flexibleKey uint32 // 下一级最左最小树所对应真实key
	city        *city
	trolleys    []*trolley
}

func (m *mall) put(originalKey Key, key uint32, value interface{}) error {
	index := uint8(m.city.flexibleKey / trolleyDistance)
	m.flexibleKey = m.city.flexibleKey - uint32(index)*trolleyDistance
	//log.Self.Debug("city", log.Uint32("key", key), log.Uint32("index", index))
	data := m.createChild(uint8(index))
	return data.put(originalKey, key, value)
}

func (m *mall) get(originalKey Key, key uint32) (interface{}, error) {
	index := uint8(m.city.flexibleKey / trolleyDistance)
	m.flexibleKey = m.city.flexibleKey - uint32(index)*trolleyDistance
	if realIndex, err := binaryMatchData(uint8(index), m); nil == err {
		return m.trolleys[realIndex].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"mall key", string(originalKey), "is nil"}, " "))
	}
}

func (m *mall) existChild(index uint8) bool {
	return matchableData(index, m)
}

func (m *mall) createChild(index uint8) database {
	if realIndex, err := binaryMatchData(index, m); nil != err {
		t := &trolley{
			index:  index,
			mall:   m,
			purses: []*purse{},
		}
		lenCity := len(m.trolleys)
		if lenCity == 0 {
			m.trolleys = append(m.trolleys, t)
			return t
		}
		m.trolleys = append(m.trolleys, nil)
		for i := len(m.trolleys) - 2; i >= 0; i-- {
			if m.trolleys[i].index < index {
				m.trolleys[i+1] = t
				break
			} else if m.trolleys[i].index > index {
				m.trolleys[i+1] = m.trolleys[i]
				m.trolleys[i] = t
			} else {
				return m.trolleys[i]
			}
		}
		return t
	} else {
		return m.trolleys[realIndex]
	}
}

func (m *mall) childCount() int {
	return len(m.trolleys)
}

func (m *mall) child(index int) nodeIndex {
	return m.trolleys[index]
}

func (m *mall) getIndex() uint8 {
	return m.index
}
