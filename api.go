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
	// FormTypeSQL 关系型数据存储方式
	FormTypeSQL = "FORM_TYPE_SQL"
	// FormTypeDoc 文档型数据存储方式
	FormTypeDoc = "FORM_TYPE_DOC"
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
	// GetDatabases 获取数据库集合
	GetDatabases() []Database
	// CreateDatabase 新建数据库
	//
	// 新建数据库会同时创建一个名为_default的表，未指定表明的情况下使用put/get等方法会操作该表
	//
	// name 数据库名称
	CreateDatabase(name, comment string) (Database, error)
	// CreateForm 创建表
	//
	// databaseName 数据库名
	//
	// 默认自增ID索引
	//
	// name 表名称
	//
	// comment 表描述
	CreateForm(databaseName, formName, comment, formType string) error
	// CreateKey 新建主键
	//
	// databaseName 数据库名
	//
	// name 表名称
	//
	// keyStructure 主键结构名，按照规范结构组成的主键字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	CreateKey(databaseName, formName string, keyStructure string) error
	// createIndex 新建索引
	//
	// databaseName 数据库名
	//
	// name 表名称
	//
	// keyStructure 索引结构名，按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	CreateIndex(databaseName, formName string, keyStructure string) error
	// PutD 新增数据
	//
	// 向_default表中新增一条数据，key相同则返回一个Error
	//
	// keyStructure 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	PutD(key string, value interface{}) (uint64, error)
	// SetD 设置数据，如果存在将被覆盖，如果不存在，则新建
	//
	// 向_default表中新增或更新一条数据，key相同则覆盖
	//
	// keyStructure 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	SetD(key string, value interface{}) (uint64, error)
	// GetD 获取数据
	//
	// 向_default表中查询一条数据并返回
	//
	// keyStructure 插入数据唯一key
	GetD(key string) (interface{}, error)
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
	Put(databaseName, formName, key string, value interface{}) (uint64, error)
	// Set 设置数据，如果存在将被覆盖，如果不存在，则新建
	//
	// 向指定表中新增或更新一条数据，key相同则覆盖
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
	Set(databaseName, formName, key string, value interface{}) (uint64, error)
	// Get 获取数据
	//
	// 向指定表中查询一条数据并返回
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// keyStructure 插入数据唯一key
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
	Insert(databaseName, formName string, value interface{}) (uint64, error)
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
	// Select 获取数据
	//
	// 向指定表中查询一条数据并返回
	//
	// databaseName 数据库名
	//
	// formName 表名
	//
	// selector 条件选择器
	Select(databaseName, formName string, selector *Selector) (int, interface{}, error)
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
	// getForms 获取数据库表集合
	getForms() map[string]Form
	// createForm 新建表方法
	//
	// 默认自增ID索引
	//
	// name 表名称
	//
	// comment 表描述
	createForm(formName, comment, formType string) error
	// createIndex 新建主键
	//
	// name 表名称
	//
	// keyStructure 主键结构名，按照规范结构组成的主键字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	createKey(formName string, keyStructure string) error
	// createIndex 新建索引
	//
	// name 表名称
	//
	// keyStructure 索引结构名，按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	createIndex(formName string, keyStructure string) error
	// Put 新增数据
	//
	// 向_default表中新增一条数据，key相同则覆盖
	//
	// keyStructure 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	//
	// update 本次是否执行更新操作
	put(formName string, key string, value interface{}, update bool) (uint64, error)
	// Get 获取数据
	//
	// 向_default表中查询一条数据并返回
	//
	// keyStructure 插入数据唯一key
	get(formName string, key string) (interface{}, error)
	// Insert 新增数据
	//
	// 向指定表中新增一条数据，key相同则覆盖
	//
	// formName 表名
	//
	// keyStructure 插入数据唯一key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	insert(formName string, value interface{}, update bool) (uint64, error)
	// querySelector 根据条件检索
	//
	// formName 表名
	//
	// selector 条件选择器
	//
	// int 返回检索条目数量
	query(formName string, selector *Selector) (int, interface{}, error)
}

