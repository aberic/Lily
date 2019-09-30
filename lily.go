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
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/lily/api"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"strings"
	"sync"
)

const (
	sysDatabase = "lily"     // 跟随‘Lily’创建的默认库
	userForm    = "_user"    // 跟随‘sysDatabase’库创建的‘Lily’用户管理表
	defaultForm = "_default" // 跟随‘sysDatabase’库创建的‘Lily’k-v表
)

var (
	lilyInstance *Lily
	onceLily     sync.Once
	// ErrDatabaseExist 自定义error信息
	ErrDatabaseExist = errors.New("database already exist")
	// ErrFormExist 自定义error信息
	ErrFormExist = errors.New("form already exist")
	// ErrDataIsNil 自定义error信息
	ErrDataIsNil = errors.New("database had never been created")
	// ErrKeyIsNil 自定义error信息
	ErrKeyIsNil = errors.New("put keyStructure can not be nil")
)

// Lily 祖宗！
//
// 全库唯一常住内存对象，并持有所有库的句柄
//
// API 入口
//
// 存储格式 {dataDir}/Data/{dataName}/{formName}/{formName}.dat/idx...
type Lily struct {
	lilyData  *api.Lily
	databases map[string]Database
	once      sync.Once
	lock      sync.Mutex
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
			lilyData:  &api.Lily{Databases: map[string]*api.Database{}},
			databases: map[string]Database{},
		}
	})
	return lilyInstance
}

// syncRPC2Store 将 api.Lily 对象同步至本地文件中
func (l *Lily) syncRPC2Store() {
	defer l.lock.Unlock()
	l.lock.Lock()
	data, err := proto.Marshal(l.lilyData)
	if nil != err {
		return
	}
	_, _ = gnomon.File().Append(obtainConf().lilyFilePath, data, true)
}

// Start 启动lily
//
// 调用后执行 initialize() 初始化方法
func (l *Lily) Start() {
	gnomon.Log().Info("lily service starting")
	l.initialize()
}

// Stop 停止lily
func (l *Lily) Stop() {
	// todo 停止lily
}

// Restart 重新启动lily
//
// 调用 Restart() 会恢复 Lily 的索引，如果 Lily 索引存在，则 Restart() 什么也不会做
func (l *Lily) Restart() {
	defer l.lock.Unlock()
	l.lock.Lock()
	if gnomon.File().PathExists(obtainConf().lilyFilePath) {
		var (
			data []byte
			lily api.Lily
			err  error
		)
		if data, err = ioutil.ReadFile(obtainConf().lilyFilePath); nil != err {
			gnomon.Log().Panic("restart failed, file read error", gnomon.Log().Err(err))
		}
		if err = proto.Unmarshal(data, &lily); nil != err {
			gnomon.Log().Panic("restart failed, proto unmarshal error", gnomon.Log().Err(err))
		}
		l.lilyData = &lily
		l.recover()
		return
	}
	l.initialize()
}

// recover Lily恢复数据
func (l *Lily) recover() {
	var wg sync.WaitGroup
	l.databases = map[string]Database{}
	for dk, dv := range l.lilyData.Databases {
		l.databases[dk] = &database{
			id:      dv.Id,
			name:    dv.Name,
			comment: dv.Comment,
			forms:   map[string]Form{},
			lily:    l,
		}
		for fk, fv := range dv.Forms {
			var formType string
			switch fv.FormType {
			default:
				formType = FormTypeSQL
			case api.FormType_Doc:
				formType = FormTypeDoc
			}
			l.databases[dk].getForms()[fk] = &form{
				id:       fv.Id,
				name:     fv.Name,
				autoID:   0,
				comment:  fv.Comment,
				formType: formType,
				database: l.databases[dk],
				indexes:  map[string]Index{},
			}
			for ik, iv := range fv.Indexes {
				index := &index{id: iv.Id, primary: iv.Primary, keyStructure: iv.KeyStructure, form: l.databases[dk].getForms()[fk]}
				node := &node{level: 1, degreeIndex: 0, preNode: nil, nodes: []Nodal{}, index: index}
				index.node = node
				l.databases[dk].getForms()[fk].getIndexes()[ik] = index
				wg.Add(1)
				go func(l *Lily, dk, fk, ik string) {
					defer wg.Done()
					l.databases[dk].getForms()[fk].getIndexes()[ik].recover()
				}(l, dk, fk, ik)
			}
		}
	}
	wg.Wait()
}

