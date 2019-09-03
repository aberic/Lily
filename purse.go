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
	"strings"
	"sync"
)

// purse 手提袋
//
// 这里面能存放很多个包装盒
//
// box 包装盒集合
//
// 存储格式 {dataDir}/checkbook/{dataName}/{formName}.shopper...
type purse struct {
	level       uint8  // 当前节点所在树层级
	degreeIndex uint8  // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	flexibleKey uint32 // 下一级最左最小树所对应真实key
	nodal       Nodal  // purse 所属 trolley
	nodes       []Nodal
	pLock       sync.RWMutex
}

func (p *purse) getFlexibleKey() uint32 {
	return p.flexibleKey
}

func (p *purse) put(indexID string, originalKey string, key uint32, value interface{}) *indexBack {
	var index uint8
	if p.level == 0 {
		index = uint8(key / mallDistance)
		p.flexibleKey = key - uint32(index)*mallDistance
	} else {
		index = uint8(p.nodal.getFlexibleKey() / distance(p.level))
		p.flexibleKey = p.nodal.getFlexibleKey() - uint32(index)*distance(p.level)
	}
	//if p.level == 2 {
	//
	//}
	//log.Self.Debug("purse", log.Uint32("key", key), log.Uint32("index", index))
	data := p.createChild(uint8(index))
	return data.put(indexID, originalKey, key, value)
}

func (p *purse) get(originalKey string, key uint32) (interface{}, error) {
	var index uint8
	if p.level == 0 {
		index = uint8(key / mallDistance)
		p.flexibleKey = key - uint32(index)*mallDistance
	} else {
		index = uint8(p.nodal.getFlexibleKey() / distance(p.level))
		p.flexibleKey = p.nodal.getFlexibleKey() - uint32(index)*distance(p.level)
	}
	if realIndex, err := binaryMatchData(uint8(index), p); nil == err {
		return p.nodes[realIndex].get(originalKey, key)
	}
	return nil, errors.New(strings.Join([]string{"purse key", originalKey, "is nil"}, " "))
}

func (p *purse) existChild(index uint8) bool {
	return matchableData(index, p)
}

func (p *purse) createChild(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	defer p.pLock.Unlock()
	p.pLock.Lock()
	if realIndex, err = binaryMatchData(index, p); nil != err {
		level := p.level + 1
		if level < levelMax {
			n := &purse{
				level:       level,
				degreeIndex: index,
				nodal:       p,
				nodes:       []Nodal{},
			}
			return p.appendNodal(index, n)
		}
		return p.appendNodal(index, &box{
			degreeIndex: index,
			nodal:       p,
			things:      []*thing{},
		})
	}
	return p.nodes[realIndex]
}

func (p *purse) appendNodal(index uint8, n Nodal) Nodal {
	lenData := len(p.nodes)
	if lenData == 0 {
		p.nodes = append(p.nodes, n)
		return n
	}
	p.nodes = append(p.nodes, nil)
	for i := len(p.nodes) - 2; i >= 0; i-- {
		if p.nodes[i].getDegreeIndex() < index {
			p.nodes[i+1] = n
			break
		} else if p.nodes[i].getDegreeIndex() > index {
			p.nodes[i+1] = p.nodes[i]
			p.nodes[i] = n
		} else {
			return p.nodes[i]
		}
	}
	return n
}

func (p *purse) childCount() int {
	return len(p.nodes)
}

func (p *purse) child(index int) Nodal {
	return p.nodes[index]
}

func (p *purse) getDegreeIndex() uint8 {
	return p.degreeIndex
}

func (p *purse) getPreNodal() Nodal {
	return p.nodal
}

func (p *purse) lock() {
	p.pLock.Lock()
}

func (p *purse) unLock() {
	p.pLock.Unlock()
}

func (p *purse) rLock() {
	p.pLock.RLock()
}

func (p *purse) rUnLock() {
	p.pLock.RUnlock()
}
