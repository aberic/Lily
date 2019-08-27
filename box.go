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

// box 包装盒一类
//
// 这里面存放的是具体物品
type box struct {
	degreeIndex uint8 // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	nodal       nodal // box 所属 purse
	things      []*thing
}

func (b *box) getFlexibleKey() uint32 {
	return 0
}

func (b *box) put(originalKey Key, key uint32, value interface{}) error {
	b.createChildSelf(originalKey, key, value)
	//log.Self.Debug("box", log.Uint32("key", key), log.Reflect("value", value))
	return b.things[len(b.things)-1].put(originalKey, key, value)
}

func (b *box) get(originalKey Key, key uint32) (interface{}, error) {
	if realIndex, exist := b.existChildSelf(originalKey, key); exist {
		return b.things[realIndex].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"box key", string(originalKey), "is nil"}, " "))
	}
}

func (b *box) existChild(index uint8) bool {
	return false
}

func (b *box) createChild(index uint8) nodal {
	return nil
}

func (b *box) existChildSelf(originalKey Key, key uint32) (int, bool) {
	for index, thg := range b.things {
		if strings.EqualFold(string(thg.originalKey), string(originalKey)) {
			return index, true
		}
	}
	return 0, false
}

func (b *box) createChildSelf(originalKey Key, key uint32, value interface{}) {
	if len(b.things) > 0 {
		for _, thg := range b.things {
			if strings.EqualFold(string(thg.originalKey), string(originalKey)) {
				return
			}
		}
	}
	b.things = append(b.things, &thing{nodal: b})
}

func (b *box) childCount() int {
	return len(b.things)
}

func (b *box) child(index int) nodeIndex {
	return nil
}

func (b *box) getDegreeIndex() uint8 {
	return b.degreeIndex
}

func (b *box) getPreNodal() nodal {
	return b.nodal
}