// initialize 初始化默认库及默认表
func (l *Lily) initialize() {
	l.once.Do(func() {
		gnomon.Log().Info("lily service is initializing")
		gnomon.Log().Info("lily service is creating default database")
		data, err := l.CreateDatabase(sysDatabase, "跟随‘Lily’创建的默认库")
		if nil != err {
			if err == ErrDatabaseExist {
				l.Restart()
				return
			}
			panic(err)
		}
		gnomon.Log().Info(strings.Join([]string{"lily service have been created default database ", sysDatabase}, ""))
		gnomon.Log().Info(strings.Join([]string{"lily service is creating default form ", userForm}, ""))
		if err = data.createForm(userForm, "default user form", FormTypeSQL); nil != err {
			_ = rmDataDir(sysDatabase)
			return
		}
		gnomon.Log().Info(strings.Join([]string{"lily service have been created ", userForm}, ""))
		gnomon.Log().Info(strings.Join([]string{"lily service is creating default form ", defaultForm}, ""))
		if err = data.createForm(defaultForm, "default Data form", FormTypeDoc); nil != err {
			_ = rmDataDir(sysDatabase)
			return
		}
		gnomon.Log().Info(strings.Join([]string{"lily service have been created ", defaultForm}, ""))
		l.databases[sysDatabase] = data
	})
}

// GetDatabases 获取数据库集合
func (l *Lily) GetDatabases() []Database {
	var dbs []Database
	for _, db := range l.databases {
		dbs = append(dbs, db)
	}
	return dbs
}

// CreateDatabase 新建数据库
//
// 新建数据库会同时创建一个名为_default的表，未指定表明的情况下使用put/get等方法会操作该表
//
// name 数据库名称
//
// comment 数据库描述
func (l *Lily) CreateDatabase(name, comment string) (Database, error) {
	// 确定库名不重复
	for k := range l.databases {
		if k == name {
			return nil, ErrDatabaseExist
		}
	}
	// 确保数据库唯一ID不重复
	id := l.name2id(name)
	if err := mkDataDir(id); nil != err {
		return nil, err
	}
	l.databases[name] = &database{name: name, id: id, comment: comment, forms: map[string]Form{}, lily: l}
	// 同步数据到 pb.Lily
	l.lilyData.Databases[name] = &api.Database{Id: id, Name: name, Comment: comment, Forms: map[string]*api.Form{}}
	l.syncRPC2Store()
	return l.databases[name], nil
}

// CreateForm 创建表
//
// 默认自增ID索引
//
// name 表名称
//
// comment 表描述
func (l *Lily) CreateForm(databaseName, formName, comment, formType string) error {
	if database := l.databases[databaseName]; nil != database {
		if err := database.createForm(formName, comment, formType); nil != err {
			return err
		}
		l.syncRPC2Store()
		return nil
	}
	return ErrDataIsNil
}

// CreateKey 新建主键
//
// databaseName 数据库名
//
// name 表名称
//
// keyStructure 主键结构名，按照规范结构组成的主键字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
func (l *Lily) CreateKey(databaseName, formName string, keyStructure string) error {
	if database := l.databases[databaseName]; nil != database {
		if err := database.createIndex(formName, keyStructure); nil != err {
			return err
		}
		l.syncRPC2Store()
		return nil
	}
	return ErrDataIsNil
}

// CreateIndex 新建索引
//
// databaseName 数据库名
//
// name 表名称
//
// keyStructure 索引结构名，按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
func (l *Lily) CreateIndex(databaseName, formName string, keyStructure string) error {
	if database := l.databases[databaseName]; nil != database {
		if err := database.createIndex(formName, keyStructure); nil != err {
			return err
		}
		l.syncRPC2Store()
		return nil
	}
	return ErrDataIsNil
}

// PutD 新增数据
//
// 向_default表中新增一条数据，key相同则返回一个Error
//
// keyStructure 插入数据唯一key
//
// value 插入数据对象
//
// 返回 hashKey
func (l *Lily) PutD(key string, value interface{}) (uint64, error) {
	if gnomon.String().IsEmpty(key) {
		return 0, ErrKeyIsNil
	}
	return l.databases[sysDatabase].put(defaultForm, key, value, false)
}

