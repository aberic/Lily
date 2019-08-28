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

package Lily

import (
	"strconv"
)

// Selector 检索选择器
type Selector struct {
	Scope      *scope       `json:"scopes"`     // Scope 范围查询
	Conditions []*condition `json:"conditions"` // Conditions 条件查询
	Matches    []*match     `json:"matches"`    // Matches 匹配查询
	Skip       int32        `json:"skip"`       // Skip 结果集跳过数量
	Limit      int32        `json:"limit"`      // Limit 结果集顺序数量
	Sort       *sort        `json:"sort"`       // Sort 排序方式
}

// 索引，默认'_id'，可选自定义参数
//
// 索引将按照编号顺序拼接参数名字符串组合成唯一键
type index struct {
	param string // 参数名
	order uint8  // 编号
}

// indexes 索引对象，封装一个可排序的索引数组
type indexes struct {
	IndexArr []*index `json:"indexArr"`
}

func (i indexes) Len() int           { return len(i.IndexArr) }
func (i indexes) Swap(m, n int)      { i.IndexArr[m], i.IndexArr[n] = i.IndexArr[n], i.IndexArr[m] }
func (i indexes) Less(m, n int) bool { return i.IndexArr[m].order < i.IndexArr[n].order }

// scope 范围查询
//
// 查询出来结果集合的起止位置
type scope struct {
	Start int32 `json:"start"` // 起始位置
	End   int32 `json:"end"`   // 终止位置
}

// condition 条件查询
//
// 查询过程中不满足条件的记录将被移除出结果集
type condition struct {
	Param string `json:"param"` // 参数名
	Cond  string `json:"cond"`  // 条件 gt/lt/eq/dif 大于/小于/等于/不等
}

// match 匹配查询
//
// 查询过程中匹配与指定参数值相等的结果
type match struct {
	Param string      `json:"param"` // 参数名
	Value interface{} `json:"value"` // 参数值 string/int/float64/bool
}

// sort 排序方式
type sort struct {
	Indexes *indexes `json:"indexes"` // Indexes 索引对象，封装一个可排序的索引数组
	ASC     bool     `json:"asc"`     // 是否升序
}

//func (s *Selector) lilyName(lilyName string) string {
//
//}

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

func (s *Selector) query(node nodal, asc bool) []interface{} {
	if asc {
		return s.leftQuery(node)
	}
	return s.rightQuery(node)
}

func (s *Selector) leftQuery(node nodal) []interface{} {
	is := make([]interface{}, 0)
	count := node.childCount()
	for i := 0; i < count; i++ {
		child := node.child(i)
		if child.childCount() > 0 {
			is = append(is, s.leftQuery(child)...)
		} else if child.childCount() == -1 {
			thg := child.(*box).things
			lenThg := len(thg)
			for ti := 0; ti < lenThg; ti++ {
				is = append(is, thg[ti].value)
			}
		}
	}
	return is
}

func (s *Selector) rightQuery(node nodal) []interface{} {
	is := make([]interface{}, 0)
	count := node.childCount()
	for i := count - 1; i >= 0; i-- {
		child := node.child(i)
		if child.childCount() > 0 {
			is = append(is, s.rightQuery(child)...)
		} else if child.childCount() == -1 {
			thg := child.(*box).things
			lenThg := len(thg)
			for ti := lenThg - 1; ti >= 0; ti-- {
				is = append(is, thg[ti].value)
			}
		}
	}
	return is
}
