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
	"strconv"
)

// Selector 检索选择器
//
// 查询顺序 scope -> conditions -> match -> skip -> sort -> limit
type Selector struct {
	Scope      []*scope     `json:"scopes"`     // Scope 范围查询
	Conditions []*condition `json:"conditions"` // Conditions 条件查询
	Matches    []*match     `json:"matches"`    // Matches 匹配查询
	Skip       int32        `json:"skip"`       // Skip 结果集跳过数量
	Sort       *sort        `json:"sort"`       // Sort 排序方式
	Limit      int32        `json:"limit"`      // Limit 结果集顺序数量
	checkbook  *checkbook   // 数据库对象
	formName   string       // 表名
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
	Start int32  `json:"Start"` // 起始位置
	End   int32  `json:"end"`   // 终止位置
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

//func (s *Selector) formName(formName string) string {
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

func (s *Selector) query() ([]interface{}, error) {
	if len(s.Scope) == 0 && len(s.Conditions) == 0 && len(s.Matches) == 0 && s.Sort == nil {
		if l := s.checkbook.forms[s.formName]; nil != l {
			// todo skip & limit 限定
			return s.leftQuery(l), nil
		}
		return nil, shopperIsInvalid(s.formName)
	}
	// todo 条件全开检索
	return s.rightQuery(s.checkbook.forms[s.formName]), nil
}

func (s *Selector) leftQuery(data Data) []interface{} {
	is := make([]interface{}, 0)
	count := data.childCount()
	for i := 0; i < count; i++ {
		child := data.child(i)
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

func (s *Selector) rightQuery(data Data) []interface{} {
	is := make([]interface{}, 0)
	count := data.childCount()
	for i := count - 1; i >= 0; i-- {
		child := data.child(i)
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
