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
	"bufio"
	"errors"
	"github.com/aberic/gnomon"
	"io"
	"os"
	"sync"
	"sync/atomic"
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

func (i *index) put(key string, hashKey uint64, update bool) IndexBack {
	return i.node.put(key, hashKey, hashKey, update)
}

func (i *index) get(key string, hashKey uint64) (interface{}, error) {
	return i.node.get(key, hashKey, hashKey)
}

func (i *index) recover() {
	i.recoverMultiReadFile()
}

func (i *index) recoverMultiReadFile() {
	indexFilePath := pathFormIndexFile(i.form.getDatabase().getID(), i.form.getID(), i.id)
	if gnomon.File().PathExists(indexFilePath) { // 索引文件存在才继续恢复
		var (
			file *os.File
			err  error
		)
		defer func() { _ = file.Close() }()
		if file, err = os.OpenFile(indexFilePath, os.O_RDONLY, 0644); nil != err {
			gnomon.Log().Panic("index recover multi read failed", gnomon.Log().Err(err))
		}
		_, err = file.Seek(0, io.SeekStart) // 文件下标置为文件的起始位置
		if err != nil {
			gnomon.Log().Panic("index recover multi read failed", gnomon.Log().Err(err))
		}
		if err = i.read(file, 0); nil != err && io.EOF != err {
			gnomon.Log().Panic("index recover multi read failed", gnomon.Log().Err(err))
		}
	}
}

func (i *index) read(file *os.File, offset int64) (err error) {
	var (
		inputReader *bufio.Reader
		data        []byte
		peekOnce          = 36000
		haveNext          = true
		position    int64 = 0
	)
	_, err = file.Seek(offset, io.SeekStart) //表示文件的起始位置，从第二个字符往后写入。
	inputReader = bufio.NewReaderSize(file, peekOnce)
	data, err = inputReader.Peek(peekOnce)
	if nil != err && io.EOF != err {
		return
	} else if nil != err && io.EOF == err {
		if len(data) == 0 {
			return
		}
		if len(data)%36 != 0 {
			return errors.New("index lens does't match")
		}
	}
	indexStr := string(data)
	indexStrLen := int64(len(indexStr))
	for haveNext {
		go func(i *index, position int64) {
			var p0, p1, p2, p3, p4 int64
			// 读取11位key及16位md5后key及5位起始seek和4位持续seek
			p0 = position
			p1 = p0 + 11
			p2 = p1 + 16
			p3 = p2 + 5
			p4 = p3 + 4
			hashKey := gnomon.Scale().DDuoStringToUint64(indexStr[p0:p1])
			md516Key := indexStr[p1:p2]
			seekStart := uint32(gnomon.Scale().DDuoStringToUint64(indexStr[p2:p3])) // value最终存储在文件中的起始位置
			seekLast := int(gnomon.Scale().DDuoStringToInt64(indexStr[p3:p4]))      // value最终存储在文件中的持续长度
			ib := i.node.put("", hashKey, hashKey, true)
			ib.getLink().setSeekStartIndex(p0)
			ib.getLink().setMD5Key(md516Key)
			ib.getLink().setSeekStart(seekStart)
			ib.getLink().setSeekLast(seekLast)
			atomic.AddUint64(i.form.getAutoID(), 1) // ID自增
		}(i, position)
		position += 36
		if indexStrLen < position+36 {
			haveNext = false
		}
	}
	if nil == err {
		offset += int64(peekOnce)
		return i.read(file, offset)
	}
	return
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
