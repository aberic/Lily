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

// shopper The Shopper
//
// hash array 模型 [00, 01, 02, 03, 04, 05, 06, 07, 08, 09, a, b, c, d, e, f]
//
// b+tree 模型 degree=128;level=4;nodes=[degree^level]/(degree-1)=2113665;
//
// purse 内范围控制数量 keyStructure=127
//
// tree 内范围控制数量 treeCount=nodes*keyStructure=268435455
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
// 存储格式 {dataDir}/checkbook/{dataName}/{formName}/{formName}.dat/idx...
//
// 索引格式
type shopper struct {
	autoID    uint32           // 自增id
	database  Database         // 数据库对象
	name      string           // 表名，根据需求可以随时变化
	id        string           // 表唯一ID，不能改变
	indexes   map[string]Index // 索引ID集合
	fileIndex int              // 数据文件存储编号
	comment   string           // 描述
	nodes     []Nodal          // 节点
	formType  string           // 表类型 SQL/Doc
	fLock     sync.RWMutex
}

func (s *shopper) getAutoID() *uint32 {
	return &s.autoID
}

func (s *shopper) getID() string {
	return s.id
}

func (s *shopper) getName() string {
	return s.name
}

func (s *shopper) getDatabase() Database {
	return s.database
}

func (s *shopper) getFileIndex() int {
	return s.fileIndex
}

func (s *shopper) getIndexes() map[string]Index {
	return s.indexes
}

func (s *shopper) getFormType() string {
	return s.formType
}

func (s *shopper) getDatabaseID() string {
	return s.database.getID()
}

func (s *shopper) lock() {
	s.fLock.Lock()
}

func (s *shopper) unLock() {
	s.fLock.Unlock()
}

func (s *shopper) rLock() {
	s.fLock.RLock()
}

func (s *shopper) rUnLock() {
	s.fLock.RUnlock()
}
