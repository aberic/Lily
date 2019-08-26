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
	key         uint32
	realKey     uint8
	flexibleKey uint32
	city        *city
	indexes     []uint8
	trolleys    map[uint8]*trolley
}

func (m *mall) put(originalKey Key, key uint32, value interface{}) error {
	m.realKey = uint8(m.city.flexibleKey / trolleyDistance)
	m.flexibleKey = m.city.flexibleKey - uint32(m.realKey)*trolleyDistance
	//log.Self.Debug("city", log.Uint32("key", key), log.Uint32("realKey", realKey))
	m.createChild(m.realKey)
	return m.trolleys[m.realKey].put(originalKey, key, value)
}

func (m *mall) get(originalKey Key, key uint32) (interface{}, error) {
	m.realKey = uint8(m.city.flexibleKey / trolleyDistance)
	m.flexibleKey = m.city.flexibleKey - uint32(m.realKey)*trolleyDistance
	if m.existChild(m.realKey) {
		return m.trolleys[m.realKey].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"mall key", string(originalKey), "is nil"}, " "))
	}
}

func (m *mall) existChild(index uint8) bool {
	return nil != m.trolleys[index]
}

func (m *mall) createChild(index uint8) {
	if !m.existChild(index) {
		m.indexes = append(m.indexes, index)
		m.trolleys[index] = &trolley{
			key:     uint32(m.city.realKey)*mallDistance + uint32(index+1)*trolleyDistance,
			mall:    m,
			indexes: []uint8{},
			purses:  map[uint8]*purse{},
		}
	}
}
