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
	"github.com/ennoo/rivet/utils/cryptos"
	"github.com/ennoo/rivet/utils/string"
	"strings"
	"sync"
)

const (
	sysDatabase  = "lily"      // 跟随‘Lily’创建的默认库
	userForm     = "_user"     // 跟随‘sysDatabase’库创建的‘Lily’用户管理表
	databaseForm = "_database" // 跟随‘sysDatabase’库创建的‘Lily’数据库管理表
	indexForm    = "_index"    // 跟随‘sysDatabase’库创建的‘Lily’索引管理表
	defaultForm  = "_default"  // 跟随‘sysDatabase’库创建的‘Lily’k-v表
)

var (
	lilyInstance     *Lily
	onceLily         sync.Once
	databaseExistErr = errors.New("database already exist")          // databaseExistErr 自定义error信息
	formExistErr     = errors.New("form already exist")              // formExistErr 自定义error信息
	errorDataIsNil   = errors.New("database had never been created") // errorDataIsNil 自定义error信息
)

// Lily 祖宗！
//
// 全库唯一常住内存对象，并持有所有库的句柄
//
// API 入口
//
// 存储格式 {dataDir}/Data/{dataName}/{formName}/{formName}.dat/idx...
type Lily struct {
	defaultDatabase Database
	databases       map[string]Database
	once            sync.Once
}

// ObtainLily 获取 Lily 对象
//
// 会初始化一个空 Lily，如果是第一次调用的话
//
// 首次调用后需要执行 initialize() 初始化方法
//
// 或者通过外部调用 Start() 来执行初始化操作
//
// 调用 Restart() 会恢复 Lily 的索引，如果 Lily 索引存在，则 Restart() 什么也不会做
//
// 会返回一个已创建的 Lily，如果非第一次调用的话
func ObtainLily() *Lily {
	onceLily.Do(func() {
		lilyInstance = &Lily{
			databases: map[string]Database{},
		}
	})
	return lilyInstance
}

// Start 启动lily
//
// 调用后执行 initialize() 初始化方法
func (l *Lily) Start() {
	l.initialize()
}

// Restart 重新启动lily
//
// 调用 Restart() 会恢复 Lily 的索引，如果 Lily 索引存在，则 Restart() 什么也不会做
func (l *Lily) Restart() {
	// todo 恢复索引
}

// initialize 初始化默认库及默认表
func (l *Lily) initialize() {
	l.once.Do(func() {
		data, err := l.CreateDatabase(sysDatabase)
		if nil != err {
			if err == databaseExistErr {
				l.Restart()
				return
			} else {
				panic(err)
			}
		}
		if err = data.createForm(userForm, "default user form"); nil != err {
			_ = rmDataDir(sysDatabase)
			return
		}
		if err = data.createForm(databaseForm, "default database form"); nil != err {
			_ = rmDataDir(sysDatabase)
			return
		}
		if err = data.createForm(indexForm, "default index form"); nil != err {
			_ = rmDataDir(sysDatabase)
			return
		}
		if err = data.createForm(defaultForm, "default Data form"); nil != err {
			_ = rmDataDir(sysDatabase)
			return
		}
		l.defaultDatabase = data
	})
}

func (l *Lily) CreateDatabase(name string) (Database, error) {
	// 确定库名不重复
	for k := range l.databases {
		if k == name {
			return nil, databaseExistErr
		}
	}
	// 确保数据库唯一ID不重复
	id := l.name2id(name)
	if err := mkDataDir(id); nil != err {
		return nil, err
	}
	data := &checkbook{name: name, id: id, forms: map[string]Form{}}
	l.databases[name] = data
	//l.defaultDatabase.insert(databaseForm, id, )
	return data, nil
}

func (l *Lily) CreateForm(databaseName, formName, comment string) error {
	if database := l.databases[databaseName]; nil != database {
		return database.createForm(formName, comment)
	}
	return errorDataIsNil
}

func (l *Lily) Put(key string, value interface{}) (uint32, error) {
	if nil == l || nil == l.databases[defaultForm] {
		return 0, errorDataIsNil
	}
	return l.Insert(sysDatabase, defaultForm, key, value)
}

func (l *Lily) Get(key string) (interface{}, error) {
	if nil == l || nil == l.databases[defaultForm] {
		return 0, errorDataIsNil
	}
	return l.Query(sysDatabase, defaultForm, key)
}

func (l *Lily) Insert(databaseName, formName string, key string, value interface{}) (uint32, error) {
	if nil == l || nil == l.databases[databaseName] {
		return 0, errorDataIsNil
	}
	return l.databases[databaseName].insert(formName, key, value)
}

func (l *Lily) Query(databaseName, formName string, key string) (interface{}, error) {
	if nil == l || nil == l.databases[databaseName] {
		return nil, errorDataIsNil
	}
	return l.databases[databaseName].query(formName, key, hash(key))
}

// name2id 确保数据库唯一ID不重复
func (l *Lily) name2id(name string) string {
	id := cryptos.MD516(name)
	have := true
	for have {
		have = false
		for _, v := range l.databases {
			if v.getID() == id {
				have = true
				id = cryptos.MD516(strings.Join([]string{id, str.RandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}
