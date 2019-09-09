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
	"strconv"
)

// Selector 检索选择器
//
// 查询顺序 scope -> match -> conditions -> skip -> sort -> limit
type Selector struct {
	Scopes     []*scope     `json:"scopes"`     // Scopes 范围查询
	Matches    []*match     `json:"matches"`    // Matches 匹配查询
	Conditions []*condition `json:"conditions"` // Conditions 条件查询
	Skip       int32        `json:"skip"`       // Skip 结果集跳过数量
	Sort       *sort        `json:"sort"`       // Sort 排序方式
	Limit      int32        `json:"limit"`      // Limit 结果集顺序数量
	database   Database     // database 数据库对象
	formName   string       // formName 表名
}

// scope 范围查询
//
// 查询出来结果集合的起止位置
type scope struct {
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
	Start int64  `json:"Start"` // 起始位置
	End   int64  `json:"end"`   // 终止位置
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

// match 匹配查询
//
// 查询过程中匹配与指定参数值相等的结果
type match struct {
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
	Value interface{} `json:"value"` // 参数值 string/int/float64/bool
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

func (s *Selector) query() ([]interface{}, error) {
	var (
		index     Index
		leftQuery bool
		sortIndex bool
		err       error
	)
	if index, leftQuery, sortIndex, err = s.getIndex(); nil != err {
		return nil, err
	}
	gnomon.Log().Debug("query", gnomon.LogField("index", index.getKeyStructure()))
	if leftQuery {
		return s.leftQueryIndex(index, sortIndex), nil
	}
	return s.rightQueryIndex(index), nil
}

// getIndex 根据检索条件获取使用索引对象
//
// index 已获取索引对象
//
// leftQuery 是否顺序查询
func (s *Selector) getIndex() (index Index, leftQuery bool, sortIndex bool, err error) {
	for _, index = range s.database.getForms()[s.formName].getIndexes() {
		if len(s.Scopes) > 0 {
			for _, scope := range s.Scopes {
				if scope.Param == index.getKeyStructure() {
					return index, true, false, nil
				}
			}
		}
		if len(s.Matches) > 0 {
			for _, match := range s.Matches {
				if match.Param == index.getKeyStructure() {
					return index, true, false, nil
				}
			}
		}
		if len(s.Conditions) > 0 {
			for _, condition := range s.Conditions {
				if condition.Param == index.getKeyStructure() {
					return index, true, false, nil
				}
			}
		}
		if s.Sort != nil && s.Sort.Param == index.getKeyStructure() {
			return index, s.Sort.ASC, true, nil
		}
	}
	// 取值默认索引来进行查询操作
	for _, idx := range s.database.getForms()[s.formName].getIndexes() {
		return idx, true, false, nil
	}
	return nil, false, false, errors.New("index not found")
}

func (s *Selector) sort(is []interface{}) []interface{} {
	if s.Sort != nil {
		// todo 排序
		return is
	}
	return is
}

func (s *Selector) leftQueryIndex(index Index, sortIndex bool) []interface{} {
	is := make([]interface{}, 0)
	if nodes := index.getNodes(); nil != nodes {
		for _, node := range index.getNodes() {
			is = append(is, s.leftQueryNode(node)...)
		}
	}
	gnomon.Log().Debug("leftQueryIndex", gnomon.LogField("is", is))
	if sortIndex {
		return is
	}
	return s.sort(is)
}

func (s *Selector) leftQueryNode(node Nodal) []interface{} {
	is := make([]interface{}, 0)
	if nodes := node.getNodes(); nil != nodes {
		for _, node := range node.getNodes() {
			is = append(is, s.leftQueryNode(node)...)
		}
	} else {
		return s.leftQueryLeaf(node.(Leaf))
	}
	return is
}

func (s *Selector) leftQueryLeaf(leaf Leaf) []interface{} {
	is := make([]interface{}, 0)
	for _, link := range leaf.getLinks() {
		if inter, err := link.get(); nil == err {
			is = append(is, inter)
		}
	}
	return is
}

func (s *Selector) rightQueryIndex(index Index) []interface{} {
	is := make([]interface{}, 0)
	if nodes := index.getNodes(); nil != nodes {
		lenNode := len(nodes)
		for i := lenNode - 1; i >= 0; i-- {
			is = append(is, s.rightQueryNode(nodes[i])...)
		}
	}
	gnomon.Log().Debug("rightQueryIndex", gnomon.LogField("is", is))
	return is
}

func (s *Selector) rightQueryNode(node Nodal) []interface{} {
	is := make([]interface{}, 0)
	if nodes := node.getNodes(); nil != nodes {
		lenNode := len(nodes)
		for i := lenNode - 1; i >= 0; i-- {
			is = append(is, s.rightQueryNode(nodes[i])...)
		}
	} else {
		return s.rightQueryLeaf(node.(Leaf))
	}
	return is
}

func (s *Selector) rightQueryLeaf(leaf Leaf) []interface{} {
	is := make([]interface{}, 0)
	links := leaf.getLinks()
	lenLink := len(links)
	for i := lenLink - 1; i >= 0; i-- {
		if inter, err := links[i].get(); nil == err {
			is = append(is, inter)
		}
	}
	return is
}
