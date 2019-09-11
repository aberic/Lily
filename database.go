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

package lily

import (
	"errors"
	"github.com/aberic/gnomon"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	indexAutoID    = "lily_biexuewo_id"
	indexDefaultID = "lily_biexuewo_default"
)

// database 数据库对象
//
// 存储格式 {dataDir}/database/{dataName}/{formName}/{formName}.dat/idx...
type database struct {
	name  string          // 数据库名称，根据需求可以随时变化
	id    string          // 数据库唯一ID，不能改变
	forms map[string]Form // 表集合
}

func (c *database) getID() string {
	return c.id
}

func (c *database) getName() string {
	return c.name
}

// getForms 获取数据库表集合
func (c *database) getForms() map[string]Form {
	return c.forms
}

func (c *database) createForm(formName, comment, formType string) error {
	// 确定库名不重复
	for k := range c.forms {
		if k == formName {
			return ErrFormExist
		}
	}
	// 确保表唯一ID不重复
	formID := c.name2id(formName)
	form := &form{
		autoID:   0,
		name:     formName,
		id:       formID,
		comment:  comment,
		database: c,
		nodes:    []Nodal{},
		formType: formType,
	}
	err := mkFormResource(c.id, formID, 0)
	if nil != err {
		return err
	}
	indexes := make(map[string]Index)
	if formType == formTypeSQL {
		// 自增索引ID
		indexID := c.name2id(strings.Join([]string{formName, indexAutoID}, "_"))
		if err = mkFormIndexResource(c.id, formID, indexID); nil != err {
			return err
		}
		indexes[indexID] = &index{id: indexID, keyStructure: indexAutoID, form: form, fileIndex: 0}
	} else {
		// 默认自定义Key生成ID
		defaultID := c.name2id(strings.Join([]string{formName, indexDefaultID}, "_"))
		if err = mkFormIndexResource(c.id, formID, defaultID); nil != err {
			return err
		}
		indexes[defaultID] = &index{id: defaultID, keyStructure: indexDefaultID, form: form, fileIndex: 0}
	}
	form.indexes = indexes
	c.forms[formName] = form
	return nil
}

func (c *database) createIndex(formName string, keyStructure string) error {
	form := c.forms[formName]
	// 自定义Key生成ID
	customID := c.name2id(strings.Join([]string{formName, keyStructure}, "_"))
	gnomon.Log().Debug("createIndex", gnomon.Log().Field("customID", customID))
	if err := mkFormIndexResource(c.id, form.getID(), customID); nil != err {
		return err
	}
	form.getIndexes()[customID] = &index{id: customID, keyStructure: keyStructure, form: form, fileIndex: 0}
	return nil
}

func (c *database) put(formName string, key string, value interface{}, update bool) (uint32, error) {
	form := c.forms[formName] // 获取待操作表
	if nil == form {
		return 0, shopperIsInvalid(formName)
	}
	if form.getFormType() != formTypeDoc {
		return 0, errors.New("put method only support doc")
	}
	indexes := form.getIndexes()                    // 获取表索引ID集合
	autoID := atomic.AddUint32(form.getAutoID(), 1) // ID自增
	return c.insertDataWithIndexInfo(form, key, autoID, indexes, value, update)
}

func (c *database) get(formName string, key string) (interface{}, error) {
	form := c.forms[formName]
	if nil == form {
		return nil, shopperIsInvalid(formName)
	}
	for _, index := range form.getIndexes() {
		if index.getKeyStructure() == indexDefaultID {
			return index.get(key, hash(key))
		}
	}
	return nil, errors.New("no key for custom id index")
}

func (c *database) insert(formName string, value interface{}, update bool) (uint32, error) {
	// todo
	return 0, nil
}

func (c *database) query(formName string, selector *Selector) (interface{}, error) {
	if nil == c {
		return nil, ErrDataIsNil
	}
	selector.formName = formName
	selector.database = c
	return selector.query()
}

