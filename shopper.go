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

// shopper The Shopper
//
// hash array 模型 [00, 01, 02, 03, 04, 05, 06, 07, 08, 09, a, b, c, d, e, f]
//
// b+tree 模型 degree=128;level=4;purses=[degree^level]/(degree-1)=2113665;
//
// purse 内范围控制数量 key=127
//
// tree 内范围控制数量 treeCount=purses*key=268435455
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
// 存储格式 {dataDir}/checkbook/{dataName}/{shopperName}/{shopperName}.dat/idx...
//
// 索引格式
type shopper struct {
	autoID    uint32     // 自增id
	checkbook *checkbook // 数据库对象
	name      string     // 表名，根据需求可以随时变化
	id        string     // 表唯一ID，不能改变
	comment   string     // 描述
	purses    []*purse   // 节点
}

func (s *shopper) put(originalKey Key, key uint32, value interface{}) error {
	//index := key / cityDistance
	index := uint32(0)
	data := s.createChild(uint8(index))
	return data.put(originalKey, key-index*cityDistance, value)
}

func (s *shopper) get(originalKey Key, key uint32) (interface{}, error) {
	//index := key / cityDistance
	index := uint32(0)
	if realIndex, err := binaryMatchData(uint8(index), s); nil == err {
		return s.purses[realIndex].get(originalKey, key-index*cityDistance)
	} else {
		return nil, errors.New(strings.Join([]string{"shopper key", string(originalKey), "is nil"}, " "))
	}
}

func (s *shopper) existChild(index uint8) bool {
	return matchableData(index, s)
}

func (s *shopper) createChild(index uint8) nodal {
	if realIndex, err := binaryMatchData(index, s); nil != err {
		nd := &purse{
			level:       0,
			degreeIndex: index,
			nodal:       s,
			nodes:       []nodal{},
		}
		lenData := len(s.purses)
		if lenData == 0 {
			s.purses = append(s.purses, nd)
			return nd
		}
		s.purses = append(s.purses, nil)
		for i := len(s.purses) - 2; i >= 0; i-- {
			if s.purses[i].getDegreeIndex() < index {
				s.purses[i+1] = nd
				break
			} else if s.purses[i].getDegreeIndex() > index {
				s.purses[i+1] = s.purses[i]
				s.purses[i] = nd
			} else {
				return s.purses[i]
			}
		}
		return nd
	} else {
		return s.purses[realIndex]
	}
}

func (s *shopper) childCount() int {
	return len(s.purses)
}

func (s *shopper) child(index int) nodal {
	return s.purses[index]
}

func (s *shopper) getDegreeIndex() uint8 {
	return 0
}

func (s *shopper) getFlexibleKey() uint32 {
	return 0
}

func (s *shopper) getPreNodal() nodal {
	return nil
}

func newShopper(name, id, comment string, checkbook *checkbook) *shopper {
	lily := &shopper{
		autoID:    0,
		name:      name,
		id:        id,
		comment:   comment,
		checkbook: checkbook,
		purses:    []*purse{},
	}
	return lily
}

func (s *shopper) lock() {

}

func (s *shopper) unLock() {

}

func (s *shopper) rLock() {

}

func (s *shopper) rUnLock() {

}
