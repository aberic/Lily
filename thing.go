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
	"github.com/aberic/gnomon"
)

type thing struct {
	nodal     Nodal // box 所属 purse
	md5Key    string
	seekStart uint32 // value最终存储在文件中的起始位置
	seekLast  int    // value最终存储在文件中的持续长度
	value     interface{}
}

func (t *thing) put(indexID string, originalKey string, key uint32, value interface{}) *indexBack {
	formIndexFilePath := t.getFormIndexFilePath(indexID)
	gnomon.Log().Debug("box",
		gnomon.LogField("originalKey", originalKey),
		gnomon.LogField("key", key),
		gnomon.LogField("value", value),
		gnomon.LogField("formIndexFilePath", formIndexFilePath))
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
