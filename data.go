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
	"strings"
	"sync"
)

const (
	defaultLily         = "_default"
	defaultSequenceLily = "_default_sequence"
)

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
//
// sequence 默认表是否启用自增ID索引
func NewData(name string, sequence bool) *Data {
	data := &Data{name: name, lilies: map[string]*lily{}}
	data.lilies[defaultLily] = newLily(defaultLily, "default data lily", data)
	if sequence {
		data.lilies[defaultSequenceLily] = newLily(defaultSequenceLily, "default data lily", data)
	}
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
		return errors.New("data had never been created")
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
func (d *Data) InsertD(key Key, value interface{}) error {
	return d.Insert(defaultLily, key, value)
}

// QueryD 获取数据
//
// 向_default表中查询一条数据并返回
//
// key 插入数据唯一key
func (d *Data) QueryD(key Key) (interface{}, error) {
	return d.Query(defaultLily, key)
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
func (d *Data) Insert(lilyName string, key Key, value interface{}) error {
	if nil == d {
		return errors.New("data had never been created")
	}
	l := d.lilies[lilyName]
	if nil == l || nil == l.cities {
		return errors.New(strings.Join([]string{"group is invalid with name ", lilyName}, ""))
	}
	sequenceName := sequenceName(lilyName)
	if nil == d.lilies[sequenceName] {
		return l.put(key, hash(key), value)
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
			err := l.put(key, hash(key), value)
			if nil != err {
				checkErr <- err
			}
		}(key, value)
		go func(key Key, value interface{}) {
			defer wg.Done()
			err := ls.put(key, hash(key), value)
			if nil != err {
				checkErr <- err
			}
		}(key, value)
		wg.Wait()
		return <-checkErr
	}
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
		return nil, errors.New("data had never been created")
	}
	l := d.lilies[lilyName]
	if nil == l || nil == l.cities {
		return nil, errors.New(strings.Join([]string{"group is invalid with name ", lilyName}, ""))
	}
	return l.get(key, hash(key))
}

func (d *Data) PutGInt(lilyName string, key int, value interface{}) error {
	l := d.lilies[lilyName]
	if nil == l || nil == l.cities {
		return errors.New(strings.Join([]string{"group is invalid with name ", lilyName}, ""))
	}
	return l.put(Key(key), uint32(key), value)
}

func (d *Data) GetGInt(lilyName string, key int) (interface{}, error) {
	l := d.lilies[lilyName]
	if nil == l || nil == l.cities {
		return nil, errors.New(strings.Join([]string{"group is invalid with name ", lilyName}, ""))
	}
	return l.get(Key(key), uint32(key))
}

func sequenceName(name string) string {
	return strings.Join([]string{name, "sequence"}, "_")
}
