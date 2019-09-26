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

func (n *node) getIndex() Index {
	return n.index
}

func (n *node) put(key string, hashKey, flexibleKey int64, update bool) IndexBack {
	var (
		index          uint8
		flexibleNewKey int64
		data           Nodal
	)
	if n.level == 0 {
		index = uint8(hashKey / mallDistance)
		flexibleNewKey = hashKey - int64(index)*mallDistance
		data = n.createNode(uint8(index))
	} else if n.level == 4 {
		link, exist := n.createLink(key)
		if !update && exist {
			return &indexBack{err: n.errDataExist(key)}
		}
		//log.Self.Debug("box", log.Uint32("keyStructure", keyStructure), log.Reflect("value", value))
		return link.put(key, hashKey)
	} else {
		index = uint8(flexibleKey / distance(n.level))
		flexibleNewKey = flexibleKey - int64(index)*distance(n.level)
		if n.level == 3 {
			data = n.createLeaf(uint8(index))
		} else {
			data = n.createNode(uint8(index))
		}
	}
	return data.put(key, hashKey, flexibleNewKey, update)
}

func (n *node) get(key string, hashKey, flexibleKey int64) (interface{}, error) {
	var (
		index          uint8
		flexibleNewKey int64
	)
	if n.level == 0 {
		index = uint8(hashKey / mallDistance)
		flexibleNewKey = hashKey - int64(index)*mallDistance
	} else if n.level == 4 {
		//gnomon.Log().Debug("box-get", gnomon.Log().Field("key", key))
		if realIndex, exist := n.existLink(key); exist {
			return n.links[realIndex].get()
		}
		return nil, errors.New(strings.Join([]string{"box key", key, "is nil"}, " "))
	} else {
		index = uint8(flexibleKey / distance(n.level))
		flexibleNewKey = flexibleKey - int64(index)*distance(n.level)
	}
	if realIndex, err := n.existNode(uint8(index)); nil == err {
		return n.nodes[realIndex].get(key, hashKey, flexibleNewKey)
	}
	return nil, errors.New(strings.Join([]string{"node key", key, "is nil"}, " "))
}

func (n *node) existNode(index uint8) (realIndex int, err error) {
	return binaryMatchData(uint8(index), n)
}

func (n *node) createNode(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	defer n.unLock()
	n.lock()
	if realIndex, err = n.existNode(uint8(index)); nil != err {
		level := n.level + 1
		newNode := &node{
			level:       level,
			degreeIndex: index,
			index:       n.index,
			preNode:     n,
			nodes:       []Nodal{},
		}
		return n.appendNodal(index, newNode)
	}
	return n.nodes[realIndex]
}

func (n *node) createLeaf(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	defer n.unLock()
	n.lock()
	if realIndex, err = binaryMatchData(index, n); nil != err {
		level := n.level + 1
		leaf := &node{
			level:       level,
			degreeIndex: index,
			index:       n.index,
			preNode:     n,
			links:       []Link{},
		}
		return n.appendNodal(index, leaf)
	}
	return n.nodes[realIndex]
}

func (n *node) createLink(key string) (Link, bool) {
	defer n.unLock()
	n.lock()
	if n.index.isPrimary() && len(n.links) > 0 {
		for _, link := range n.links {
			//gnomon.Log().Debug("createLink", gnomon.Log().Field("exist", true))
			if strings.EqualFold(link.getMD5Key(), gnomon.CryptoHash().MD516(key)) {
				return link, true
			}
		}
	}
	link := &link{preNode: n, seekStartIndex: -1}
	n.links = append(n.links, link)
	return link, false
}

func (n *node) existLink(key string) (int, bool) {
	for index, link := range n.links {
		//gnomon.Log().Debug("existLink", gnomon.Log().Field("link.md5Key", link.getMD5Key()), gnomon.Log().Field("md516", gnomon.CryptoHash().MD516(key)))
		if strings.EqualFold(link.getMD5Key(), gnomon.CryptoHash().MD516(key)) {
			return index, true
		}
	}
	return 0, false
}

func (n *node) appendNodal(index uint8, nodal Nodal) Nodal {
	lenData := len(n.nodes)
	if lenData == 0 {
		n.nodes = append(n.nodes, nodal)
		return nodal
	}
	n.nodes = append(n.nodes, nil)
	for i := len(n.nodes) - 2; i >= 0; i-- {
		if n.nodes[i].getDegreeIndex() < index {
			n.nodes[i+1] = nodal
			break
		} else if n.nodes[i].getDegreeIndex() > index {
			n.nodes[i+1] = n.nodes[i]
			n.nodes[i] = nodal
		} else {
			return n.nodes[i]
		}
	}
	return nodal
}

func (n *node) getNodes() []Nodal {
	return n.nodes
}

func (n *node) getLinks() []Link {
	return n.links
}

func (n *node) getDegreeIndex() uint8 {
	return n.degreeIndex
}

func (n *node) getPreNode() Nodal {
	return n.preNode
}

func (n *node) lock() {
	n.pLock.Lock()
}

func (n *node) unLock() {
	n.pLock.Unlock()
}

func (n *node) rLock() {
	n.pLock.RLock()
}

func (n *node) rUnLock() {
	n.pLock.RUnlock()
}

// // errDataExist 自定义error信息
func (n *node) errDataExist(key string) error {
	return errors.New(strings.Join([]string{"data ", key, " already exist"}, ""))
}
