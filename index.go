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
	"strconv"
	"strings"
	"sync"
)

const (
	Int = iota
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	String
)

// index 索引对象
//
// 5位key及16位md5后key及5位起始seek和4位持续seek
type index struct {
	id           string  // id 索引唯一ID
	keyStructure string  // keyStructure 按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	form         Form    // form 索引所属表对象
	fileIndex    int     // 数据文件存储编号
	nodes        []Nodal // 节点
	fLock        sync.RWMutex
}

// getID 索引唯一ID
func (c *index) getID() string {
	return c.id
}

// getKey 索引字段名称，由对象结构层级字段通过'.'组成，如
func (c *index) getKeyStructure() string {
	return c.keyStructure
}

// getForm 索引所属表对象
func (c *index) getForm() Form {
	return c.form
}

func (c *index) put(originalKey string, key uint32, value interface{}, update bool) IndexBack {
	index := key / cityDistance
	//index := uint32(0)
	node := c.createNode(uint8(index))
	return node.put(originalKey, key-index*cityDistance, 0, value, update)
}

func (c *index) get(originalKey string, key uint32) (interface{}, error) {
	index := key / cityDistance
	//index := uint32(0)
	if realIndex, err := c.existNode(uint8(index)); nil == err {
		return c.nodes[realIndex].get(originalKey, key-index*cityDistance, 0)
	}
	return nil, errors.New(strings.Join([]string{"index originalKey =", originalKey, "and keyStructure =", strconv.Itoa(int(key)), ", index =", strconv.Itoa(int(index)), "is nil"}, " "))
}

func (c *index) existNode(index uint8) (realIndex int, err error) {
	return binaryMatchData(uint8(index), c)
}

func (c *index) createNode(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	if realIndex, err = c.existNode(uint8(index)); nil != err {
		nd := &node{
			level:       0,
			degreeIndex: index,
			index:       c,
			preNode:     nil,
			nodes:       []Nodal{},
		}
		lenData := len(c.nodes)
		if lenData == 0 {
			c.nodes = append(c.nodes, nd)
			return nd
		}
		c.nodes = append(c.nodes, nd)
		for i := lenData - 1; i >= 0; i-- {
			if c.nodes[i].getDegreeIndex() < index {
				break
			} else if c.nodes[i].getDegreeIndex() > index {
				c.nodes[i+1] = c.nodes[i]
				c.nodes[i] = nd
				continue
			}
		}
		return nd
	}
	return c.nodes[realIndex]
}

func (c *index) getNodes() []Nodal {
	return c.nodes
}

func (c *index) lock() {
	c.fLock.Lock()
}

func (c *index) unLock() {
	c.fLock.Unlock()
}

func (c *index) rLock() {
	c.fLock.RLock()
}

func (c *index) rUnLock() {
	c.fLock.RUnlock()
}
