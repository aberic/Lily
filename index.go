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
	"os"
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
	id           string // id 索引唯一ID
	primary      bool   // 是否主键
	keyStructure string // keyStructure 按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	form         Form   // form 索引所属表对象
	node         Nodal  // 节点
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

func (i *index) put(originalKey string, key int64, update bool) IndexBack {
	return i.node.put(originalKey, key, key, update)
	//index := key / cityDistance
	//node := i.createNode(uint8(index))
	//return node.put(originalKey, key-index*cityDistance, 0, update)
}

func (i *index) get(originalKey string, key int64) (interface{}, error) {
	return i.node.get(originalKey, key, key)
	//index := key / cityDistance
	//if realIndex, err := i.existNode(uint8(index)); nil == err {
	//	return i.nodes[realIndex].get(originalKey, key-index*cityDistance, 0)
	//}
	//return nil, errors.New(strings.Join([]string{"index originalKey =", originalKey, "and keyStructure =", strconv.Itoa(int(key)), ", index =", strconv.Itoa(int(index)), "is nil"}, " "))
}

func (i *index) recover() error {
	// todo 恢复索引，注意Form的autoID
	indexFilePath := pathFormIndexFile(i.form.getDatabase().getID(), i.form.getID(), i.id)
	if gnomon.File().PathExists(indexFilePath) { // 索引文件不存在，则无需操作
		var (
			file *os.File
			err  error
		)
		if file, err = os.OpenFile(indexFilePath, os.O_CREATE|os.O_RDWR, 0644); nil != err {
			gnomon.Log().Panic("index recover failed", gnomon.Log().Err(err))
		}
		_ = file.Close()
	}
	return nil
}

func (i *index) getNode() Nodal {
	return i.node
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