// SetD 新增数据
//
// 向_default表中新增一条数据，key相同则覆盖
//
// keyStructure 插入数据唯一key
//
// value 插入数据对象
//
// 返回 hashKey
func (l *Lily) SetD(key string, value interface{}) (uint64, error) {
	if gnomon.String().IsEmpty(key) {
		return 0, ErrKeyIsNil
	}
	return l.databases[sysDatabase].put(defaultForm, key, value, true)
}

// GetD 获取数据
//
// 向_default表中查询一条数据并返回
//
// keyStructure 插入数据唯一key
func (l *Lily) GetD(key string) (interface{}, error) {
	return l.databases[sysDatabase].get(defaultForm, key)
}

// Put 新增数据
//
// 向指定表中新增一条数据，key相同则返回一个Error
//
// databaseName 数据库名
//
// formName 表名
//
// keyStructure 插入数据唯一key
//
// value 插入数据对象
//
// 返回 hashKey
func (l *Lily) Put(databaseName, formName, key string, value interface{}) (uint64, error) {
	if gnomon.String().IsEmpty(key) {
		return 0, ErrKeyIsNil
	}
	if nil == l || nil == l.databases[databaseName] {
		return 0, ErrDataIsNil
	}
	return l.databases[databaseName].put(formName, key, value, false)
}

// Set 新增数据
//
// 向指定表中新增一条数据，key相同则覆盖
//
// databaseName 数据库名
//
// formName 表名
//
// keyStructure 插入数据唯一key
//
// value 插入数据对象
//
// 返回 hashKey
func (l *Lily) Set(databaseName, formName, key string, value interface{}) (uint64, error) {
	if gnomon.String().IsEmpty(key) {
		return 0, ErrKeyIsNil
	}
	if nil == l || nil == l.databases[databaseName] {
		return 0, ErrDataIsNil
	}
	return l.databases[databaseName].put(formName, key, value, true)
}

// Get 获取数据
//
// 向指定表中查询一条数据并返回
//
// databaseName 数据库名
//
// formName 表名
//
// keyStructure 插入数据唯一key
func (l *Lily) Get(databaseName, formName, key string) (interface{}, error) {
	if nil == l || nil == l.databases[databaseName] {
		return 0, ErrDataIsNil
	}
	return l.databases[databaseName].get(formName, key)
}

// Insert 新增数据
//
// 向指定表中新增一条数据，key相同则返回一个Error
//
// formName 表名
//
// keyStructure 插入数据唯一key
//
// value 插入数据对象
func (l *Lily) Insert(databaseName, formName string, value interface{}) (uint64, error) {
	if nil == l || nil == l.databases[databaseName] {
		return 0, ErrDataIsNil
	}
	return l.databases[databaseName].insert(formName, value, false)
}

// Update 更新数据
//
// 向指定表中更新一条数据，key相同则覆盖
//
// databaseName 数据库名
//
// formName 表名
//
// value 插入数据对象
func (l *Lily) Update(databaseName, formName string, value interface{}) error {
	if nil == l || nil == l.databases[databaseName] {
		return ErrDataIsNil
	}
	_, err := l.databases[databaseName].insert(formName, value, true)
	return err
}

// Select 获取数据
//
// 向指定表中查询一条数据并返回
//
// formName 表名
//
// keyStructure 插入数据唯一key
func (l *Lily) Select(databaseName, formName string, selector *Selector) (int, interface{}, error) {
	if nil == l || nil == l.databases[databaseName] {
		return 0, nil, ErrDataIsNil
	}
	return l.databases[databaseName].query(formName, selector)
}

// Delete 删除数据
//
// 向指定表中删除一条数据并返回
//
// databaseName 数据库名
//
// formName 表名
//
// selector 条件选择器
func (l *Lily) Delete(databaseName, formName string, selector *Selector) error {
	// todo 删除数据
	if nil == l || nil == l.databases[databaseName] {
		return ErrDataIsNil
	}
	_, _, err := l.databases[databaseName].query(formName, selector)
	return err
}

// name2id 确保数据库唯一ID不重复
func (l *Lily) name2id(name string) string {
	id := gnomon.CryptoHash().MD516(name)
	have := true
	for have {
		have = false
		for _, v := range l.databases {
			if v.getID() == id {
				have = true
				id = gnomon.CryptoHash().MD516(strings.Join([]string{id, gnomon.String().RandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}
