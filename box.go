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

package lily

import (
	"errors"
	"github.com/aberic/gnomon"
	"strings"
)

// box 包装盒一类
//
// 这里面存放的是具体物品
type box struct {
	degreeIndex uint8 // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	nodal       Nodal // box 所属 purse
	things      []*thing
}

func (b *box) getFlexibleKey() uint32 {
	return 0
}

func (b *box) put(originalKey string, key uint32, value interface{}, update bool) IndexBack {
	thg, exist := b.createChildSelf(originalKey, key, value)
	if !update && exist {
		return &indexBack{err: ErrDataExist}
	}
	//log.Self.Debug("box", log.Uint32("keyStructure", keyStructure), log.Reflect("value", value))
	return thg.put(originalKey, key, value)
}

func (b *box) get(originalKey string, key uint32) (interface{}, error) {
	gnomon.Log().Debug("box-get", gnomon.LogField("originalKey", originalKey))
	if realIndex, exist := b.existChildSelf(originalKey, key); exist {
		return b.things[realIndex].get()
	}
	return nil, errors.New(strings.Join([]string{"box keyStructure", originalKey, "is nil"}, " "))
}

func (b *box) existChild(index uint8) bool {
	return false
}

func (b *box) createChild(index uint8) Nodal {
	return nil
}

func (b *box) existChildSelf(originalKey string, key uint32) (int, bool) {
	for index, thg := range b.things {
		gnomon.Log().Debug("existChildSelf", gnomon.LogField("thg.md5Key", thg.md5Key), gnomon.LogField("md516", gnomon.CryptoHash().MD516(originalKey)))
		if strings.EqualFold(thg.md5Key, gnomon.CryptoHash().MD516(originalKey)) {
			return index, true
		}
	}
	return 0, false
}

func (b *box) createChildSelf(originalKey string, key uint32, value interface{}) (*thing, bool) {
	if len(b.things) > 0 {
		for _, thg := range b.things {
			if strings.EqualFold(thg.md5Key, gnomon.CryptoHash().MD516(originalKey)) {
				return thg, true
			}
		}
	}
	thg := &thing{nodal: b}
	b.things = append(b.things, thg)
	return thg, false
}

func (b *box) childCount() int {
	return -1
}

func (b *box) child(index int) Nodal {
	return nil
}

func (b *box) getDegreeIndex() uint8 {
	return b.degreeIndex
}

func (b *box) getPreNodal() Nodal {
	return b.nodal
}

func (b *box) lock() {

}

func (b *box) unLock() {

}

func (b *box) rLock() {

}

func (b *box) rUnLock() {

}
