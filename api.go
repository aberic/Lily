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

const (
	formTypeSQL = "FORM_TYPE_SQL" // formTypeSQL 关系型数据存储方式
	formTypeDoc = "FORM_TYPE_DOC" // formTypeDoc 文档型数据存储方式
)

// API 暴露公共API接口
//
// 提供通用 k-v 方法，无需创建新的数据库和表等对象
//
// 在创建 Lily 服务的时候，会默认创建‘sysDatabase’库，同时在该库中创建‘defaultForm’表
//
// 该接口的数据默认在上表中进行操作
type API interface {
	// Start 启动lily
	Start()
	// Restart 重新启动lily
	Restart()
	// CreateDatabase 新建数据库
	//
	// 新建数据库会同时创建一个名为_default的表，未指定表明的情况下使用put/get等方法会操作该表
	//
	// name 数据库名称
	CreateDatabase(name string) (Database, error)
	// CreateForm 创建表
	//
	// 默认自增ID索引
	//
	// name 表名称
	//
	// comment 表描述
	CreateForm(databaseName, formName, comment, formType string) error
	// PutD 新增数据
	//
	// 向_default表中新增一条数据，key相同则返回一个Error
	//
	// key 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	PutD(key string, value interface{}) (uint32, error)
	// SetD 设置数据，如果存在将被覆盖，如果不存在，则新建
	//
	// 向_default表中新增或更新一条数据，key相同则覆盖
	//
	// key 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	SetD(key string, value interface{}) (uint32, error)
	// GetD 获取数据
	//
	// 向_default表中查询一条数据并返回
	//
	// key 插入数据唯一key
	GetD(key string) (interface{}, error)
	// Put 新增数据
	//
	// 向指定表中新增一条数据，key相同则返回一个Error
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// key 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Put(databaseName, formName, key string, value interface{}) (uint32, error)
	// Set 设置数据，如果存在将被覆盖，如果不存在，则新建
	//
	// 向指定表中新增或更新一条数据，key相同则覆盖
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// key 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Set(databaseName, formName, key string, value interface{}) (uint32, error)
	// Get 获取数据
	//
	// 向指定表中查询一条数据并返回
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// key 插入数据唯一key
	Get(databaseName, formName, key string) (interface{}, error)
	// Insert 新增数据
	//
	// 向指定表中新增一条数据，key相同则返回一个Error
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// value 插入数据对象
	Insert(databaseName, formName string, value interface{}) (uint32, error)
	// Update 更新数据
	//
	// 向指定表中新增或更新一条数据，key相同则覆盖
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// value 插入数据对象
	Update(databaseName, formName string, value interface{}) error
	// Query 获取数据
	//
	// 向指定表中查询一条数据并返回
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// selector 条件选择器
	Query(databaseName, formName string, selector *Selector) (interface{}, error)
	// Delete 删除数据
	//
	// 向指定表中删除一条数据并返回
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// selector 条件选择器
	Delete(databaseName, formName string, selector *Selector) error
}

// Database 数据库接口
//
// 提供数据库基本操作方法
type Database interface {
	// getID 返回数据库唯一ID
	getID() string
	// getName 返回数据库名称
	getName() string
	// createForm 新建表方法
	//
	// 默认自增ID索引
	//
	// name 表名称
	//
	// comment 表描述
	createForm(formName, comment, formType string) error
	// Put 新增数据
	//
	// 向_default表中新增一条数据，key相同则覆盖
	//
	// key 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	//
	// update 本次是否执行更新操作
	put(formName string, key string, value interface{}, update bool) (uint32, error)
	// Get 获取数据
	//
	// 向_default表中查询一条数据并返回
	//
	// key 插入数据唯一key
	get(formName string, key string) (interface{}, error)
	// Insert 新增数据
	//
	// 向指定表中新增一条数据，key相同则覆盖
	//
	// formName 表名
	//
	// key 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	insert(formName string, value interface{}, update bool) (uint32, error)
	// querySelector 根据条件检索
	//
	// formName 表名
	//
	// selector 条件选择器
	query(formName string, selector *Selector) (interface{}, error)
}

// Form 表接口
//
// 提供表基本操作方法
type Form interface {
	Data                           // Data 表内数据操作接口
	getAutoID() *uint32            // getAutoID 返回表当前自增ID值
	getID() string                 // getID 返回表唯一ID
	getName() string               // getName 返回表名称
	getFileIndex() int             // getFileIndex 获取表索引文件ID，该ID根据容量满载自增
	getIndexes() map[string]*index // getIndexes 获取表下索引集合
	getFormType() string           // getFormType 获取表类型
}

type Index interface {
	Data
	// id 索引唯一ID
	getID() string
	// 索引字段名称，由对象结构层级字段通过'.'组成，如
	//
	// ref := &ref{
	//		i: 1,
	//		s: "2",
	//		in: refIn{
	//			i: 3,
	//			s: "4",
	//		},
	//	}
	//
	// key可取'i','in.s'
	getKey() string
}

// Nodal 节点对象接口
type Nodal interface {
	Data                           // Data 表内数据操作接口
	existChild(index uint8) bool   // existChild 根据下标判定是否存在子节点
	createChild(index uint8) Nodal // createChild 根据下标创建新的子节点
	getFlexibleKey() uint32        // getFlexibleKey 下一级最左最小树所对应真实key
	getDegreeIndex() uint8         // getDegreeIndex 获取节点所在树中度集合中的数组下标
	getPreNodal() Nodal            // getPreNodal 获取父节点对象
}

type IndexBack interface {
	getFormIndexFilePath() string // 索引文件所在路径
	getNodal() Nodal              // 索引文件所对应level2层级度节点
	getThing() *thing             // 索引对应节点对象子集
	getHashKey() uint32           // put hash key
	getErr() error
}

// Data 表内数据操作接口
//
// 表对象、节点对象、叶子结点对象以及存储节点对象都会实现该接口
type Data interface {
	// put 插入数据
	//
	// originalKey 真实key，必须string类型
	//
	// key 索引key，可通过hash转换string生成
	//
	// value 存储对象
	//
	// update 本次是否执行更新操作
	put(indexID string, originalKey string, key uint32, value interface{}, update bool) IndexBack
	// get 获取数据，返回存储对象
	//
	// originalKey 真实key，必须string类型
	//
	// key 索引key，可通过hash转换string生成
	get(originalKey string, key uint32) (interface{}, error)
	// childCount binaryMatcher 二分查询辅助方法，获取子节点集合数量
	childCount() int
	// child binaryMatcher 二分查询辅助方法，根据子节点集合下标获取树-度对象
	child(index int) Nodal
	// lock 写锁
	lock()
	// unLock 写解锁
	unLock()
	// rLock 读锁
	rLock()
	// rUnLock 读解锁
	rUnLock()
}
