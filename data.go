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
	s "sort"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	defaultLily         = "_default"
	defaultSequenceLily = "_default_id"
)

// errorDataIsNil 自定义error信息
var errorDataIsNil = errors.New("data had never been created")

// Data 数据库对象
type Data struct {
	name   string           // 数据库名称
	lilies map[string]*lily // 数据库表集合
}

// NewData 新建数据库
//
// 新建数据库会同时创建一个名为_default的表，未指定表明的情况下使用put/get等方法会操作该表
//
// name 数据库名称
func NewData(name string) *Data {
	data := &Data{name: name, lilies: map[string]*lily{}}
	data.lilies[defaultLily] = newLily(defaultLily, "default data lily", data)
	data.lilies[defaultSequenceLily] = newLily(defaultSequenceLily, "default data lily", data)
	return data
}

// CreateLily 创建表
//
// name 表名称
//
// comment 表描述
//
// sequence 是否启用自增ID索引
func (d *Data) CreateLily(name, comment string, sequence bool) error {
	if nil == d {
		return errorDataIsNil
	}
	d.lilies[name] = newLily(name, comment, d)
	if sequence {
		sequenceName := sequenceName(name)
		d.lilies[sequenceName] = newLily(sequenceName, comment, d)
	}
	return nil
}

// InsertD 新增数据
//
// 向_default表中新增一条数据，key相同则覆盖
//
// key 插入数据唯一key
//
// value 插入数据对象
func (d *Data) InsertD(key Key, value interface{}) (uint32, error) {
	if nil == d {
		return 0, errorDataIsNil
	}
	return d.Insert(defaultLily, key, value)
}

// QueryD 获取数据
//
// 向_default表中查询一条数据并返回
//
// key 插入数据唯一key
func (d *Data) QueryD(key Key) (interface{}, error) {
	if nil == d {
		return nil, errorDataIsNil
	}
	return d.Query(defaultLily, key)
}

// InsertGInt 新增数据
//
// 向指定表中新增一条数据，key相同则覆盖
//
// lilyName 表名
//
// key 插入数据唯一key
//
// value 插入数据对象
func (d *Data) InsertGInt(lilyName string, key int, value interface{}) (uint32, error) {
	if nil == d {
		return 0, errorDataIsNil
	}
	return d.insert(lilyName, Key(key), uint32(key), value)
}

// QueryGInt 获取数据
//
// 向指定表中查询一条数据并返回
//
// lilyName 表名
//
// key 插入数据唯一key
func (d *Data) QueryGInt(lilyName string, key int) (interface{}, error) {
	if nil == d {
		return nil, errorDataIsNil
	}
	return d.query(lilyName, Key(key), uint32(key))
}

// Insert 新增数据
//
// 向指定表中新增一条数据，key相同则覆盖
//
// lilyName 表名
//
// key 插入数据唯一key
//
// value 插入数据对象
func (d *Data) Insert(lilyName string, key Key, value interface{}) (uint32, error) {
	if nil == d {
		return 0, errorDataIsNil
	}
	return d.insert(lilyName, key, hash(key), value)
}

// Query 获取数据
//
// 向指定表中查询一条数据并返回
//
// lilyName 表名
//
// key 插入数据唯一key
func (d *Data) Query(lilyName string, key Key) (interface{}, error) {
	if nil == d {
		return nil, errorDataIsNil
	}
	//return l.get(key, hash(key))
	return d.query(lilyName, key, hash(key))
}

// Insert 新增数据
//
// 向指定表中新增一条数据，key相同则覆盖
//
// lilyName 表名
//
// key 插入数据唯一key
//
// value 插入数据对象
func (d *Data) insert(lilyName string, key Key, hashKey uint32, value interface{}) (uint32, error) {
	l := d.lilies[lilyName]
	if nil == l || nil == l.purses {
		return 0, groupIsInvalid(lilyName)
	}
	sequenceName := sequenceName(lilyName)
	if nil == d.lilies[sequenceName] {
		atomic.AddUint32(&l.count, 1)
		return hashKey, l.put(key, hashKey, value)
	} else {
		var (
			ls       *lily
			wg       sync.WaitGroup
			checkErr chan error
		)
		ls = d.lilies[sequenceName]
		checkErr = make(chan error, 2)
		wg.Add(2)
		go func(key Key, value interface{}) {
			defer wg.Done()
			err := l.put(key, hashKey, value)
			if nil != err {
				checkErr <- err
			} else {
				checkErr <- nil
			}
		}(key, value)
		go func(key Key, value interface{}) {
			defer wg.Done()
			err := ls.put(key, atomic.AddUint32(&ls.autoID, 1), value)
			if nil != err {
				checkErr <- err
			} else {
				checkErr <- nil
			}
		}(key, value)
		wg.Wait()
		err := <-checkErr
		// todo 回滚策略待完成
		if nil == err {
			err = <-checkErr
		} else {
			return 0, err
		}
		if nil != err {
			return 0, err
		}
		atomic.AddUint32(&l.count, 1)
		atomic.AddUint32(&ls.count, 1)
		return hashKey, nil
	}
}

// Query 获取数据
//
// 向指定表中查询一条数据并返回
//
// lilyName 表名
//
// key 插入数据唯一key
func (d *Data) query(lilyName string, key Key, hashKey uint32) (interface{}, error) {
	l := d.lilies[lilyName]
	if nil == l || nil == l.purses {
		return nil, groupIsInvalid(lilyName)
	}
	return l.get(key, hashKey)
}

// QuerySelector 根据条件检索
//
// lilyName 表名
//
// selector 条件选择器
func (d *Data) QuerySelector(lilyName string, selector *Selector) (interface{}, error) {
	var l *lily
	if nil == d {
		return nil, errorDataIsNil
	}
	if nil != selector.Indexes {
		s.Stable(selector.Indexes)
		indexStr := ""
		for _, index := range selector.Indexes.IndexArr {
			indexStr = strings.Join([]string{indexStr, index.param}, "")
		}
		lilyIndexName := strings.Join([]string{lilyName, indexStr}, "")
		l = d.lilies[lilyIndexName]
	} else {
		l = d.lilies[lilyName]
	}
	if nil == l || nil == l.purses {
		return nil, groupIsInvalid(lilyName)
	}
	return selector.query(l), nil
}

// groupIsInvalid 自定义error信息
func groupIsInvalid(lilyName string) error {
	return errors.New(strings.Join([]string{"group is invalid with name ", lilyName}, ""))
}

// sequenceName 开启自增主键索引后新的组合固定表明
func sequenceName(name string) string {
	return strings.Join([]string{name, "id"}, "_")
}
