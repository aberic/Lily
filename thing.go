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
	"github.com/ennoo/rivet/utils/log"
	"path/filepath"
	"strconv"
	"strings"
)

type thing struct {
	nodal       nodal // box 所属 purse
	originalKey Key
	seek        int64 // value最终存储在文件中的位置
	value       interface{}
}

func (t *thing) put(originalKey Key, key uint32, value interface{}) error {
	indexPath, formPath, form := t.getPath()
	log.Self.Debug("box", log.Uint32("key", key), log.Reflect("value", value), log.String("indexPath", indexPath), log.String("pathForm", formPath))
	wr := make(chan *writeResult, 1)
	go store().appendForm(form, formPath, value, wr)
	go store().appendIndex(t.nodal.getPreNodal().getPreNodal(), indexPath, t.uint32toFullState(key), wr)
	return nil
}

func (t *thing) get(originalKey Key, key uint32) (interface{}, error) {
	if t.originalKey == originalKey {
		return t.value, nil
	}
	return nil, errors.New(strings.Join([]string{"had no value for key ", string(originalKey)}, ""))
}

func (t *thing) saveIndex(form Form, indexPath string, key uint32) error {
	defer form.unLock()
	form.lock()
	return nil
}

func (t *thing) saveForm(formPath string, originalKey Key, key uint32, value interface{}) error {
	defer t.nodal.getPreNodal().getPreNodal().unLock() // level 2
	t.nodal.getPreNodal().getPreNodal().lock()         // level 2
	t.originalKey = originalKey
	t.value = value
	return nil
}

// getFormPath 获取表存储文件路径
func (t *thing) getPath() (indexPath, formPath string, form Form) {
	spr := t.nodal.getPreNodal().getPreNodal().getPreNodal().getPreNodal().(*shopper)
	dataID := spr.database.getID()
	formID := spr.id
	return pathIndex(dataID, formID, t.nodal.getPreNodal().getPreNodal().getPreNodal().getDegreeIndex()),
		filepath.Join(dataDir,
			dataID,
			formID,
			strconv.Itoa(int(t.nodal.getPreNodal().getPreNodal().getPreNodal().getDegreeIndex())), // level 1
			//strconv.Itoa(int(t.nodal.getPreNodal().getDegreeIndex())),
			strings.Join([]string{
				t.uint8toFullState(t.nodal.getPreNodal().getPreNodal().getDegreeIndex()), // level 2
				t.uint8toFullState(t.nodal.getPreNodal().getDegreeIndex()),               // level 3
				t.uint8toFullState(t.nodal.getDegreeIndex()),                             // level 4
				".dat"}, "",
			),
		),
		spr
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
