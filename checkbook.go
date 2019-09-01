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
	"github.com/aberic/common/log"
	"github.com/ennoo/rivet/utils/cryptos"
	"github.com/ennoo/rivet/utils/string"
	"strconv"
	"strings"
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
// 默认自增ID索引
//
// name 表名称
//
// comment 表描述
func (c *checkbook) createForm(formName, comment string) error {
	// 确定库名不重复
	for k := range c.forms {
		if k == formName {
			return formExistErr
		}
	}
	// 确保表唯一ID不重复
	formID := c.name2id(formName)
	indexID := c.name2id(strings.Join([]string{formName, "id"}, "_"))
	customID := c.name2id(strings.Join([]string{formName, "custom"}, "_"))
	fileIndex := 0
	if err := mkFormResource(c.id, formID, indexID, customID, fileIndex); nil != err {
		return err
	}
	c.forms[formName] = &shopper{
		autoID:    0,
		name:      formName,
		id:        formID,
		indexIDs:  []string{indexID, customID},
		fileIndex: fileIndex,
		comment:   comment,
		database:  c,
		nodes:     []nodal{},
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
	var (
		chanIndex chan *indexBack
		err       error
	)
	indexIDs := form.getIndexIDs()
	indexLen := len(indexIDs)
	chanIndex = make(chan *indexBack, indexLen)
	autoID := atomic.AddUint32(form.getAutoID(), 1)
	for _, indexID := range indexIDs {
		if err = pool().submit(func() {
			if indexID == c.name2id(strings.Join([]string{formName, "id"}, "_")) {
				chanIndex <- form.put(indexID, key, autoID, value)
			} else if indexID == c.name2id(strings.Join([]string{formName, "custom"}, "_")) {
				chanIndex <- form.put(indexID, key, hashKey, value)
			}
		}); nil != err {
			return 0, err
		}
	}
	wrTo := make(chan *writeResult, 1)
	wrIndexBack := make(chan *writeResult, 1)
	wrFormBack := make(chan *writeResult, 1)
	if err = pool().submit(func() {
		wrFormBack <- store().appendForm(form, pathFormDataFile(c.id, form.getID(), form.getFileIndex()), value, wrTo)
	}); nil != err {
		return 0, err
	}
	md5Key := cryptos.MD516(string(key))
	for i := 0; i < indexLen; i++ {
		ib := <-chanIndex
		if err = pool().submit(func() {
			appendStr := strings.Join([]string{c.uint32toFullState(autoID), md5Key}, "")
			log.Self.Debug("insert", log.Reflect("formIndexFilePath", ib.formIndexFilePath))
			wr := store().appendIndex(ib.indexNodal, ib.formIndexFilePath, appendStr, wrTo)
			if nil == wr.err {
				ib.thing.md5Key = md5Key
				ib.thing.seekStart = wr.seekStart
				ib.thing.seekLast = wr.seekLast
			}
			wrIndexBack <- wr
		}); nil != err {
			return 0, err
		}
	}
	for {
		select {
		case wrForm := <-wrFormBack:
			if nil != wrForm.err {
				return 0, wrForm.err
			}
		case wrIndex := <-wrIndexBack:
			if nil != wrIndex.err {
				return 0, wrIndex.err
			}
			// todo 回滚策略待完成
			return autoID, nil
		}
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

// uint32toFullState 补全不满十位数状态，如1->0000000001、34->0000000034、215->0000000215
func (c *checkbook) uint32toFullState(index uint32) string {
	pos := 0
	for index > 1 {
		index /= 10
		pos++
	}
	backZero := 10 - pos
	backZeroStr := strconv.Itoa(int(index))
	for i := 0; i < backZero; i++ {
		backZeroStr = strings.Join([]string{"0", backZeroStr}, "")
	}
	return backZeroStr
}
