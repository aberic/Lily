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
		nc        *nodeCondition
		count     int
		is        []interface{}
		err       error
	)
	if index, leftQuery, nc, err = s.getIndex(); nil != err {
		return 0, nil, err
	}
	gnomon.Log().Debug("query", gnomon.Log().Field("index", index.getKeyStructure()))
	if leftQuery {
		count, is = s.leftQueryIndex(index, nc)
	} else {
		count, is = s.rightQueryIndex(index, nc)
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
func (s *Selector) getIndex() (index Index, leftQuery bool, node *nodeCondition, err error) {
	var idx Index
	// 优先尝试采用条件作为索引，缩小索引范围以提高检索效率
	idx, leftQuery, node, err = s.getIndexCondition()
	if idx != nil { // 如果存在条件查询，则优先条件查询
		return idx, leftQuery, node, err
	}
	for _, idx := range s.database.getForms()[s.formName].getIndexes() { // 如果存在排序查询，则优先排序查询
		if s.Sort != nil && s.Sort.Param == idx.getKeyStructure() {
			return idx, s.Sort.ASC, nil, nil
		}
	}
	// 取值默认索引来进行查询操作
	for _, idx := range s.database.getForms()[s.formName].getIndexes() {
		gnomon.Log().Debug("getIndex", gnomon.Log().Field("index", index))
		return idx, true, nil, nil
	}
	return nil, false, nil, errors.New("index not found")
}

// getIndexCondition 优先尝试采用条件作为索引，缩小索引范围以提高检索效率
//
// 优先匹配有多个相同Param参数的条件，如果相同数量一样，则按照先后顺序选择最先匹配的
func (s *Selector) getIndexCondition() (index Index, leftQuery bool, nc *nodeCondition, err error) {
	ncs := make(map[string]*nodeCondition)
	leftQuery = true
	for _, condition := range s.Conditions { // 遍历检索条件
		for _, idx := range s.database.getForms()[s.formName].getIndexes() {
			if condition.Param == idx.getKeyStructure() { // 匹配条件是否存在已有索引
				if nil != s.Sort && s.Sort.Param == idx.getKeyStructure() { // 如果有，则继续判断该索引是否存在排序需求
					index = idx
					leftQuery = s.Sort.ASC
					if ncs[condition.Param] == nil {
						ncs[condition.Param] = &nodeCondition{nss: []*nodeSelector{}}
					}
					s.getConditionNode(ncs[condition.Param], condition)
					break
				}
				if index != nil {
					continue
				}
				index = idx
				if ncs[condition.Param] == nil {
					ncs[condition.Param] = &nodeCondition{nss: []*nodeSelector{}}
				}
				s.getConditionNode(ncs[condition.Param], condition)
			}
		}
	}
	var nodeCount = 0
	for _, ncNow := range ncs {
		ncNowCount := len(ncNow.nss)
		if ncNowCount > nodeCount {
			nodeCount = ncNowCount
			nc = ncNow
		}
	}
	return
}

// leftQueryIndex 索引顺序检索
//
// index 已获取索引对象
func (s *Selector) leftQueryIndex(index Index, ns *nodeCondition) (int, []interface{}) {
	var (
		count, nc int
		nis       []interface{}
		is        = make([]interface{}, 0)
	)
	for _, node := range index.getNode().getNodes() {
		if nil == ns {
			nc, nis = s.leftQueryNode(node, nil)
		} else {
			nc, nis = s.leftQueryNode(node, ns.nextNode)
		}

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

// leftQueryNode 节点顺序检索
func (s *Selector) leftQueryNode(node Nodal, ns *nodeCondition) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	if nodes := node.getNodes(); nil != nodes {
		for _, nd := range node.getNodes() {
			if ns == nil {
				nc, nis := s.leftQueryNode(nd, nil)
				count += nc
				is = append(is, nis...)
			} else if s.conditions(nd, ns.nextNode.nss) { // 判断当前条件是否满足，如果满足则继续下一步
				nc, nis := s.leftQueryNode(nd, ns.nextNode)
				count += nc
				is = append(is, nis...)
			}
		}
	} else {
		if ns == nil {
			return s.leftQueryLeaf(node.(Leaf), nil)
		}
		return s.leftQueryLeaf(node.(Leaf), ns.nextNode)
	}
	return count, is
}

// leftQueryLeaf 叶子节点顺序检索
func (s *Selector) leftQueryLeaf(leaf Leaf, ns *nodeCondition) (int, []interface{}) {
	is := make([]interface{}, 0)
	for _, link := range leaf.getLinks() {
		if inter, err := link.get(); nil == err {
			// todo 条件确定
			for _, cond := range s.Conditions {
				if nil != ns && cond.Param == ns.nss[0].cond.Param {
					continue
				}
			}
			is = append(is, inter)
		}
	}
	return len(leaf.getLinks()), is
}

// rightQueryIndex 索引倒序检索
//
// index 已获取索引对象
func (s *Selector) rightQueryIndex(index Index, ns *nodeCondition) (int, []interface{}) {
	var (
		count, nc int
		nis       []interface{}
		is        = make([]interface{}, 0)
	)
	lenNode := len(index.getNode().getNodes())
	for i := lenNode - 1; i >= 0; i-- {
		if nil == ns {
			nc, nis = s.rightQueryNode(index.getNode().getNodes()[i], nil)
		} else {
			nc, nis = s.rightQueryNode(index.getNode().getNodes()[i], ns.nextNode)
		}
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
	return count, is
}

// rightQueryNode 节点倒序检索
func (s *Selector) rightQueryNode(node Nodal, ns *nodeCondition) (int, []interface{}) {
	count := 0
	is := make([]interface{}, 0)
	if nodes := node.getNodes(); nil != nodes {
		lenNode := len(nodes)
		for i := lenNode - 1; i >= 0; i-- {
			if ns == nil {
				nc, nis := s.rightQueryNode(nodes[i], nil)
				count += nc
				is = append(is, nis...)
			} else if s.conditions(nodes[i], ns.nextNode.nss) { // 判断当前条件是否满足，如果满足则继续下一步
				nc, nis := s.rightQueryNode(nodes[i], ns.nextNode)
				count += nc
				is = append(is, nis...)
			}
		}
	} else {
		if ns == nil {
			return s.rightQueryLeaf(node.(Leaf), nil)
		}
		return s.rightQueryLeaf(node.(Leaf), ns.nextNode)
	}
	return count, is
}

// rightQueryLeaf 叶子节点倒序检索
func (s *Selector) rightQueryLeaf(leaf Leaf, ns *nodeCondition) (int, []interface{}) {
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

// conditions 判断当前条件集合是否满足
func (s *Selector) conditions(node Nodal, nss []*nodeSelector) bool {
	for _, ns := range nss {
		if !s.condition(node, ns) {
			return false
		}
	}
	return true
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
				return s.conditionGT(node, ns)
			case "lt":
				return s.conditionLT(node, ns)
			case "eq":
				return ns.degreeIndex == node.getDegreeIndex()
			case "dif":
				return ns.level == 4 && ns.degreeIndex != node.getDegreeIndex()
			}
		}
	}
	return true
}

// conditionGT 条件大于判断
func (s *Selector) conditionGT(node Nodal, ns *nodeSelector) bool {
	switch ns.level {
	default:
		return ns.degreeIndex < node.getDegreeIndex()
	case 1, 2, 3, 4:
		return ns.degreeIndex <= node.getDegreeIndex()
	}
}

// conditionLT 条件小于判断
func (s *Selector) conditionLT(node Nodal, ns *nodeSelector) bool {
	switch ns.level {
	default:
		return ns.degreeIndex > node.getDegreeIndex()
	case 1, 2, 3, 4:
		return ns.degreeIndex >= node.getDegreeIndex()
	}
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

// hashKeyFromValue 通过Param获取该参数所属hashKey
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

// getConditionNode 根据条件匹配节点单元
//
// 该方法可以用更优雅或正确的方式实现，但烧脑，性能无影响，就这样吧
func (s *Selector) getConditionNode(nc *nodeCondition, cond *condition) {
	var (
		hashKey, flexibleKey, nextFlexibleKey, distance uint64
		nextDegree                                      uint16
		ok                                              bool
	)
	if _, hashKey, ok = type2index(cond.Value); !ok {
		return
	}

	nodeLevel1 := &nodeSelector{level: 1, degreeIndex: 0, cond: cond}
	nc.nss = append(nc.nss, nodeLevel1)
	flexibleKey = hashKey
	distance = levelDistance(nodeLevel1.level)
	nextDegree = uint16(flexibleKey / distance)
	nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance

	nodeLevel2 := &nodeSelector{level: 2, degreeIndex: nextDegree, cond: cond}
	nodeLevel1.nextNode = nodeLevel2
	if nil == nc.nextNode {
		nc.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nss = append(nc.nextNode.nss, nodeLevel2)
	flexibleKey = nextFlexibleKey
	distance = levelDistance(nodeLevel2.level)
	nextDegree = uint16(flexibleKey / distance)
	nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance

	nodeLevel3 := &nodeSelector{level: 3, degreeIndex: nextDegree, cond: cond}
	nodeLevel2.nextNode = nodeLevel3
	if nil == nc.nextNode.nextNode {
		nc.nextNode.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nextNode.nss = append(nc.nextNode.nextNode.nss, nodeLevel3)
	flexibleKey = nextFlexibleKey
	distance = levelDistance(nodeLevel3.level)
	nextDegree = uint16(flexibleKey / distance)
	nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance

	nodeLevel4 := &nodeSelector{level: 4, degreeIndex: nextDegree, cond: cond}
	nodeLevel3.nextNode = nodeLevel4
	if nil == nc.nextNode.nextNode.nextNode {
		nc.nextNode.nextNode.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nextNode.nextNode.nss = append(nc.nextNode.nextNode.nextNode.nss, nodeLevel4)
	flexibleKey = nextFlexibleKey
	distance = levelDistance(nodeLevel4.level)
	nextDegree = uint16(flexibleKey / distance)
	nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance

	nodeLevel5 := &nodeSelector{level: 5, degreeIndex: nextDegree, cond: cond}
	nodeLevel4.nextNode = nodeLevel5
	nodeLevel3.nextNode = nodeLevel4
	if nil == nc.nextNode.nextNode.nextNode.nextNode {
		nc.nextNode.nextNode.nextNode.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nextNode.nextNode.nextNode.nss = append(nc.nextNode.nextNode.nextNode.nextNode.nss, nodeLevel5)
}

// nodeCondition 多个相同Param条件检索预匹配的节点单元
type nodeCondition struct {
	nextNode *nodeCondition
	nss      []*nodeSelector
}

// nodeSelector 条件检索预匹配的节点单元
type nodeSelector struct {
	level       uint8  // 当前节点所在树层级
	degreeIndex uint16 // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	nextNode    *nodeSelector
	cond        *condition
}
