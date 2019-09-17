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

import (
	"errors"
	"github.com/aberic/gnomon"
	"reflect"
	"strconv"
	"strings"
	"sync"
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

func (d *database) getID() string {
	return d.id
}

func (d *database) getName() string {
	return d.name
}

// getForms 获取数据库表集合
func (d *database) getForms() map[string]Form {
	return d.forms
}

func (d *database) createForm(formName, comment, formType string) error {
	// 确定库名不重复
	for k := range d.forms {
		if k == formName {
			return ErrFormExist
		}
	}
	// 确保表唯一ID不重复
	formID := d.name2id(formName)
	form := &form{
		autoID:   0,
		name:     formName,
		id:       formID,
		comment:  comment,
		database: d,
		indexes:  map[string]Index{},
		nodes:    []Nodal{},
		formType: formType,
	}
	err := mkFormResource(d.id, formID, 0)
	if nil != err {
		return err
	}
	d.forms[formName] = form
	if formType == FormTypeSQL {
		// 自增索引ID
		if err = d.createKey(formName, indexAutoID); nil != err {
			return err
		}
	} else {
		// 默认自定义Key生成ID
		if err = d.createKey(formName, indexDefaultID); nil != err {
			return err
		}
	}
	return nil
}

func (d *database) createKey(formName string, keyStructure string) error {
	form := d.forms[formName]
	// 自定义Key生成ID
	customID := d.name2id(strings.Join([]string{formName, keyStructure}, "_"))
	//gnomon.Log().Debug("createIndex", gnomon.Log().Field("customID", customID))
	if err := mkFormIndexResource(d.id, form.getID(), customID); nil != err {
		return err
	}
	form.getIndexes()[customID] = &index{id: customID, primary: true, keyStructure: keyStructure, form: form, fileIndex: 0}
	return nil
}

func (d *database) createIndex(formName string, keyStructure string) error {
	form := d.forms[formName]
	// 自定义Key生成ID
	customID := d.name2id(strings.Join([]string{formName, keyStructure}, "_"))
	//gnomon.Log().Debug("createIndex", gnomon.Log().Field("customID", customID))
	if err := mkFormIndexResource(d.id, form.getID(), customID); nil != err {
		return err
	}
	form.getIndexes()[customID] = &index{id: customID, primary: false, keyStructure: keyStructure, form: form, fileIndex: 0}
	return nil
}

func (d *database) put(formName string, key string, value interface{}, update bool) (uint32, error) {
	form := d.forms[formName] // 获取待操作表
	if nil == form {
		return 0, shopperIsInvalid(formName)
	}
	if form.getFormType() != FormTypeDoc {
		return 0, errors.New("put method only support doc")
	}
	indexes := form.getIndexes()                    // 获取表索引ID集合
	autoID := atomic.AddUint32(form.getAutoID(), 1) // ID自增
	return d.insertDataWithIndexInfo(form, key, autoID, indexes, value, update)
}

func (d *database) get(formName string, key string) (interface{}, error) {
	form := d.forms[formName]
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

func (d *database) insert(formName string, value interface{}, update bool) (uint32, error) {
	// todo
	return 0, nil
}

func (d *database) query(formName string, selector *Selector) (int, interface{}, error) {
	if nil == d {
		return 0, nil, ErrDataIsNil
	}
	selector.formName = formName
	selector.database = d
	return selector.query()
}

func (d *database) insertDataWithIndexInfo(form Form, key string, autoID uint32, indexes map[string]Index, value interface{}, update bool) (uint32, error) {
	var (
		ibs []IndexBack
		err error
	)
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	if ibs, err = d.rangeIndexes(form, key, autoID, indexes, value, update); nil != err {
		return 0, err
	}
	//gnomon.Log().Debug("insertDataWithIndexInfo", gnomon.Log().Field("ibs", ibs))
	wrIndexBack := make(chan *writeResult, 1) // 索引存储结果通道
	// 存储数据到表文件
	wrf := store().storeData(form, pathFormDataFile(d.id, form.getID(), form.getFileIndex()), value)
	if nil != wrf.err {
		return 0, wrf.err
	}
	for _, ib := range ibs {
		go func(key string, ib IndexBack) {
			wrIndexBack <- store().storeIndex(ib, wrf)
		}(key, ib)
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
func (d *database) rangeIndexes(form Form, key string, autoID uint32, indexes map[string]Index, value interface{}, update bool) ([]IndexBack, error) {
	var (
		wg        sync.WaitGroup
		chanIndex chan IndexBack
		err       error
	)
	indexLen := len(indexes)
	chanIndex = make(chan IndexBack, indexLen) // 创建索引ID结果返回通道
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	for _, index := range indexes {
		wg.Add(1)
		go func(autoID uint32, index Index) {
			defer wg.Done()
			//gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("index.id", index.getID()), gnomon.Log().Field("index.keyStructure", index.getKeyStructure()))
			if index.getKeyStructure() == indexAutoID {
				chanIndex <- form.getIndexes()[index.getID()].put(strconv.Itoa(int(autoID)), autoID, update)
			} else if index.getKeyStructure() == indexDefaultID {
				chanIndex <- form.getIndexes()[index.getID()].put(key, hash(key), update)
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
				//gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("checkValue", checkValue))
				if keyNew, hashKeyNew, valid := valueType2index(&checkValue); valid {
					chanIndex <- form.getIndexes()[index.getID()].put(keyNew, hashKeyNew, update)
				}
			}
		}(autoID, index)
	}
	wg.Wait()
	var ibs []IndexBack
	for i := 0; i < indexLen; i++ {
		ib := <-chanIndex
		//gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("ib.formIndexFilePath", ib.getFormIndexFilePath()))
		if err = ib.getErr(); nil != err {
			//gnomon.Log().Debug("rangeIndexes", gnomon.Log().Err(err))
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
func (d *database) indexID(formName, indexName string) string {
	return strings.Join([]string{formName, indexName}, "_")
}

// name2id 确保数据库唯一ID不重复
func (d *database) name2id(name string) string {
	id := gnomon.CryptoHash().MD516(name)
	have := true
	for have {
		have = false
		for _, v := range d.forms {
			if v.getID() == id {
				have = true
				id = gnomon.CryptoHash().MD516(strings.Join([]string{id, gnomon.String().RandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}
