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
	"github.com/aberic/lily/api"
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
	id      string          // 数据库唯一ID，不能改变
	name    string          // 数据库名称，根据需求可以随时变化
	comment string          // 描述
	forms   map[string]Form // 表集合
	lily    *Lily           // 数据库引擎
}

func (d *database) getID() string {
	return d.id
}

func (d *database) getName() string {
	return d.name
}

func (d *database) getComment() string {
	return d.comment
}

// getForms 获取数据库表集合
func (d *database) getForms() map[string]Form {
	return d.forms
}

func (d *database) createDoc(formName, comment string) error {
	if err := d.createForm(formName, comment, FormTypeDoc); nil != err {
		return err
	}
	// 默认自定义Key生成ID
	_ = d.createKey(formName, indexDefaultID)
	d.lily.lilyData.Databases[d.name].Forms[formName].FormType = api.FormType_Doc
	return nil
}

func (d *database) createSQL(formName, comment string) error {
	if err := d.createForm(formName, comment, FormTypeSQL); nil != err {
		return err
	}
	// 自增索引ID
	_ = d.createKey(formName, indexAutoID)
	d.lily.lilyData.Databases[d.name].Forms[formName].FormType = api.FormType_SQL
	return nil
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
		formType: formType,
	}
	err := mkFormResource(d.id, formID)
	if nil != err {
		return err
	}
	d.forms[formName] = form
	// 同步数据到 pb.Lily
	d.lily.lilyData.Databases[d.name].Forms[formName] = &api.Form{
		Id:      formID,
		Name:    formName,
		Comment: comment,
		Indexes: map[string]*api.Index{},
	}
	return nil
}

func (d *database) createKey(formName string, keyStructure string) error {
	// 确定key名不重复
	for _, v := range d.forms[formName].getIndexes() {
		if v.getKeyStructure() == keyStructure {
			return ErrKeyExist
		}
	}
	form := d.forms[formName]
	// 自定义Key生成ID
	customID := d.name2id(strings.Join([]string{formName, keyStructure}, "_"))
	//gnomon.Log().Debug("createIndex", gnomon.Log().Field("customID", customID))
	index := &index{id: customID, primary: true, keyStructure: keyStructure, form: form}
	node := &node{level: 1, degreeIndex: 0, preNode: nil, nodes: []Nodal{}, index: index}
	index.node = node
	form.getIndexes()[customID] = index
	// 同步数据到 pb.Lily
	d.lily.lilyData.Databases[d.name].Forms[formName].Indexes[customID] = &api.Index{
		Id:           customID,
		Primary:      true,
		KeyStructure: keyStructure,
	}
	return nil
}

func (d *database) createIndex(formName string, keyStructure string) error {
	// 确定index名不重复
	for _, v := range d.forms[formName].getIndexes() {
		if v.getKeyStructure() == keyStructure {
			return ErrIndexExist
		}
	}
	form := d.forms[formName]
	// 自定义Key生成ID
	customID := d.name2id(strings.Join([]string{formName, keyStructure}, "_"))
	//gnomon.Log().Debug("createIndex", gnomon.Log().Field("customID", customID))
	index := &index{id: customID, primary: false, keyStructure: keyStructure, form: form}
	node := &node{level: 1, degreeIndex: 0, preNode: nil, nodes: []Nodal{}, index: index}
	index.node = node
	form.getIndexes()[customID] = index
	// 同步数据到 pb.Lily
	d.lily.lilyData.Databases[d.name].Forms[formName].Indexes[customID] = &api.Index{
		Id:           customID,
		Primary:      false,
		KeyStructure: keyStructure,
	}
	return nil
}

func (d *database) put(formName string, key string, value interface{}, update bool) (uint64, error) {
	form := d.forms[formName] // 获取待操作表
	if nil == form {
		return 0, formIsInvalid(formName)
	}
	indexes := form.getIndexes() // 获取表索引ID集合
	return d.insertDataWithIndexInfo(form, key, indexes, value, update, true)
}

func (d *database) get(formName string, key string) (interface{}, error) {
	form := d.forms[formName]
	if nil == form {
		return nil, formIsInvalid(formName)
	}
	for _, index := range form.getIndexes() {
		if index.getKeyStructure() == indexDefaultID {
			v, err := index.get(key, hash(key))
			if nil != err {
				return nil, err
			}
			switch v.(type) {
			default:
				return index.get(key, hash(key))
			case string:
				if gnomon.String().IsEmpty(v.(string)) {
					return nil, errors.New("value is invalid")
				}
			}
		}
	}
	return nil, errors.New("no key for custom id index")
}

func (d *database) remove(formName string, key string) error {
	value, err := d.get(formName, key)
	if nil != err {
		return err
	}
	form := d.forms[formName] // 获取待操作表
	if nil == form {
		return formIsInvalid(formName)
	}
	indexes := form.getIndexes() // 获取表索引ID集合
	_, err = d.insertDataWithIndexInfo(form, key, indexes, value, true, false)
	return err
}

func (d *database) delete(formName string, selector *Selector) error {
	// todo 删除
	return nil
}

