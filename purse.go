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
	index       uint8    // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	flexibleKey uint32   // 下一级最左最小树所对应真实key
	trolley     *trolley // purse 所属 trolley
	box         []*box
}

func (p *purse) put(originalKey Key, key uint32, value interface{}) error {
	index := uint8(p.trolley.flexibleKey / boxDistance)
	p.flexibleKey = p.trolley.flexibleKey - uint32(index)*boxDistance
	//log.Self.Debug("purse", log.Uint32("key", key), log.Uint32("index", index))
	data := p.createChild(uint8(index))
	return data.put(originalKey, key, value)
}

func (p *purse) get(originalKey Key, key uint32) (interface{}, error) {
	index := uint8(p.trolley.flexibleKey / boxDistance)
	p.flexibleKey = p.trolley.flexibleKey - uint32(index)*boxDistance
	if realIndex, err := binaryMatchData(uint8(index), p); nil == err {
		return p.box[realIndex].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"purse key", string(originalKey), "is nil"}, " "))
	}
}

func (p *purse) existChild(index uint8) bool {
	return matchableData(index, p)
}

func (p *purse) createChild(index uint8) database {
	if realIndex, err := binaryMatchData(index, p); nil != err {
		b := &box{
			index:  index,
			purse:  p,
			things: map[uint32]*thing{},
		}
		lenCity := len(p.box)
		if lenCity == 0 {
			p.box = append(p.box, b)
			return b
		}
		p.box = append(p.box, nil)
		for i := len(p.box) - 2; i >= 0; i-- {
			if p.box[i].index < index {
				p.box[i+1] = b
				break
			} else if p.box[i].index > index {
				p.box[i+1] = p.box[i]
				p.box[i] = b
			} else {
				return p.box[i]
			}
		}
		return b
	} else {
		return p.box[realIndex]
	}
}

func (p *purse) childCount() int {
	return len(p.box)
}

func (p *purse) child(index int) nodeIndex {
	return p.box[index]
}

func (p *purse) getIndex() uint8 {
	return p.index
}
