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
	"github.com/ennoo/rivet/utils/string"
	"strings"
	"sync"
	"sync/atomic"
)

// checkbook 数据库对象
//
// 存储格式 {dataDir}/checkbook/{dataName}/{formName}/{formName}.dat/idx...
type checkbook struct {
	name  string          // 数据库名称，根据需求可以随时变化
	id    string          // 数据库唯一ID，不能改变
	forms map[string]Form // 表集合
}

func (c *checkbook) getID() string {
	return c.id
}

func (c *checkbook) getName() string {
	return c.name
}

// createForm 创建表
//
// name 表名称
//
// comment 表描述
//
// sequence 是否启用自增ID索引
func (c *checkbook) createForm(formName, comment string, sequence bool) error {
	// 确定库名不重复
	for k := range c.forms {
		if k == formName {
			return formExistErr
		}
	}
	// 确保表唯一ID不重复
	id := c.name2id(formName)
	if err := mkFormResource(c.id, id); nil != err {
		return err
	}
	c.forms[formName] = &shopper{
		autoID:   0,
		name:     formName,
		id:       id,
		comment:  comment,
		database: c,
		nodes:    []nodal{},
	}
	if sequence {
		sequenceName := c.sequenceName(formName)
		sequenceId := c.name2id(sequenceName)
		c.forms[sequenceName] = &shopper{
			autoID:   0,
			name:     sequenceName,
			id:       sequenceId,
			comment:  comment,
			database: c,
			nodes:    []nodal{},
		}

	}
	return nil
}

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
func (c *checkbook) insert(formName string, key Key, hashKey uint32, value interface{}) (uint32, error) {
	form := c.forms[formName]
	if nil == form {
		return 0, shopperIsInvalid(formName)
	}
	sequenceName := c.sequenceName(formName)
	if nil == c.forms[sequenceName] {
		return hashKey, form.put(key, hashKey, value)
	} else {
		var (
			formSequence Form
			err          error
			wg           sync.WaitGroup
			checkErr     chan error
		)
		formSequence = c.forms[sequenceName]
		checkErr = make(chan error, 2)
		wg.Add(2)
		err = pool().submit(func() {
			defer wg.Done()
			err := form.put(key, hashKey, value)
			if nil != err {
				checkErr <- err
			} else {
				checkErr <- nil
			}
		})
		if nil != err {
			return 0, err
		}
		err = pool().submit(func() {
			defer wg.Done()
			err := formSequence.put(key, atomic.AddUint32(formSequence.getAutoID(), 1), value)
			if nil != err {
				checkErr <- err
			} else {
				checkErr <- nil
			}
		})
		if nil != err {
			return 0, err
		}
		wg.Wait()
		err = <-checkErr
		if nil == err {
			err = <-checkErr
		} else {
			return 0, err
		}
		if nil != err {
			return 0, err
		}
		//err = form.put(key, hashKey, value)
		//if nil != err {
		//	return 0, err
		//}
		//err = formSequence.put(key, hashKey, value)
		//if nil != err {
		//	return 0, err
		//}
		// todo 回滚策略待完成
		return hashKey, nil
	}
}

// Query 获取数据
//
// 向指定表中查询一条数据并返回
//
// formName 表名
//
// key 插入数据唯一key
func (c *checkbook) query(formName string, key Key, hashKey uint32) (interface{}, error) {
	form := c.forms[formName]
	if nil == form {
		return nil, shopperIsInvalid(formName)
	}
	return form.get(key, hashKey)
}

// querySelector 根据条件检索
//
// formName 表名
//
// selector 条件选择器
func (c *checkbook) querySelector(formName string, selector *Selector) (interface{}, error) {
	if nil == c {
		return nil, errorDataIsNil
	}
	selector.formName = formName
	selector.checkbook = c
	return selector.query()
}

// shopperIsInvalid 自定义error信息
func shopperIsInvalid(formName string) error {
	return errors.New(strings.Join([]string{"invalid name ", formName}, ""))
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
		for _, v := range c.forms {
			if v.getID() == id {
				have = true
				id = cryptos.MD516(strings.Join([]string{id, str.RandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}
