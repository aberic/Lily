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
	preNode     Nodal // node 所属 trolley
	nodes       []Nodal
	links       []Link
	pLock       sync.RWMutex
}

func (p *node) getIndex() Index {
	return p.index
}

func (p *node) put(key string, hashKey, flexibleKey uint32, value interface{}, update bool) IndexBack {
	var (
		index          uint8
		flexibleNewKey uint32
		data           Nodal
	)
	if p.level == 0 {
		index = uint8(hashKey / mallDistance)
		flexibleNewKey = hashKey - uint32(index)*mallDistance
		data = p.createNode(uint8(index))
	} else if p.level == 4 {
		link, exist := p.createLink(key)
		if !update && exist {
			return &indexBack{err: p.errDataExist(key)}
		}
		formIndexFilePath := p.getFormIndexFilePath()
		if exist {
			return &indexBack{
				formIndexFilePath: formIndexFilePath,
				locker:            p.index,
				link:              link,
				key:               key,
				hashKey:           hashKey,
				err:               nil,
			}
		}
		//log.Self.Debug("box", log.Uint32("keyStructure", keyStructure), log.Reflect("value", value))
		return link.put(key, hashKey, value, formIndexFilePath)
	} else {
		index = uint8(flexibleKey / distance(p.level))
		flexibleNewKey = flexibleKey - uint32(index)*distance(p.level)
		if p.level == 3 {
			data = p.createLeaf(uint8(index))
		} else {
			data = p.createNode(uint8(index))
		}
	}
	return data.put(key, hashKey, flexibleNewKey, value, update)
}

func (p *node) get(key string, hashKey, flexibleKey uint32) (interface{}, error) {
	var (
		index          uint8
		flexibleNewKey uint32
	)
	if p.level == 0 {
		index = uint8(hashKey / mallDistance)
		flexibleNewKey = hashKey - uint32(index)*mallDistance
	} else if p.level == 4 {
		//gnomon.Log().Debug("box-get", gnomon.Log().Field("key", key))
		if realIndex, exist := p.existLink(key, hashKey); exist {
			return p.links[realIndex].get()
		}
		return nil, errors.New(strings.Join([]string{"box key", key, "is nil"}, " "))
	} else {
		index = uint8(flexibleKey / distance(p.level))
		flexibleNewKey = flexibleKey - uint32(index)*distance(p.level)
	}
	if realIndex, err := p.existNode(uint8(index)); nil == err {
		return p.nodes[realIndex].get(key, hashKey, flexibleNewKey)
	}
	return nil, errors.New(strings.Join([]string{"node key", key, "is nil"}, " "))
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
			preNode:     p,
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
			preNode:     p,
			links:       []Link{},
		}
		return p.appendNodal(index, n)
	}
	return p.nodes[realIndex]
}

func (p *node) createLink(key string) (Link, bool) {
	defer p.unLock()
	p.lock()
	if len(p.links) > 0 {
		for _, link := range p.links {
			//gnomon.Log().Debug("createLink", gnomon.Log().Field("exist", true))
			if strings.EqualFold(link.getMD5Key(), gnomon.CryptoHash().MD516(key)) {
				return link, true
			}
		}
	}
	link := &link{preNode: p, seekStartIndex: -1}
	p.links = append(p.links, link)
	return link, false
}

func (p *node) existLink(key string, hashKey uint32) (int, bool) {
	for index, link := range p.links {
		//gnomon.Log().Debug("existLink", gnomon.Log().Field("link.md5Key", link.getMD5Key()), gnomon.Log().Field("md516", gnomon.CryptoHash().MD516(key)))
		if strings.EqualFold(link.getMD5Key(), gnomon.CryptoHash().MD516(key)) {
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
	return pathFormIndexFile(dataID, formID, index.getID(), index.getKeyStructure())
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

func (p *node) getPreNode() Nodal {
	return p.preNode
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

// // errDataExist 自定义error信息
func (p *node) errDataExist(key string) error {
	return errors.New(strings.Join([]string{"data ", key, " already exist"}, ""))
}
