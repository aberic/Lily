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

// purse 手提袋
//
// 这里面能存放很多个包装盒
//
// box 包装盒集合
type purse struct {
	key         uint32
	flexibleKey uint32
	trolley     *trolley // purse 所属 trolley
	indexes     []uint8
	box         map[uint8]*box
}

func (p *purse) put(originalKey Key, key uint32, value interface{}) error {
	realKey := uint8(p.trolley.flexibleKey / boxDistance)
	p.flexibleKey = p.trolley.flexibleKey - uint32(realKey)*boxDistance
	//log.Self.Debug("purse", log.Uint32("key", key), log.Uint32("realKey", realKey))
	p.createChild(realKey)
	return p.box[realKey].put(originalKey, key, value)
}

func (p *purse) get(originalKey Key, key uint32) (interface{}, error) {
	realKey := uint8(p.trolley.flexibleKey / boxDistance)
	p.flexibleKey = p.trolley.flexibleKey - uint32(realKey)*boxDistance
	if p.existChild(realKey) {
		return p.box[realKey].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"purse key", string(originalKey), "is nil"}, " "))
	}
}

func (p *purse) existChild(index uint8) bool {
	return matchable(index, p.indexes)
}

func (p *purse) createChild(index uint8) {
	if !p.existChild(index) {
		p.indexes = append(p.indexes, index)
		p.box[index] = &box{
			key:    uint32(p.trolley.mall.city.realKey)*mallDistance + uint32(p.trolley.mall.realKey)*trolleyDistance + uint32(p.trolley.realKey)*purseDistance + uint32(index+1)*boxDistance,
			purse:  p,
			things: map[uint32]*thing{},
		}
	}
}
