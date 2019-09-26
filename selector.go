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
)

// Selector 检索选择器
//
// 查询顺序 scope -> match -> conditions -> skip -> sort -> limit
type Selector struct {
	Conditions []*condition `json:"conditions"` // Conditions 条件查询
	Skip       int32        `json:"skip"`       // Skip 结果集跳过数量
	Sort       *sort        `json:"sort"`       // Sort 排序方式
	Limit      int32        `json:"limit"`      // Limit 结果集顺序数量
	database   Database     // database 数据库对象
	formName   string       // formName 表名
}

// condition 条件查询
//
// 查询过程中不满足条件的记录将被移除出结果集
type condition struct {
	// 参数名，由对象结构层级字段通过'.'组成，如
	//
	// ref := &ref{
	//		i: 1,
	//		s: "2",
	//		in: refIn{
	//			i: 3,
	//			s: "4",
	//		},
	//	}
	//
	// key可取'i','in.s'
	Param string      `json:"param"`
	Cond  string      `json:"cond"`  // 条件 gt/lt/eq/dif 大于/小于/等于/不等
	Value interface{} `json:"value"` // 比较对象，支持int、string、float和bool
}

// sort 排序方式
type sort struct {
	// 参数名，由对象结构层级字段通过'.'组成，如
	//
	// ref := &ref{
	//		i: 1,
	//		s: "2",
	//		in: refIn{
	//			i: 3,
	//			s: "4",
	//		},
	//	}
	//
	// key可取'i','in.s'
	Param string `json:"param"`
	ASC   bool   `json:"asc"` // 是否升序
}

func (s *Selector) match2String(inter interface{}) string {
	switch inter.(type) {
	case string:
		return inter.(string)
	case int:
		return strconv.Itoa(inter.(int))
	case float64:
		return strconv.FormatFloat(inter.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(inter.(bool))
	}
	return ""
}

func (s *Selector) query() (int, []interface{}, error) {
	var (
		index     Index
		leftQuery bool
		count     int
		is        []interface{}
		err       error
	)
	if index, leftQuery, _, err = s.getIndex(); nil != err {
		return 0, nil, err
	}
	gnomon.Log().Debug("query", gnomon.Log().Field("index", index.getKeyStructure()))
	if leftQuery {
		count, is = s.leftQueryIndex(index)
	} else {
		count, is = s.rightQueryIndex(index)
	}
	return count, is, nil
}

// getIndex 根据检索条件获取使用索引对象
//
// index 已获取索引对象
//
// leftQuery 是否顺序查询
func (s *Selector) getIndex() (index Index, leftQuery bool, sortIndex bool, err error) {
	if len(s.Conditions) > 0 {
		for _, condition := range s.Conditions {
			for _, index = range s.database.getForms()[s.formName].getIndexes() {
				if condition.Param == index.getKeyStructure() {
					return index, true, false, nil
				}
			}
		}
	}
	for _, index = range s.database.getForms()[s.formName].getIndexes() {
		if s.Sort != nil && s.Sort.Param == index.getKeyStructure() {
			return index, s.Sort.ASC, true, nil
		}
	}
	// 取值默认索引来进行查询操作
	for _, idx := range s.database.getForms()[s.formName].getIndexes() {
		gnomon.Log().Debug("getIndex", gnomon.Log().Field("index", index))
		return idx, true, false, nil
	}
	return nil, false, false, errors.New("index not found")
}

// leftQueryIndex 索引顺序检索
func (s *Selector) leftQueryIndex(index Index) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	if nodes := index.getNodes(); nil != nodes {
		for _, node := range index.getNodes() {
			nc, nis := s.leftQueryNode(node)
			count += nc
			is = append(is, nis...)
		}
	}
	//gnomon.Log().Debug("leftQueryIndex", gnomon.Log().Field("is", is))
	if s.Sort == nil {
		gnomon.Log().Debug("leftQueryIndex", gnomon.Log().Field("s.Sort", s.Sort))
		return count, is
	}
	return count, s.shellSort(is)
}

// leftQueryNode 节点顺序检索
func (s *Selector) leftQueryNode(node Nodal) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	if nodes := node.getNodes(); nil != nodes {
		for _, node := range node.getNodes() {
			nc, nis := s.leftQueryNode(node)
			count += nc
			is = append(is, nis...)
		}
	} else {
		return s.leftQueryLeaf(node.(Leaf))
	}
	return count, is
}