// Form 表接口
//
// 提供表基本操作方法
type Form interface {
	WriteLocker
	getAutoID() *uint64           // getAutoID 返回表当前自增ID值
	getID() string                // getID 返回表唯一ID
	getName() string              // getName 返回表名称
	getDatabase() Database        // getDatabase 返回数据库对象
	getIndexes() map[string]Index // getIndexes 获取表下索引集合
	getFormType() string          // getFormType 获取表类型
}

// Index 索引接口
type Index interface {
	WriteLocker
	// getID 索引唯一ID
	getID() string
	// isPrimary 是否主键
	isPrimary() bool
	// getKey 索引字段名称，由对象结构层级字段通过'.'组成，如
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
	getKeyStructure() string
	// getForm 索引所属表对象
	getForm() Form
	getNode() Nodal // getNode 获取树根节点
	// put 插入数据
	//
	// key 真实key，必须string类型
	//
	// hashKey 索引key，可通过hash转换string生成
	//
	// value 存储对象
	//
	// update 本次是否执行更新操作
	put(key string, hashKey uint64, update bool) IndexBack
	// get 获取数据，返回存储对象
	//
	// key 真实key，必须string类型
	//
	// hashKey 索引key，可通过hash转换string生成
	get(key string, hashKey uint64) (interface{}, error)
	// recover 重置索引数据
	recover()
}

// Nodal 节点对象接口
type Nodal interface {
	WriteLocker      // WriteLocker 读写锁接口
	getIndex() Index // 获取索引对象
	// put 插入数据
	//
	// key 真实key，必须string类型
	//
	// hashKey 索引key，可通过hash转换string生成
	//
	// flexibleKey 下一级最左最小树所对应真实key
	//
	// value 存储对象
	//
	// update 本次是否执行更新操作
	put(key string, hashKey, flexibleKey uint64, update bool) IndexBack
	// get 获取数据，返回存储对象
	//
	// key 真实key，必须string类型
	//
	// hashKey 索引key，可通过hash转换string生成
	//
	// flexibleKey 下一级最左最小树所对应真实key
	get(key string, hashKey, flexibleKey uint64) (interface{}, error)
	getDegreeIndex() uint16 // getDegreeIndex 获取节点所在树中度集合中的数组下标
	getPreNode() Nodal      // getPreNode 获取父节点对象
	getNodes() []Nodal      // getNodes 获取下属节点集合
}

// Leaf 叶子节点对象接口
type Leaf interface {
	Nodal
	getLinks() []Link // getLinks 获取叶子节点下的链表对象集合
}

// Link 叶子节点下的链表对象接口
type Link interface {
	WriteLocker                   // WriteLocker 读写锁接口
	setMD5Key(md5Key string)      // 设置md5Key
	setSeekStartIndex(seek int64) // 设置索引最终存储在文件中的起始位置
	setSeekStart(seek uint32)     // 设置value最终存储在文件中的起始位置
	setSeekLast(seek int)         // 设置value最终存储在文件中的持续长度
	getNodal() Nodal              // box 所属 node
	getMD516Key() string          // 获取md516Key
	getSeekStartIndex() int64     // 索引最终存储在文件中的起始位置
	getSeekStart() uint32         // value最终存储在文件中的起始位置
	getSeekLast() int             // value最终存储在文件中的持续长度
	put(key string, hashKey uint64) *indexBack
	get() (interface{}, error)
}

// IndexBack 索引检索回调结果接口
type IndexBack interface {
	getFormIndexFilePath() string // 索引文件所在路径
	getLocker() WriteLocker       // 索引文件所对应level2层级度节点
	getLink() Link                // 索引对应节点对象子集
	getKey() string               // 索引对应字符串key
	getHashKey() uint64           // put hash keyStructure
	getErr() error
}

// WriteLocker 读写锁接口
type WriteLocker interface {
	lock()    // lock 写锁
	unLock()  // unLock 写解锁
	rLock()   // rLock 读锁
	rUnLock() // rUnLock 读解锁
}
