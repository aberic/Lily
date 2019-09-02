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
	"github.com/ennoo/rivet/utils/log"
	"strconv"
	"strings"
)

type thing struct {
	nodal     nodal // box 所属 purse
	md5Key    string
	seekStart uint32 // value最终存储在文件中的起始位置
	seekLast  int    // value最终存储在文件中的持续长度
	value     interface{}
}

func (t *thing) put(indexID string, originalKey string, key uint32, value interface{}) *indexBack {
	formIndexFilePath := t.getFormIndexFilePath(indexID)
	log.Self.Debug("box", log.Uint32("key", key), log.Reflect("value", value), log.String("formIndexFilePath", formIndexFilePath))
	return &indexBack{
		formIndexFilePath: formIndexFilePath,
		indexNodal:        t.nodal.getPreNodal().getPreNodal(),
		thing:             t,
		originalKey:       originalKey,
		key:               key,
		err:               nil,
	}
}

func (t *thing) get() (interface{}, error) {
	spr := t.nodal.getPreNodal().getPreNodal().getPreNodal().getPreNodal().(*shopper)
	rrFormBack := make(chan *readResult, 1)
	if err := pool().submit(func() {
		store().read(pathFormDataFile(spr.database.getID(), spr.id, spr.fileIndex), t.seekStart, t.seekLast, rrFormBack)
	}); nil != err {
		return nil, err
	}
	rr := <-rrFormBack
	return rr.value, rr.err
}

// getFormIndexFilePath 获取表索引文件路径
func (t *thing) getFormIndexFilePath(indexID string) (formIndexFilePath string) {
	spr := t.nodal.getPreNodal().getPreNodal().getPreNodal().getPreNodal().(*shopper)
	dataID := spr.database.getID()
	formID := spr.id
	rootNodeDegreeIndex := t.nodal.getPreNodal().getPreNodal().getPreNodal().getDegreeIndex()
	return pathFormIndexFile(dataID, formID, indexID, rootNodeDegreeIndex)
}

// uint8toFullState 补全不满三位数状态，如1->001、34->034、215->215
func (t *thing) uint8toFullState(index uint8) string {
	result := strconv.Itoa(int(index))
	if index < 10 {
		return strings.Join([]string{"00", result}, "")
	} else if index < 100 {
		return strings.Join([]string{"0", result}, "")
	}
	return result
}

// uint32toFullState 补全不满十位数状态，如1->0000000001、34->0000000034、215->0000000215
func (t *thing) uint32toFullState(index uint32) string {
	pos := 0
	for index > 1 {
		index /= 10
		pos++
	}
	backZero := 10 - pos
	backZeroStr := strconv.Itoa(int(index))
	for i := 0; i < backZero; i++ {
		backZeroStr = strings.Join([]string{"0", backZeroStr}, "")
	}
	return backZeroStr
}
