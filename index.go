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
	primary      bool    // 是否主键
	keyStructure string  // keyStructure 按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	form         Form    // form 索引所属表对象
	fileIndex    int     // 数据文件存储编号
	nodes        []Nodal // 节点
	fLock        sync.RWMutex
}

// getID 索引唯一ID
func (i *index) getID() string {
	return i.id
}

// isPrimary 是否主键
func (i *index) isPrimary() bool {
	return i.primary
}

// getKey 索引字段名称，由对象结构层级字段通过'.'组成，如
func (i *index) getKeyStructure() string {
	return i.keyStructure
}

// getForm 索引所属表对象
func (i *index) getForm() Form {
	return i.form
}

func (i *index) put(originalKey string, key uint32, update bool) IndexBack {
	index := key / cityDistance
	//index := uint32(0)
	node := i.createNode(uint8(index))
	return node.put(originalKey, key-index*cityDistance, 0, update)
}

func (i *index) get(originalKey string, key uint32) (interface{}, error) {
	index := key / cityDistance
	//index := uint32(0)
	if realIndex, err := i.existNode(uint8(index)); nil == err {
		return i.nodes[realIndex].get(originalKey, key-index*cityDistance, 0)
	}
	return nil, errors.New(strings.Join([]string{"index originalKey =", originalKey, "and keyStructure =", strconv.Itoa(int(key)), ", index =", strconv.Itoa(int(index)), "is nil"}, " "))
}

func (i *index) existNode(index uint8) (realIndex int, err error) {
	return binaryMatchData(uint8(index), i)
}

func (i *index) createNode(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	if realIndex, err = i.existNode(uint8(index)); nil != err {
		nd := &node{
			level:       0,
			degreeIndex: index,
			index:       i,
			preNode:     nil,
			nodes:       []Nodal{},
		}
		lenData := len(i.nodes)
		if lenData == 0 {
			i.nodes = append(i.nodes, nd)
			return nd
		}
		i.nodes = append(i.nodes, nd)
		for j := lenData - 1; j >= 0; j-- {
			if i.nodes[j].getDegreeIndex() < index {
				break
			} else if i.nodes[j].getDegreeIndex() > index {
				i.nodes[j+1] = i.nodes[j]
				i.nodes[j] = nd
				continue
			}
		}
		return nd
	}
	return i.nodes[realIndex]
}

func (i *index) getNodes() []Nodal {
	return i.nodes
}

func (i *index) lock() {
	i.fLock.Lock()
}

func (i *index) unLock() {
	i.fLock.Unlock()
}

func (i *index) rLock() {
	i.fLock.RLock()
}

func (i *index) rUnLock() {
	i.fLock.RUnlock()
}
