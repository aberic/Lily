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
	"strings"
)

// Selector 检索选择器
//
// 查询顺序 scope -> match -> conditions -> skip -> sort -> limit
type Selector struct {
	Conditions []*condition `json:"conditions"` // Conditions 条件查询
	Skip       uint32       `json:"skip"`       // Skip 结果集跳过数量
	Sort       *sort        `json:"sort"`       // Sort 排序方式
	Limit      uint32       `json:"limit"`      // Limit 结果集顺序数量
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

//func (s *Selector) match2String(inter interface{}) string {
//	switch inter.(type) {
//	case string:
//		return inter.(string)
//	case int:
//		return strconv.Itoa(inter.(int))
//	case float64:
//		return strconv.FormatFloat(inter.(float64), 'f', -1, 64)
//	case bool:
//		return strconv.FormatBool(inter.(bool))
//	}
//	return ""
//}

func (s *Selector) query() (int, []interface{}, error) {
	var (
		index     Index
		leftQuery bool
		ns        *nodeSelector
		count     int
		is        []interface{}
	)
	index, leftQuery, ns = s.getIndex()
	if nil == index {
		return 0, nil, errors.New("index not found")
	}
	gnomon.Log().Debug("query", gnomon.Log().Field("index", index.getKeyStructure()))
	if leftQuery {
		count, is = s.leftQueryIndex(index, ns)
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
//
// cond 条件对象
func (s *Selector) getIndex() (index Index, leftQuery bool, node *nodeSelector) {
	var idx Index
	idx, leftQuery, node = s.getIndexCondition()
	if idx != nil {
		return idx, leftQuery, node
	}
	for _, idx = range s.database.getForms()[s.formName].getIndexes() {
		if s.Sort != nil && s.Sort.Param == idx.getKeyStructure() {
			return idx, s.Sort.ASC, nil
		}
	}
	// 取值默认索引来进行查询操作
	for _, idx = range s.database.getForms()[s.formName].getIndexes() {
		gnomon.Log().Debug("getIndex", gnomon.Log().Field("index", index))
		return idx, true, nil
	}
	return nil, false, nil
}

func (s *Selector) getIndexCondition() (index Index, leftQuery bool, node *nodeSelector) {
	if len(s.Conditions) > 0 { // 优先尝试采用条件作为索引，缩小索引范围以提高检索效率
		for _, condition := range s.Conditions { // 遍历检索条件
			meetRight := false
			for _, idx := range s.database.getForms()[s.formName].getIndexes() {
				if condition.Param == idx.getKeyStructure() { // 匹配条件是否存在已有索引
					if nil != s.Sort && s.Sort.Param == idx.getKeyStructure() { // 如果有，则继续判断该索引是否存在排序需求
						index = idx
						leftQuery = s.Sort.ASC
						node = s.getConditionNode(condition)
						meetRight = true
						break
					}
					if index != nil {
						continue
					}
					index = idx
					node = s.getConditionNode(condition)
				}
			}
			if meetRight {
				break
			}
		}
	}
	return
}

// leftQueryIndex 索引顺序检索
//
// index 已获取索引对象
func (s *Selector) leftQueryIndex(index Index, ns *nodeSelector) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	for _, node := range index.getNode().getNodes() {
		nc, nis := s.leftQueryNode(node, ns)
		count += nc
		is = append(is, nis...)
	}
	if s.Skip > 0 {
		if s.Skip < uint32(len(is)) {
			is = is[s.Skip:]
		} else {
			is = is[0:0]
		}
	}
	if s.Limit > 0 && s.Limit < uint32(len(is)) {
		is = is[:s.Limit]
	}
	if s.Sort == nil {
		return count, is
	}
	return count, s.shellSort(is)
}

// condition 判断当前条件是否满足
func (s *Selector) condition(node Nodal, ns *nodeSelector) bool {
	if ns != nil {
		for _, cond := range s.Conditions {
			if cond != ns.cond {
				continue
			}
			switch cond.Cond {
			case "gt":
				return ns.degreeIndex > node.getDegreeIndex()
			case "lt":
				return ns.degreeIndex < node.getDegreeIndex()
			case "eq":
				return ns.degreeIndex == node.getDegreeIndex()
			case "dif":
				return ns.degreeIndex != node.getDegreeIndex()
			}
		}
	}
	return true
}

// leftQueryNode 节点顺序检索
func (s *Selector) leftQueryNode(node Nodal, ns *nodeSelector) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	if nodes := node.getNodes(); nil != nodes {
		for _, nd := range node.getNodes() {
			// 判断当前条件是否满足，如果满足则继续下一步
			if s.condition(nd, ns) {
				nc, nis := s.leftQueryNode(nd, ns)
				count += nc
				is = append(is, nis...)
			}
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
			// todo 条件确定
			//for _, cond := range s.Conditions {
			//
			//}
			is = append(is, inter)
		}
	}
	return len(leaf.getLinks()), is
}

// rightQueryIndex 索引倒序检索
//
// index 已获取索引对象
func (s *Selector) rightQueryIndex(index Index) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	lenNode := len(index.getNode().getNodes())
	for i := lenNode - 1; i >= 0; i-- {
		nc, nis := s.rightQueryNode(index.getNode().getNodes()[i])
		count += nc
		is = append(is, nis...)
	}
	if s.Skip > 0 {
		if s.Skip < uint32(len(is)) {
			is = is[s.Skip:]
		} else {
			is = is[0:0]
		}
	}
	if s.Limit > 0 && s.Limit < uint32(len(is)) {
		is = is[:s.Limit]
	}
	if s.Sort == nil {
		return count, is
	}
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
	gnomon.Log().Debug("shellSort 希尔排序", gnomon.Log().Field("s.Sort", s.Sort))
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

func (s *Selector) hashKeyFromValue(params []string, value interface{}) uint64 {
	hashKey, support := s.getInterValue(params, value)
	if !support {
		return 0
	}
	return hashKey
}

// getInterValue 根据索引描述和当前检索到的value对象获取当前value对象所在索引的hashKey
func (s *Selector) getInterValue(params []string, value interface{}) (hashKey uint64, support bool) {
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

func (s *Selector) getConditionNode(cond *condition) *nodeSelector {
	var (
		hashKey uint64
		ok      bool
	)
	if _, hashKey, ok = type2index(cond.Value); !ok {
		return nil
	}

	node1 := &nodeSelector{level: 1, degreeIndex: 0, cond: cond}
	nowKey := hashKey
	distance := levelDistance(node1.level)
	nextDegree := uint16(nowKey / distance)
	nextFlexibleKey := nowKey - uint64(nextDegree)*distance

	node2 := &nodeSelector{level: 2, degreeIndex: nextDegree, cond: cond}
	node1.nextNode = node2
	nowKey = nextFlexibleKey
	distance = levelDistance(node1.level)
	nextDegree = uint16(nowKey / distance)
	nextFlexibleKey = nowKey - uint64(nextDegree)*distance

	node3 := &nodeSelector{level: 3, degreeIndex: nextDegree, cond: cond}
	node2.nextNode = node3
	nowKey = nextFlexibleKey
	distance = levelDistance(node1.level)
	nextDegree = uint16(nowKey / distance)
	nextFlexibleKey = nowKey - uint64(nextDegree)*distance

	node4 := &nodeSelector{level: 4, degreeIndex: nextDegree, cond: cond}
	node3.nextNode = node4

	return node1
}

type nodeSelector struct {
	level       uint8  // 当前节点所在树层级
	degreeIndex uint16 // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	nextNode    *nodeSelector
	cond        *condition
}