func (d *database) query(formName string, selector *Selector) (int32, interface{}, error) {
	if nil == d {
		return 0, nil, ErrDataIsNil
	}
	selector.formName = formName
	selector.database = d
	return selector.query()
}

func (d *database) insertDataWithIndexInfo(form Form, key string, indexes map[string]Index, value interface{}, update, valid bool) (uint64, error) {
	var (
		ibs []IndexBack
		wg  sync.WaitGroup
		err error
	)
	//gnomon.Log().Debug("insertDataWithIndexInfo", gnomon.Log().Field("ibs", ibs))
	defer form.unLock()
	form.lock()
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	ibs = d.rangeIndexes(form, key, indexes, value, update)
	// 存储数据到表文件
	dataWriteResult := store().storeData(pathFormDataFile(d.id, form.getID()), value, valid)
	if nil != dataWriteResult.err {
		return 0, dataWriteResult.err
	}
	errBack := make(chan error, len(ibs)) // 索引存储结果通道
	for _, ib := range ibs {
		wg.Add(1)
		go func(key string, ib IndexBack) {
			defer wg.Done()
			wrIndexBack := store().storeIndex(ib, dataWriteResult)
			if wrIndexBack.err != nil {
				errBack <- wrIndexBack.err
			}
		}(key, ib)
	}
	wg.Wait()
	if len(errBack) > 0 {
		if err = <-errBack; nil != err {
			// todo 回滚策略待完成
			return 0, err
		}
	}
	return *form.getAutoID(), nil
}

// rangeIndexes 遍历表索引ID集合，检索并计算所有索引返回对象集合
func (d *database) rangeIndexes(form Form, key string, indexes map[string]Index, value interface{}, update bool) []IndexBack {
	var (
		wg        sync.WaitGroup
		chanIndex chan IndexBack
	)
	indexLen := len(indexes)
	chanIndex = make(chan IndexBack, indexLen) // 创建索引ID结果返回通道
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	for _, index := range indexes {
		wg.Add(1)
		go func(index Index) {
			defer wg.Done()
			//gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("index.id", index.getID()), gnomon.Log().Field("index.keyStructure", index.getKeyStructure()))
			if index.getKeyStructure() == indexAutoID {
				autoID := atomic.AddUint64(form.getAutoID(), 1) // ID自增
				chanIndex <- form.getIndexes()[index.getID()].put(strconv.FormatUint(autoID, 10), autoID, update)
			} else if index.getKeyStructure() == indexDefaultID {
				chanIndex <- form.getIndexes()[index.getID()].put(key, hash(key), update)
			} else {
				chanIndex <- d.getCustomIndex(form, index, value, update)
			}
		}(index)
	}
	wg.Wait()
	var ibs []IndexBack
	for i := 0; i < indexLen; i++ {
		ib := <-chanIndex
		if ib.getErr() == nil {
			ibs = append(ibs, ib)
		}
	}
	return ibs
}

// getCustomIndex 获取自定义索引预插入返回对象
func (d *database) getCustomIndex(form Form, idx Index, value interface{}, update bool) IndexBack {
	reflectValue := reflect.ValueOf(value) // 反射对象，通过reflectObj获取存储在里面的值，还可以去改变值
	params := strings.Split(idx.getKeyStructure(), ".")
	switch reflectValue.Kind() {
	default:
		return &indexBack{err: errors.New(strings.Join([]string{"index", idx.getKeyStructure(), "with type is invalid"}, " "))}
	case reflect.Map:
		var (
			item      interface{}
			paramsLen = len(params)
			position  int
			itemMap   = value.(map[string]interface{})
		)
		for _, param := range params {
			position++
			item = itemMap[param]
			if position == paramsLen { // 表示没有后续参数
				break
			}
			switch item := item.(type) {
			default:
				return &indexBack{err: errors.New(strings.Join([]string{"index", idx.getKeyStructure(), "with map is invalid"}, " "))}
			case map[string]interface{}:
				itemMap = item
				continue
			}
		}
		if keyNew, hashKeyNew, valid := type2index(item); valid {
			return form.getIndexes()[idx.getID()].put(keyNew, hashKeyNew, update)
		}
		return &indexBack{err: errors.New(strings.Join([]string{"index", idx.getKeyStructure(), "with map value is invalid"}, " "))}
	case reflect.Ptr:
		checkValue := reflectValue
		for _, param := range params {
			checkNewValue := checkValue.Elem().FieldByName(param)
			if checkNewValue.IsValid() { // 子字段有效
				checkValue = checkNewValue
				continue
			}
			return &indexBack{err: errors.New(strings.Join([]string{"index", idx.getKeyStructure(), "with ptr is invalid"}, " "))}
		}
		if keyNew, hashKeyNew, valid := valueType2index(&checkValue); valid {
			return form.getIndexes()[idx.getID()].put(keyNew, hashKeyNew, update)
		}
		return &indexBack{err: errors.New(strings.Join([]string{"index", idx.getKeyStructure(), "with ptr value is invalid"}, " "))}
	}
}

// formIsInvalid 自定义error信息
func formIsInvalid(formName string) error {
	return errors.New(strings.Join([]string{"invalid name ", formName}, ""))
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
