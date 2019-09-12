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
	"sync"
)

type link struct {
	preNode        Nodal // box 所属 node
	md5Key         string
	seekStartIndex int64  // 索引最终存储在文件中的起始位置
	seekStart      uint32 // value最终存储在文件中的起始位置
	seekLast       int    // value最终存储在文件中的持续长度
	value          interface{}
	tLock          sync.RWMutex
}

func (t *link) setMD5Key(md5Key string) {
	t.md5Key = md5Key
}

func (t *link) setSeekStartIndex(seek int64) {
	t.seekStartIndex = seek
}

func (t *link) setSeekStart(seek uint32) {
	t.seekStart = seek
}

func (t *link) setSeekLast(seek int) {
	t.seekLast = seek
}

func (t *link) getNodal() Nodal {
	return t.preNode
}

func (t *link) getMD5Key() string {
	return t.md5Key
}

func (t *link) getSeekStartIndex() int64 {
	return t.seekStartIndex
}

func (t *link) getSeekStart() uint32 {
	return t.seekStart
}

func (t *link) getSeekLast() int {
	return t.seekLast
}

func (t *link) getValue() interface{} {
	return t.value
}

func (t *link) lock() {
	t.tLock.Lock()
}

func (t *link) unLock() {
	t.tLock.Unlock()
}

func (t *link) rLock() {
	t.tLock.RLock()
}

func (t *link) rUnLock() {
	t.tLock.RUnlock()
}

func (t *link) put(key string, hashKey uint32, value interface{}, formIndexFilePath string) *indexBack {
	gnomon.Log().Debug("box",
		gnomon.Log().Field("key", key),
		gnomon.Log().Field("hashKey", hashKey),
		gnomon.Log().Field("value", value),
		gnomon.Log().Field("formIndexFilePath", formIndexFilePath))
	return &indexBack{
		formIndexFilePath: formIndexFilePath,
		locker:            t.preNode.getIndex(),
		link:              t,
		key:               key,
		hashKey:           hashKey,
		err:               nil,
	}
}

func (t *link) get() (interface{}, error) {
	index := t.preNode.getIndex()
	rrFormBack := make(chan *readResult, 1)
	go store().read(pathFormDataFile(index.getForm().getDatabase().getID(), index.getForm().getID(), index.getForm().getFileIndex()), t.seekStart, t.seekLast, rrFormBack)
	rr := <-rrFormBack
	return rr.value, rr.err
}

// indexBack 索引对象
type indexBack struct {
	formIndexFilePath string      // 索引文件所在路径
	locker            WriteLocker // 索引文件所对应level2层级度节点
	link              Link        // 索引对应节点对象子集
	key               string      // 索引对应字符串key
	hashKey           uint32      // put hash hashKey
	err               error
}

// getFormIndexFilePath 索引文件所在路径
func (i *indexBack) getFormIndexFilePath() string {
	return i.formIndexFilePath
}

// getLocker 索引文件所对应level2层级度节点
func (i *indexBack) getLocker() WriteLocker {
	return i.locker
}

// getLink 索引对应节点对象子集
func (i *indexBack) getLink() Link {
	return i.link
}

// getKey 索引对应字符串key
func (i *indexBack) getKey() string {
	return i.key
}

// getHashKey put hash keyStructure
func (i *indexBack) getHashKey() uint32 {
	return i.hashKey
}

// getErr error信息
func (i *indexBack) getErr() error {
	return i.err
}
