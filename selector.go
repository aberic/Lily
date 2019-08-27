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

type Selector struct {
	Indexes    []*index     `json:"index"`
	Scopes     *scope       `json:"scopes"`
	Conditions []*condition `json:"conditions"`
	Matches    []*match     `json:"matches"`
	Skip       int32        `json:"skip"`
	Limit      int32        `json:"limit"`
	Sort       *sort        `json:"sort"`
}

// 索引，默认'_id'，可选自定义参数
type index struct {
	param string
	sort  uint8
}

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
	Param string `json:"param"` // 参数名 默认 _id
	ASC   bool   `json:"asc"`   // 是否升序
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

func (s *Selector) query(node nodal) []interface{} {
	if nil != s.Sort && s.Sort.Param == "_id" && !s.Sort.ASC {
		return s.rightQuery(node)
	}
	return s.leftQuery(node)
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
