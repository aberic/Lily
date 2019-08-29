/*
 * Copyright (c) 2019. Aberic - All Rights Reservec.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliec.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package Lily

import (
	"errors"
	"github.com/ennoo/rivet/utils/cryptos"
	str "github.com/ennoo/rivet/utils/string"
	"strings"
)

// checkbook 数据库对象
//
// 存储格式 {dataDir}/checkbook/{dataName}/{shopperName}/{shopperName}.dat/idx...
type checkbook struct {
	name     string              // 数据库名称，根据需求可以随时变化
	id       string              // 数据库唯一ID，不能改变
	shoppers map[string]*shopper // 表集合
}

// CreateShopper 创建表
//
// name 表名称
//
// comment 表描述
//
// sequence 是否启用自增ID索引
func (c *checkbook) createShopper(shopperName, comment string, sequence bool) error {
	// 确定库名不重复
	for k := range c.shoppers {
		if k == shopperName {
			return shopperExistErr
		}
	}
	// 确保表唯一ID不重复
	id := c.name2id(shopperName)
	if err := mkFormPath(c.id, id); nil != err {
		return err
	}
	c.shoppers[shopperName] = newShopper(shopperName, id, comment, c)
	if sequence {
		sequenceName := c.sequenceName(shopperName)
		sequenceId := c.name2id(sequenceName)
		c.shoppers[sequenceName] = newShopper(sequenceName, sequenceId, comment, c)
	}
	return nil
}

// InsertInt 新增数据
//
// 向指定表中新增一条数据，key相同则覆盖
//
// shopperName 表名
//
// key 插入数据唯一key
//
// value 插入数据对象
func (c *checkbook) InsertInt(shopperName string, key int, value interface{}) (uint32, error) {
	if nil == c {
		return 0, errorDataIsNil
	}
	return c.insert(shopperName, Key(key), uint32(key), value)
}

// QueryInt 获取数据
//
// 向指定表中查询一条数据并返回
//
// shopperName 表名
//
// key 插入数据唯一key
func (c *checkbook) QueryInt(shopperName string, key int) (interface{}, error) {
	if nil == c {
		return nil, errorDataIsNil
	}
	return c.query(shopperName, Key(key), uint32(key))
}

// Insert 新增数据
//
// 向指定表中新增一条数据，key相同则覆盖
//
// shopperName 表名
//
// key 插入数据唯一key
//
// value 插入数据对象
func (c *checkbook) Insert(shopperName string, key Key, value interface{}) (uint32, error) {
	if nil == c {
		return 0, errorDataIsNil
	}
	return c.insert(shopperName, key, hash(key), value)
}

// Query 获取数据
//
// 向指定表中查询一条数据并返回
//
// shopperName 表名
//
// key 插入数据唯一key
func (c *checkbook) Query(shopperName string, key Key) (interface{}, error) {
	if nil == c {
		return nil, errorDataIsNil
	}
	//return l.get(key, hash(key))
	return c.query(shopperName, key, hash(key))
}

// Insert 新增数据
//
// 向指定表中新增一条数据，key相同则覆盖
//
// shopperName 表名
//
// key 插入数据唯一key
//
// value 插入数据对象
//
// 返回 hashKey
func (c *checkbook) insert(shopperName string, key Key, hashKey uint32, value interface{}) (uint32, error) {
	l := c.shoppers[shopperName]
	if nil == l || nil == l.purses {
		return 0, shopperIsInvalid(shopperName)
	}
	sequenceName := c.sequenceName(shopperName)
	if nil == c.shoppers[sequenceName] {
		return hashKey, l.put(key, hashKey, value)
	} else {
		var (
			ls  *shopper
			err error
			//wg       sync.WaitGroup
			//checkErr chan error
		)
		ls = c.shoppers[sequenceName]
		//checkErr = make(chan error, 2)
		//wg.Add(2)
		//err = pool().submit(func() {
		//	defer wg.Done()
		//	err := l.put(key, hashKey, value)
		//	if nil != err {
		//		checkErr <- err
		//	} else {
		//		checkErr <- nil
		//	}
		//})
		//if nil != err {
		//	return 0, err
		//}
		//err = pool().submit(func() {
		//	defer wg.Done()
		//	err := ls.put(key, atomic.AddUint32(&ls.autoID, 1), value)
		//	if nil != err {
		//		checkErr <- err
		//	} else {
		//		checkErr <- nil
		//	}
		//})
		//if nil != err {
		//	return 0, err
		//}
		//wg.Wait()
		//err = <-checkErr
		//if nil == err {
		//	err = <-checkErr
		//} else {
		//	return 0, err
		//}
		//if nil != err {
		//	return 0, err
		//}
		err = l.put(key, hashKey, value)
		if nil != err {
			return 0, err
		}
		err = ls.put(key, hashKey, value)
		if nil != err {
			return 0, err
		}
		// todo 回滚策略待完成
		return hashKey, nil
	}
}

// Query 获取数据
//
// 向指定表中查询一条数据并返回
//
// shopperName 表名
//
// key 插入数据唯一key
func (c *checkbook) query(shopperName string, key Key, hashKey uint32) (interface{}, error) {
	l := c.shoppers[shopperName]
	if nil == l || nil == l.purses {
		return nil, shopperIsInvalid(shopperName)
	}
	return l.get(key, hashKey)
}

// QuerySelector 根据条件检索
//
// shopperName 表名
//
// selector 条件选择器
func (c *checkbook) QuerySelector(shopperName string, selector *Selector) (interface{}, error) {
	if nil == c {
		return nil, errorDataIsNil
	}
	selector.shopperName = shopperName
	selector.checkbook = c
	return selector.query()
}

// shopperIsInvalid 自定义error信息
func shopperIsInvalid(shopperName string) error {
	return errors.New(strings.Join([]string{"invalid name ", shopperName}, ""))
}

// sequenceName 开启自增主键索引后新的组合固定表明
func (c *checkbook) sequenceName(name string) string {
	return strings.Join([]string{name, "id"}, "_")
}

// name2id 确保数据库唯一ID不重复
func (c *checkbook) name2id(name string) string {
	id := cryptos.MD516(name)
	have := true
	for have {
		have = false
		for _, v := range c.shoppers {
			if v.id == id {
				have = true
				id = cryptos.MD516(strings.Join([]string{id, str.RandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}
