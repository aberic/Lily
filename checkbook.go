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
	indexAutoID   = "lily_id"
	indexCustomID = "lily_custom"
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

func (c *checkbook) createForm(formName, comment, formType string) error {
	// 确定库名不重复
	for k := range c.forms {
		if k == formName {
			return ErrFormExist
		}
	}
	indexes := make(map[string]*index)
	// 确保表唯一ID不重复
	formID := c.name2id(formName)
	// 自增索引ID
	indexID := c.name2id(strings.Join([]string{formName, indexAutoID}, "_"))
	indexes[indexID] = &index{id: indexID, key: indexAutoID}
	fileIndex := 0
	if formType == formTypeSQL {
		if err := mkFormResourceSQL(c.id, formID, indexID, fileIndex); nil != err {
			return err
		}
	} else {
		// 默认自定义Key生成ID
		customID := c.name2id(strings.Join([]string{formName, indexCustomID}, "_"))
		if err := mkFormResource(c.id, formID, indexID, customID, fileIndex); nil != err {
			return err
		}
		indexes[customID] = &index{id: customID, key: indexCustomID}
	}
	c.forms[formName] = &shopper{
		autoID:    0,
		name:      formName,
		id:        formID,
		indexes:   indexes,
		fileIndex: fileIndex,
		comment:   comment,
		database:  c,
		nodes:     []Nodal{},
		formType:  formType,
	}
	return nil
}

func (c *checkbook) createIndex(formName string, key string, value interface{}) (uint32, error) {
	// todo
	return 0, nil
}

func (c *checkbook) put(formName string, key string, value interface{}, update bool) (uint32, error) {
	form := c.forms[formName] // 获取待操作表
	if nil == form {
		return 0, shopperIsInvalid(formName)
	}
	indexes := form.getIndexes()                    // 获取表索引ID集合
	autoID := atomic.AddUint32(form.getAutoID(), 1) // ID自增
	return c.insertDataWithIndexInfo(form, key, autoID, indexes, value, update)
}

func (c *checkbook) get(formName string, key string) (interface{}, error) {
	form := c.forms[formName]
	if nil == form {
		return nil, shopperIsInvalid(formName)
	}
	return form.get(key, hash(key))

}

func (c *checkbook) insert(formName string, value interface{}, update bool) (uint32, error) {
	// todo
	return 0, nil
}

func (c *checkbook) query(formName string, selector *Selector) (interface{}, error) {
	if nil == c {
		return nil, ErrDataIsNil
	}
	selector.formName = formName
	selector.checkbook = c
	return selector.query()
}

func (c *checkbook) valueTypeCheckKey(value *reflect.Value) (key string, hashKey uint32, support bool) {
	support = true
	switch value.Kind() {
	default:
		return "", 0, false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := value.Int()
		key = strconv.FormatInt(i64, 10)
		hashKey = int64ToUint32Index(i64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ui64 := value.Uint()
		key = gnomon.String().PrefixSupplementZero(gnomon.Scale().Uint64ToDDuoString(ui64), 5)
		hashKey = uint64ToUint32Index(ui64)
	case reflect.Float32, reflect.Float64:
		i64 := gnomon.Scale().Wrap(value.Float(), 4)
		key = strconv.FormatInt(i64, 10)
		hashKey = int64ToUint32Index(i64)
	case reflect.String:
		key = value.String()
		hashKey = hash(key)
	}
	return
}

func (c *checkbook) insertDataWithIndexInfo(form Form, key string, autoID uint32, indexes map[string]*index, value interface{}, update bool) (uint32, error) {
	var (
		chanIndex chan *indexBack
		err       error
	)
	indexLen := len(indexes)
	chanIndex = make(chan *indexBack, indexLen) // 创建索引ID结果返回通道
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	for _, info := range indexes {
		if err = pool().submitIndexInfo(autoID, info, func(autoID uint32, index *index) {
			if index.key == indexAutoID {
				chanIndex <- form.put(index.id, strconv.Itoa(int(autoID)), autoID, value, update)
			} else if index.key == indexCustomID {
				chanIndex <- form.put(index.id, key, hash(key), value, update)
			} else {
				reflectObj := reflect.ValueOf(value) // 反射对象，通过reflectObj获取存储在里面的值，还可以去改变值
				params := strings.Split(index.key, ".")
				checkValue := reflectObj.Elem()
				for _, param := range params {
					if checkValue = checkValue.FieldByName(param); checkValue.IsValid() { // 子字段有效
						continue
					}
					chanIndex <- &indexBack{err: errors.New(strings.Join([]string{"index", index.key, "is invalid"}, " "))}
					return
				}
				if keyNew, hashKeyNew, valid := c.valueTypeCheckKey(&checkValue); valid {
					chanIndex <- form.put(index.id, keyNew, hashKeyNew, value, update)
				}
			}
		}); nil != err {
			return 0, err
		}
	}
	var ibs []*indexBack
	for i := 0; i < indexLen; i++ {
		ib := <-chanIndex
		if nil != ib.err {
			return 0, ib.err
		}
		ibs = append(ibs, ib)
	}
	wrIndexBack := make(chan *writeResult, 1) // 索引存储结果通道
	// 存储数据到表文件
	wf := store().appendForm(form, pathFormDataFile(c.id, form.getID(), form.getFileIndex()), value)
	if nil != wf.err {
		return 0, wf.err
	}
	for _, ib := range ibs {
		if err = pool().submitChanIndex(ib, func(ib *indexBack) {
			md5Key := gnomon.CryptoHash().MD516(ib.originalKey) // hash(originalKey) 会发生碰撞，因此这里存储md5结果进行反向验证
			// 写入5位key及16位md5后key
			appendStr := strings.Join([]string{gnomon.String().PrefixSupplementZero(gnomon.Scale().Uint32ToDDuoString(ib.key), 5), md5Key}, "")
			gnomon.Log().Debug("insert", gnomon.LogField("appendStr", appendStr), gnomon.LogField("formIndexFilePath", ib.formIndexFilePath))
			// 写入5位key及16位md5后key及16位起始seek和8位持续seek
			wr := store().appendIndex(ib.indexNodal, ib.formIndexFilePath, appendStr, wf)
			if nil == wr.err {
				gnomon.Log().Debug("insert", gnomon.LogField("md5Key", md5Key))
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

// shopperIsInvalid 自定义error信息
func shopperIsInvalid(formName string) error {
	return errors.New(strings.Join([]string{"invalid name ", formName}, ""))
}

// indexID 索引ID新的组合名称
func (c *checkbook) indexID(formName, indexName string) string {
	return strings.Join([]string{formName, indexName}, "_")
}

// name2id 确保数据库唯一ID不重复
func (c *checkbook) name2id(name string) string {
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
