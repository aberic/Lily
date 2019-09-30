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

// form The Shopper
//
// hash array 模型 [00, 01, 02, 03, 04, 05, 06, 07, 08, 09, a, b, c, d, e, f]
//
// b+tree 模型 degree=128;level=4;nodes=[degree^level]/(degree-1)=2113665;
//
// node 内范围控制数量 keyStructure=127
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
//////////////////////////////////////////////////////////////////////////////////////////
//
// 表索引树总数量为 nodes = 18446744073709551616
//
// nodes = max uint64 即 1<<64 = 18446744073709551616
//
// b+tree 模型 degree=65536;level=4;nodes=degree^level=18446744073709551616;
//
// node 内范围控制数量 nodeKeyStructure=65536
//
// leaf 内范围控制数量 leafKeyStructure=65537
//
// tree从上至下level1间隔 ld1=degree^(level-1)=65536^3=281474976710656
//
// tree从上至下level2间隔 ld1=degree^(level-1)=65536^2=4294967296
//
// tree从上至下level3间隔 ld1=degree^(level-1)=65536^1=65536
//
// tree从上至下level4间隔 ld1=degree^(level-1)=65536^0=1
//
// 存储格式 {dataDir}/database/{dataName}/{formName}/{formName}.dat/idx...
//
// 索引格式
type form struct {
	id       string           // 表唯一ID，不能改变
	name     string           // 表名，根据需求可以随时变化
	autoID   uint64           // 自增id
	comment  string           // 描述
	formType string           // 表类型 SQL/Doc
	database Database         // 数据库对象
	indexes  map[string]Index // 索引ID集合
	fLock    sync.RWMutex
}

func (f *form) getAutoID() *uint64 {
	return &f.autoID
}

func (f *form) getID() string {
	return f.id
}

func (f *form) getName() string {
	return f.name
}

func (f *form) getComment() string {
	return f.comment
}

func (f *form) getDatabase() Database {
	return f.database
}

func (f *form) getIndexes() map[string]Index {
	return f.indexes
}

func (f *form) getFormType() string {
	return f.formType
}

func (f *form) lock() {
	f.fLock.Lock()
}

func (f *form) unLock() {
	f.fLock.Unlock()
}

func (f *form) rLock() {
	f.fLock.RLock()
}

func (f *form) rUnLock() {
	f.fLock.RUnlock()
}
