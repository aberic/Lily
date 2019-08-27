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
// b+tree 模型 degree=128;level=4;nodes=[degree^level]/(degree-1)=2113665;
//
// node 内范围控制数量 key=127
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
type lily struct {
	data    *Data  // 数据库对象
	name    string // 表明
	comment string // 描述
	cities  []*city
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
		return l.cities[realIndex].get(originalKey, key-index*cityDistance)
	} else {
		return nil, errors.New(strings.Join([]string{"lily key", string(originalKey), "is nil"}, " "))
	}
}

func (l *lily) existChild(index uint8) bool {
	return matchableData(index, l)
}

func (l *lily) createChild(index uint8) database {
	if realIndex, err := binaryMatchData(index, l); nil != err {
		c := &city{
			index: index,
			lily:  l,
			malls: []*mall{},
		}
		lenCity := len(l.cities)
		if lenCity == 0 {
			l.cities = append(l.cities, c)
			return c
		}
		l.cities = append(l.cities, nil)
		for i := len(l.cities) - 2; i >= 0; i-- {
			if l.cities[i].index < index {
				l.cities[i+1] = c
				break
			} else if l.cities[i].index > index {
				l.cities[i+1] = l.cities[i]
				l.cities[i] = c
			} else {
				return l.cities[i]
			}
		}
		return c
	} else {
		return l.cities[realIndex]
	}
}

func (l *lily) childCount() int {
	return len(l.cities)
}

func (l *lily) child(index int) nodeIndex {
	return l.cities[index]
}

func (l *lily) getIndex() uint8 {
	return 0
}

func newLily(name, comment string, data *Data) *lily {
	lily := &lily{
		name:    name,
		comment: comment,
		data:    data,
		cities:  []*city{},
	}
	return lily
}
