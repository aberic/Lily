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

// lily The Shopper
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
type lily struct {
	data    *Data  // 数据库对象
	name    string // 表明
	comment string // 描述
	purses  []*purse
}

func (l *lily) put(originalKey Key, key uint32, value interface{}) error {
	index := key / cityDistance
	//index := uint32(0)
	data := l.createChild(uint8(index))
	return data.put(originalKey, key-index*cityDistance, value)
}

func (l *lily) get(originalKey Key, key uint32) (interface{}, error) {
	index := key / cityDistance
	//index := uint32(0)
	if realIndex, err := binaryMatchData(uint8(index), l); nil == err {
		return l.purses[realIndex].get(originalKey, key-index*cityDistance)
	} else {
		return nil, errors.New(strings.Join([]string{"lily key", string(originalKey), "is nil"}, " "))
	}
}

func (l *lily) existChild(index uint8) bool {
	return matchableData(index, l)
}

func (l *lily) createChild(index uint8) nodal {
	if realIndex, err := binaryMatchData(index, l); nil != err {
		nd := &purse{
			level:       0,
			degreeIndex: index,
			nodal:       l,
			nodes:       []nodal{},
		}
		lenData := len(l.purses)
		if lenData == 0 {
			l.purses = append(l.purses, nd)
			return nd
		}
		l.purses = append(l.purses, nil)
		for i := len(l.purses) - 2; i >= 0; i-- {
			if l.purses[i].getDegreeIndex() < index {
				l.purses[i+1] = nd
				break
			} else if l.purses[i].getDegreeIndex() > index {
				l.purses[i+1] = l.purses[i]
				l.purses[i] = nd
			} else {
				return l.purses[i]
			}
		}
		return nd
	} else {
		return l.purses[realIndex]
	}
}

func (l *lily) childCount() int {
	return len(l.purses)
}

func (l *lily) child(index int) nodal {
	return l.purses[index]
}

func (l *lily) getDegreeIndex() uint8 {
	return 0
}

func (l *lily) getFlexibleKey() uint32 {
	return 0
}

func (l *lily) getPreNodal() nodal {
	return nil
}

func newLily(name, comment string, data *Data) *lily {
	lily := &lily{
		name:    name,
		comment: comment,
		data:    data,
		purses:  []*purse{},
	}
	return lily
}
