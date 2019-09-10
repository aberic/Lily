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
	"sync"
)

// node 手提袋
//
// 这里面能存放很多个包装盒
//
// box 包装盒集合
//
// 存储格式 {dataDir}/database/{dataName}/{formName}.form...
type node struct {
	level       uint8 // 当前节点所在树层级
	degreeIndex uint8 // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	index       Index // 所属索引对象
	nodal       Nodal // node 所属 trolley
	nodes       []Nodal
	links       []Link
	pLock       sync.RWMutex
}

func (p *node) getIndex() Index {
	return p.index
}

func (p *node) put(originalKey string, key, flexibleKey uint32, value interface{}, update bool) IndexBack {
	var (
		index          uint8
		flexibleNewKey uint32
		data           Nodal
	)
	if p.level == 0 {
		index = uint8(key / mallDistance)
		flexibleNewKey = key - uint32(index)*mallDistance
		data = p.createNode(uint8(index))
	} else if p.level == 4 {
		link, exist := p.createLink(originalKey, key, value)
		if !update && exist {
			return &indexBack{err: ErrDataExist}
		}
		formIndexFilePath := p.getFormIndexFilePath()
		if exist {
			return &indexBack{
				formIndexFilePath: formIndexFilePath,
				indexNodal:        p.nodal.getPreNodal(),
				link:              link,
				key:               key,
				err:               nil,
			}
		}
		//log.Self.Debug("box", log.Uint32("keyStructure", keyStructure), log.Reflect("value", value))
		return link.put(originalKey, key, value, formIndexFilePath)
	} else {
		index = uint8(flexibleKey / distance(p.level))
		flexibleNewKey = flexibleKey - uint32(index)*distance(p.level)
		if p.level == 3 {
			data = p.createLeaf(uint8(index))
		} else {
			data = p.createNode(uint8(index))
		}
	}
	return data.put(originalKey, key, flexibleNewKey, value, update)
}

func (p *node) get(originalKey string, key, flexibleKey uint32) (interface{}, error) {
	var (
		index          uint8
		flexibleNewKey uint32
	)
	if p.level == 0 {
		index = uint8(key / mallDistance)
		flexibleNewKey = key - uint32(index)*mallDistance
	} else if p.level == 4 {
		gnomon.Log().Debug("box-get", gnomon.LogField("originalKey", originalKey))
		if realIndex, exist := p.existLink(originalKey, key); exist {
			return p.links[realIndex].get()
		}
		return nil, errors.New(strings.Join([]string{"box keyStructure", originalKey, "is nil"}, " "))
	} else {
		index = uint8(flexibleKey / distance(p.level))
		flexibleNewKey = flexibleKey - uint32(index)*distance(p.level)
	}
	if realIndex, err := p.existNode(uint8(index)); nil == err {
		return p.nodes[realIndex].get(originalKey, key, flexibleNewKey)
	}
	return nil, errors.New(strings.Join([]string{"node keyStructure", originalKey, "is nil"}, " "))
}

func (p *node) existNode(index uint8) (realIndex int, err error) {
	return binaryMatchData(uint8(index), p)
}

func (p *node) createNode(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	defer p.unLock()
	p.lock()
	if realIndex, err = p.existNode(uint8(index)); nil != err {
		level := p.level + 1
		n := &node{
			level:       level,
			degreeIndex: index,
			index:       p.index,
			nodal:       p,
			nodes:       []Nodal{},
		}
		return p.appendNodal(index, n)
	}
	return p.nodes[realIndex]
}

func (p *node) createLeaf(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	defer p.unLock()
	p.lock()
	if realIndex, err = binaryMatchData(index, p); nil != err {
		level := p.level + 1
		n := &node{
			level:       level,
			degreeIndex: index,
			index:       p.index,
			nodal:       p,
			links:       []Link{},
		}
		return p.appendNodal(index, n)
	}
	return p.nodes[realIndex]
}

func (p *node) createLink(originalKey string, key uint32, value interface{}) (Link, bool) {
	defer p.unLock()
	p.lock()
	if len(p.links) > 0 {
		for _, thg := range p.links {
			if strings.EqualFold(thg.getMD5Key(), gnomon.CryptoHash().MD516(originalKey)) {
				return thg, true
			}
		}
	}
	thg := &link{nodal: p, seekStartIndex: -1}
	p.links = append(p.links, thg)
	return thg, false
}

func (p *node) existLink(originalKey string, key uint32) (int, bool) {
	for index, thg := range p.links {
		gnomon.Log().Debug("existChildSelf", gnomon.LogField("link.md5Key", thg.getMD5Key()), gnomon.LogField("md516", gnomon.CryptoHash().MD516(originalKey)))
		if strings.EqualFold(thg.getMD5Key(), gnomon.CryptoHash().MD516(originalKey)) {
			return index, true
		}
	}
	return 0, false
}

// getFormIndexFilePath 获取表索引文件路径
func (p *node) getFormIndexFilePath() (formIndexFilePath string) {
	index := p.getIndex()
	dataID := index.getForm().getDatabase().getID()
	formID := index.getForm().getID()
	rootNodeDegreeIndex := p.nodal.getPreNodal().getPreNodal().getDegreeIndex()
	return pathFormIndexFile(dataID, formID, index.getID(), rootNodeDegreeIndex)
}

func (p *node) appendNodal(index uint8, n Nodal) Nodal {
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

func (p *node) getNodes() []Nodal {
	return p.nodes
}

func (p *node) getLinks() []Link {
	return p.links
}

func (p *node) getDegreeIndex() uint8 {
	return p.degreeIndex
}

func (p *node) getPreNodal() Nodal {
	return p.nodal
}

func (p *node) lock() {
	p.pLock.Lock()
}

func (p *node) unLock() {
	p.pLock.Unlock()
}

func (p *node) rLock() {
	p.pLock.RLock()
}

func (p *node) rUnLock() {
	p.pLock.RUnlock()
}
