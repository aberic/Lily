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

// index 索引对象
//
// 5位key及16位md5后key及5位起始seek和4位持续seek
type catalog struct {
	id        string  // id 索引唯一ID
	key       string  // 索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	fileIndex int     // 数据文件存储编号
	nodes     []Nodal // 节点
	fLock     sync.RWMutex
}

// id 索引唯一ID
func (c *catalog) getID() string {
	return c.id
}

// 索引字段名称，由对象结构层级字段通过'.'组成，如
func (c *catalog) getKey() string {
	return c.key
}

func (c *catalog) put(indexID string, originalKey string, key uint32, value interface{}, update bool) *indexBack {
	index := key / cityDistance
	//index := uint32(0)
	data := c.createChild(uint8(index))
	return data.put(indexID, originalKey, key-index*cityDistance, value, update)
}

func (c *catalog) get(originalKey string, key uint32) (interface{}, error) {
	index := key / cityDistance
	//index := uint32(0)
	if realIndex, err := binaryMatchData(uint8(index), c); nil == err {
		return c.nodes[realIndex].get(originalKey, key-index*cityDistance)
	}
	return nil, errors.New(strings.Join([]string{"shopper originalKey =", originalKey, "and key =", strconv.Itoa(int(key)), ", index =", strconv.Itoa(int(index)), "is nil"}, " "))
}

func (c *catalog) existChild(index uint8) bool {
	return matchableData(index, c)
}

func (c *catalog) createChild(index uint8) Nodal {
	var (
		realIndex int
		err       error
	)
	if realIndex, err = binaryMatchData(index, c); nil != err {
		nd := &purse{
			level:       0,
			degreeIndex: index,
			nodal:       c,
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

func (c *catalog) childCount() int {
	return len(c.nodes)
}

func (c *catalog) child(index int) Nodal {
	return c.nodes[index]
}

func (c *catalog) getDegreeIndex() uint8 {
	return 0
}

func (c *catalog) getFlexibleKey() uint32 {
	return 0
}

func (c *catalog) getPreNodal() Nodal {
	return nil
}

func (c *catalog) lock() {
	c.fLock.Lock()
}

func (c *catalog) unLock() {
	c.fLock.Unlock()
}

func (c *catalog) rLock() {
	c.fLock.RLock()
}

func (c *catalog) rUnLock() {
	c.fLock.RUnlock()
}
