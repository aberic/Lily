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
	"sync"
)

// shopper The Shopper
//
// hash array 模型 [00, 01, 02, 03, 04, 05, 06, 07, 08, 09, a, b, c, d, e, f]
//
// b+tree 模型 degree=128;level=4;nodes=[degree^level]/(degree-1)=2113665;
//
// purse 内范围控制数量 key=127
//
// tree 内范围控制数量 treeCount=nodes*key=268435455
//
// hash array 内范围控制数量 t*16=4294967280
//
// level1间隔 ld1=(treeCount+1)/128=2097152
//
// level2间隔 ld2=(16513*127+1)/128=16384
//
// level3间隔 ld3=(129*127+1)/128=128
//
// level4间隔 ld3=(1*127+1)/128=1
//
// 存储格式 {dataDir}/checkbook/{dataName}/{formName}/{formName}.dat/idx...
//
// 索引格式
type shopper struct {
	autoID    uint32   // 自增id
	database  Database // 数据库对象
	name      string   // 表名，根据需求可以随时变化
	id        string   // 表唯一ID，不能改变
	indexIDs  []string // 索引ID集合
	fileIndex int      // 数据文件存储编号
	comment   string   // 描述
	nodes     []Nodal  // 节点
	fLock     sync.RWMutex
}

func (s *shopper) getAutoID() *uint32 {
	return &s.autoID
}

func (s *shopper) getID() string {
	return s.id
}

func (s *shopper) getName() string {
	return s.name
}

func (s *shopper) getFileIndex() int {
	return s.fileIndex
}

func (s *shopper) getIndexIDs() []string {
	return s.indexIDs
}

func (s *shopper) getDatabaseID() string {
	return s.database.getID()
}

func (s *shopper) put(indexID string, originalKey string, key uint32, value interface{}) *indexBack {
	index := key / cityDistance
	//index := uint32(0)
	data := s.createChild(uint8(index))
	return data.put(indexID, originalKey, key-index*cityDistance, value)
}

func (s *shopper) get(originalKey string, key uint32) (interface{}, error) {
	index := key / cityDistance
	//index := uint32(0)
	if realIndex, err := binaryMatchData(uint8(index), s); nil == err {
		return s.nodes[realIndex].get(originalKey, key-index*cityDistance)
	} else {
		return nil, errors.New(strings.Join([]string{"shopper key", originalKey, "is nil"}, " "))
	}
}

func (s *shopper) existChild(index uint8) bool {
	return matchableData(index, s)
}

func (s *shopper) createChild(index uint8) Nodal {
	if realIndex, err := binaryMatchData(index, s); nil != err {
		nd := &purse{
			level:       0,
			degreeIndex: index,
			nodal:       s,
			nodes:       []Nodal{},
		}
		lenData := len(s.nodes)
		if lenData == 0 {
			s.nodes = append(s.nodes, nd)
			return nd
		}
		s.nodes = append(s.nodes, nil)
		for i := len(s.nodes) - 2; i >= 0; i-- {
			if s.nodes[i].getDegreeIndex() < index {
				s.nodes[i+1] = nd
				break
			} else if s.nodes[i].getDegreeIndex() > index {
				s.nodes[i+1] = s.nodes[i]
				s.nodes[i] = nd
			} else {
				return s.nodes[i]
			}
		}
		return nd
	} else {
		return s.nodes[realIndex]
	}
}

func (s *shopper) childCount() int {
	return len(s.nodes)
}

func (s *shopper) child(index int) Nodal {
	return s.nodes[index]
}

func (s *shopper) getDegreeIndex() uint8 {
	return 0
}

func (s *shopper) getFlexibleKey() uint32 {
	return 0
}

func (s *shopper) getPreNodal() Nodal {
	return nil
}

func (s *shopper) lock() {
	s.fLock.Lock()
}

func (s *shopper) unLock() {
	s.fLock.Unlock()
}

func (s *shopper) rLock() {
	s.fLock.RLock()
}

func (s *shopper) rUnLock() {
	s.fLock.RUnlock()
}
