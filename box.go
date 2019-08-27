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
	index  uint8  // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	purse  *purse // box 所属 purse
	things map[uint32]*thing
}

func (b *box) put(originalKey Key, key uint32, value interface{}) error {
	b.createChildSelf(originalKey, key, value)
	//log.Self.Debug("box", log.Uint32("key", key), log.Reflect("value", value))
	return b.things[key].put(originalKey, key, value)
}

func (b *box) get(originalKey Key, key uint32) (interface{}, error) {
	if b.existChildSelf(key) {
		return b.things[key].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"box key", string(originalKey), "is nil"}, " "))
	}
}

func (b *box) existChild(key uint8) bool {
	return false
}

func (b *box) createChild(index uint8) database {
	return nil
}

func (b *box) existChildSelf(key uint32) bool {
	return nil != b.things[key]
}

func (b *box) createChildSelf(originalKey Key, key uint32, value interface{}) {
	if !b.existChildSelf(key) {
		b.things[key] = &thing{
			box:         b,
			originalKey: originalKey,
			value:       value,
		}
	}
}

func (b *box) getIndex() uint8 {
	return b.index
}
