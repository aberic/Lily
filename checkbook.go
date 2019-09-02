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
	"github.com/ennoo/rivet/utils/log"
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

func (c *checkbook) createForm(formName, comment string) error {
	// 确定库名不重复
	for k := range c.forms {
		if k == formName {
			return formExistErr
		}
	}
	// 确保表唯一ID不重复
	formID := c.name2id(formName)
	// 自增索引ID
	indexID := c.name2id(strings.Join([]string{formName, "id"}, "_"))
	// 默认自定义Key生成ID
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
		nodes:     []Nodal{},
	}
	return nil
}

func (c *checkbook) createIndex(formName string, key string, value interface{}) (uint32, error) {
	return 0, nil
}

func (c *checkbook) insert(formName string, key string, value interface{}) (uint32, error) {
	form := c.forms[formName] // 获取待操作表
	if nil == form {
		return 0, shopperIsInvalid(formName)
	}
	var (
		chanIndex chan *indexBack
		err       error
	)
	indexIDs := form.getIndexIDs() // 获取表索引ID集合
	indexLen := len(indexIDs)
	chanIndex = make(chan *indexBack, indexLen)     // 创建索引ID结果返回通道
	autoID := atomic.AddUint32(form.getAutoID(), 1) // ID自增
	for _, indexID := range indexIDs {              // 遍历表索引ID集合，检索并计算当前索引所在文件位置
		if err = pool().submitIndex(indexID, func(indexID string) {
			if indexID == c.name2id(strings.Join([]string{formName, "id"}, "_")) {
				chanIndex <- form.put(indexID, strconv.Itoa(int(autoID)), autoID, value)
			} else if indexID == c.name2id(strings.Join([]string{formName, "custom"}, "_")) {
				chanIndex <- form.put(indexID, key, hash(key), value)
			}
		}); nil != err {
			return 0, err
		}
	}
	wrIndexBack := make(chan *writeResult, 1) // 索引存储结果通道
	// 存储数据到表文件
	wf := store().appendForm(form, pathFormDataFile(c.id, form.getID(), form.getFileIndex()), value)
	if nil != wf.err {
		return 0, wf.err
	}
	for i := 0; i < indexLen; i++ {
		ib := <-chanIndex
		if err = pool().submitChanIndex(ib, func(ib *indexBack) {
			md5Key := cryptos.MD516(ib.originalKey) // hash(originalKey) 会发生碰撞，因此这里存储md5结果进行反向验证
			// 写入5位key及16位md5后key
			appendStr := strings.Join([]string{uint32ToDDuoString(ib.key), md5Key}, "")
			log.Self.Debug("insert", log.String("appendStr", appendStr), log.Reflect("formIndexFilePath", ib.formIndexFilePath))
			// 写入5位key及16位md5后key及16位起始seek和8位持续seek
			wr := store().appendIndex(ib.indexNodal, ib.formIndexFilePath, appendStr, wf)
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
		case wrIndex := <-wrIndexBack:
			if nil != wrIndex.err {
				return 0, wrIndex.err
			}
			// todo 回滚策略待完成
			return autoID, nil
		}
	}
}
func (c *checkbook) query(formName string, key string, hashKey uint32) (interface{}, error) {
	form := c.forms[formName]
	if nil == form {
		return nil, shopperIsInvalid(formName)
	}
	return form.get(key, hashKey)
}

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
