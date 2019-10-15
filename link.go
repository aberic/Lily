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
	"sync"
)

type link struct {
	preNode        Nodal // box 所属 node
	md516Key       string
	seekStartIndex int64 // 索引最终存储在文件中的起始位置
	seekStart      int64 // value最终存储在文件中的起始位置
	seekLast       int   // value最终存储在文件中的持续长度
	tLock          sync.RWMutex
}

func (l *link) setMD5Key(md5Key string) {
	l.md516Key = md5Key
}

func (l *link) setSeekStartIndex(seek int64) {
	l.seekStartIndex = seek
}

func (l *link) setSeekStart(seek int64) {
	l.seekStart = seek
}

func (l *link) setSeekLast(seek int) {
	l.seekLast = seek
}

func (l *link) getNodal() Nodal {
	return l.preNode
}

func (l *link) getMD516Key() string {
	return l.md516Key
}

func (l *link) getSeekStartIndex() int64 {
	return l.seekStartIndex
}

func (l *link) getSeekStart() int64 {
	return l.seekStart
}

func (l *link) getSeekLast() int {
	return l.seekLast
}

func (l *link) lock() {
	l.tLock.Lock()
}

func (l *link) unLock() {
	l.tLock.Unlock()
}

func (l *link) rLock() {
	l.tLock.RLock()
}

func (l *link) rUnLock() {
	l.tLock.RUnlock()
}

func (l *link) put(key string, hashKey uint64) *indexBack {
	formIndexFilePath := l.getFormIndexFilePath()
	//gnomon.Log().Debug("box",
	//	gnomon.Log().Field("key", key),
	//	gnomon.Log().Field("hashKey", hashKey),
	//	gnomon.Log().Field("formIndexFilePath", formIndexFilePath))
	return &indexBack{
		formIndexFilePath: formIndexFilePath,
		locker:            l.preNode.getIndex(),
		link:              l,
		key:               key,
		hashKey:           hashKey,
		err:               nil,
	}
}

func (l *link) get() *readResult {
	index := l.preNode.getIndex()
	return store().read(pathFormDataFile(index.getForm().getDatabase().getID(), index.getForm().getID()), l.seekStart, l.seekLast)
}

// getFormIndexFilePath 获取表索引文件路径
func (l *link) getFormIndexFilePath() (formIndexFilePath string) {
	index := l.preNode.getIndex()
	dataID := index.getForm().getDatabase().getID()
	formID := index.getForm().getID()
	return pathFormIndexFile(dataID, formID, index.getID())
}

// indexBack 索引对象
type indexBack struct {
	formIndexFilePath string      // 索引文件所在路径
	locker            WriteLocker // 索引文件所对应level2层级度节点
	link              Link        // 索引对应节点对象子集
	key               string      // 索引对应字符串key
	hashKey           uint64      // put hash hashKey
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
func (i *indexBack) getHashKey() uint64 {
	return i.hashKey
}

// getErr error信息
func (i *indexBack) getErr() error {
	return i.err
}