func (c *database) insertDataWithIndexInfo(form Form, key string, autoID uint32, indexes map[string]Index, value interface{}, update bool) (uint32, error) {
	var (
		ibs []IndexBack
		err error
	)
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	if ibs, err = c.rangeIndexes(form, key, autoID, indexes, value, update); nil != err {
		return 0, err
	}
	wrIndexBack := make(chan *writeResult, 1) // 索引存储结果通道
	// 存储数据到表文件
	wf := store().appendForm(form, pathFormDataFile(c.id, form.getID(), form.getFileIndex()), value)
	if nil != wf.err {
		return 0, wf.err
	}
	for _, ib := range ibs {
		if err = pool().submitChanIndex(key, ib, func(key string, ib IndexBack) {
			md5Key := gnomon.CryptoHash().MD516(key) // hash(keyStructure) 会发生碰撞，因此这里存储md5结果进行反向验证
			// 写入5位key及16位md5后key
			appendStr := strings.Join([]string{gnomon.String().PrefixSupplementZero(gnomon.Scale().Uint32ToDDuoString(ib.getHashKey()), 5), md5Key}, "")
			gnomon.Log().Debug("insert", gnomon.Log().Field("appendStr", appendStr), gnomon.Log().Field("formIndexFilePath", ib.getFormIndexFilePath()))
			// 将获取到的索引存储位置传入。如果为0，则表示没有存储过；如果不为0，则覆盖旧的存储记录
			// 写入5位key及16位md5后key及16位起始seek和8位持续seek
			wr := store().appendIndex(ib, appendStr, wf)
			if nil == wr.err {
				gnomon.Log().Debug("insert", gnomon.Log().Field("md5Key", md5Key), gnomon.Log().Field("seekStartIndex", wr.seekStartIndex))
				ib.getLink().setMD5Key(md5Key)
				ib.getLink().setSeekStart(wr.seekStart)
				ib.getLink().setSeekLast(wr.seekLast)
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

// rangeIndexes 遍历表索引ID集合，检索并计算所有索引返回对象集合
func (c *database) rangeIndexes(form Form, key string, autoID uint32, indexes map[string]Index, value interface{}, update bool) ([]IndexBack, error) {
	var (
		chanIndex chan IndexBack
		err       error
	)
	indexLen := len(indexes)
	chanIndex = make(chan IndexBack, indexLen) // 创建索引ID结果返回通道
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	for _, info := range indexes {
		if err = pool().submitIndexInfo(autoID, info, func(autoID uint32, index Index) {
			gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("index.id", index.getID()), gnomon.Log().Field("index.keyStructure", index.getKeyStructure()))
			if index.getKeyStructure() == indexAutoID {
				chanIndex <- form.getIndexes()[index.getID()].put(strconv.Itoa(int(autoID)), autoID, value, update)
			} else if index.getKeyStructure() == indexDefaultID {
				chanIndex <- form.getIndexes()[index.getID()].put(key, hash(key), value, update)
			} else {
				reflectObj := reflect.ValueOf(value) // 反射对象，通过reflectObj获取存储在里面的值，还可以去改变值
				params := strings.Split(index.getKeyStructure(), ".")
				var checkValue reflect.Value
				for _, param := range params {
					checkNewValue := reflectObj.Elem().FieldByName(param)
					if checkNewValue.IsValid() { // 子字段有效
						checkValue = checkNewValue
						continue
					}
					chanIndex <- &indexBack{err: errors.New(strings.Join([]string{"index", index.getKeyStructure(), "is invalid"}, " "))}
					return
				}
				gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("checkValue", checkValue))
				if keyNew, hashKeyNew, valid := valueType2index(&checkValue); valid {
					chanIndex <- form.getIndexes()[index.getID()].put(keyNew, hashKeyNew, value, update)
				}
			}
		}); nil != err {
			return nil, err
		}
	}
	var ibs []IndexBack
	for i := 0; i < indexLen; i++ {
		ib := <-chanIndex
		gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("ib.formIndexFilePath", ib.getFormIndexFilePath()))
		if err = ib.getErr(); nil != err {
			return nil, err
		}
		ibs = append(ibs, ib)
	}
	return ibs, nil
}

// shopperIsInvalid 自定义error信息
func shopperIsInvalid(formName string) error {
	return errors.New(strings.Join([]string{"invalid name ", formName}, ""))
}

// indexID 索引ID新的组合名称
func (c *database) indexID(formName, indexName string) string {
	return strings.Join([]string{formName, indexName}, "_")
}

// name2id 确保数据库唯一ID不重复
func (c *database) name2id(name string) string {
	id := gnomon.CryptoHash().MD516(name)
	have := true
	for have {
		have = false
		for _, v := range c.forms {
			if v.getID() == id {
				have = true
				id = gnomon.CryptoHash().MD516(strings.Join([]string{id, gnomon.String().RandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}