// leftQueryLeaf 叶子节点顺序检索
func (s *Selector) leftQueryLeaf(leaf Leaf) (int, []interface{}) {
	is := make([]interface{}, 0)
	for _, link := range leaf.getLinks() {
		if inter, err := link.get(); nil == err {
			is = append(is, inter)
		}
	}
	return len(leaf.getLinks()), is
}

// rightQueryIndex 索引倒序检索
func (s *Selector) rightQueryIndex(index Index) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	if nodes := index.getNodes(); nil != nodes {
		lenNode := len(nodes)
		for i := lenNode - 1; i >= 0; i-- {
			nc, nis := s.rightQueryNode(nodes[i])
			count += nc
			is = append(is, nis...)
		}
	}
	gnomon.Log().Debug("rightQueryIndex", gnomon.Log().Field("is", is))
	return count, s.shellSort(is)
}

// rightQueryNode 节点倒序检索
func (s *Selector) rightQueryNode(node Nodal) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	if nodes := node.getNodes(); nil != nodes {
		lenNode := len(nodes)
		for i := lenNode - 1; i >= 0; i-- {
			nc, nis := s.rightQueryNode(nodes[i])
			count += nc
			is = append(is, nis...)
		}
	} else {
		return s.rightQueryLeaf(node.(Leaf))
	}
	return count, is
}

// rightQueryLeaf 叶子节点倒序检索
func (s *Selector) rightQueryLeaf(leaf Leaf) (int, []interface{}) {
	is := make([]interface{}, 0)
	links := leaf.getLinks()
	lenLink := len(links)
	for i := lenLink - 1; i >= 0; i-- {
		if inter, err := links[i].get(); nil == err {
			is = append(is, inter)
		}
	}
	return len(leaf.getLinks()), is
}

// shellSort 希尔排序
func (s *Selector) shellSort(is []interface{}) []interface{} {
	gnomon.Log().Debug("shellSort 希尔排序")
	if s.Sort.ASC {
		gnomon.Log().Debug("shellAsc 希尔顺序排序")
		return s.shellAsc(is)
	}
	gnomon.Log().Debug("shellDesc 希尔倒序排序")
	return s.shellDesc(is)
}

// shellAsc 希尔顺序排序
func (s *Selector) shellAsc(is []interface{}) []interface{} {
	length := len(is)
	gap := length / 2
	params := strings.Split(s.Sort.Param, ".")
	for gap > 0 {
		for i := gap; i < length; i++ {
			tempI := is[i]
			temp := s.hashKeyFromValue(params, is[i])
			preIndex := i - gap
			for preIndex >= 0 && s.hashKeyFromValue(params, is[preIndex]) > temp {
				is[preIndex+gap] = is[preIndex]
				preIndex -= gap
			}
			is[preIndex+gap] = tempI
		}
		gap /= 2
	}
	return is
}

// shellDesc 希尔倒序排序
func (s *Selector) shellDesc(is []interface{}) []interface{} {
	length := len(is)
	gap := length / 2
	params := strings.Split(s.Sort.Param, ".")
	for gap > 0 {
		for i := gap; i < length; i++ {
			tempI := is[i]
			temp := s.hashKeyFromValue(params, is[i])
			preIndex := i - gap
			for preIndex >= 0 && s.hashKeyFromValue(params, is[preIndex]) < temp {
				is[preIndex+gap] = is[preIndex]
				preIndex -= gap
			}
			is[preIndex+gap] = tempI
		}
		gap /= 2
	}
	return is
}

func (s *Selector) hashKeyFromValue(params []string, value interface{}) int64 {
	hashKey, support := s.getInterValue(params, value)
	if !support {
		return 0
	}
	return hashKey
}

// getInterValue 根据索引描述和当前检索到的value对象获取当前value对象所在索引的hashKey
func (s *Selector) getInterValue(params []string, value interface{}) (hashKey int64, support bool) {
	reflectObj := reflect.ValueOf(value) // 反射对象，通过reflectObj获取存储在里面的值，还可以去改变值
	if reflectObj.Kind() == reflect.Map {
		interMap := value.(map[string]interface{})
		lenParams := len(params)
		var valueResult interface{}
		for i, param := range params {
			if i == lenParams-1 {
				valueResult = interMap[param]
				break
			}
			interMap = interMap[param].(map[string]interface{})
		}
		//gnomon.Log().Debug("getInterValue", gnomon.Log().Field("valueResult", valueResult))
		checkValue := reflect.ValueOf(valueResult)
		return value2hashKey(&checkValue)
	}
	gnomon.Log().Debug("getInterValue", gnomon.Log().Field("kind", reflectObj.Kind()), gnomon.Log().Field("support", false))
	return 0, false
}
