/*
 * Copyright (c) 2019.. Aberic - All Rights Reserved.
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
	"hash/crc32"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// levelDistance 根据节点所在层级获取当前节点内部子节点之间的差
func levelDistance(level uint8) uint64 {
	switch level {
	case 1:
		return level1Distance
	case 2:
		return level2Distance
	case 3:
		return level3Distance
	case 4:
		return level4Distance
	}
	return 0
}

// String hashes a string to a unique hashcode.
func hash(key string) uint64 {
	return uint64(crc32.ChecksumIEEE([]byte(key)))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// matchableData 'Nodal'内子节点数组二分查找是否存在指定值
//
// matchIndex 要查找的值
//
// data 'binaryMatcher'接口支持的获取‘Nodal’接口的内置方法对象
//
// bool 返回存在与否
//func matchableData(matchIndex uint16, data Data) bool {
//	_, err := binaryMatchData(matchIndex, data)
//	return nil == err
//}

// binaryMatchData 'Nodal'内子节点数组二分查找基本方法
//
// matchIndex 要查找的值
//
// data 'binaryMatcher'接口支持的获取‘Nodal’接口的内置方法对象
//
// realIndex 返回查找到的真实的元素下标，该下标是对应数组内的下标，并非树中节点数组原型的下标
//
// 如果没找到，则返回err
func binaryMatchData(matchIndex uint16, node Nodal) (realIndex int, err error) {
	var (
		leftIndex   int
		middleIndex int
		rightIndex  int
	)
	leftIndex = 0
	nodes := node.getNodes()
	rightIndex = len(nodes) - 1
	for leftIndex <= rightIndex {
		middleIndex = (leftIndex + rightIndex) / 2
		// 如果要找的数比midVal大
		if nodes[middleIndex].getDegreeIndex() > matchIndex {
			// 在arr数组的左边找
			rightIndex = middleIndex - 1
		} else if nodes[middleIndex].getDegreeIndex() < matchIndex {
			// 在arr数组的右边找
			leftIndex = middleIndex + 1
		} else if nodes[middleIndex].getDegreeIndex() == matchIndex {
			return middleIndex, nil
		}
	}
	return 0, errors.New("index is nil")
}

// binaryMatch 数组内二分查找基本方法
//
// matchVal 要查找的值
//
// uintArr 在‘uintArr’数组中检索
//
// index 返回查找到的在数组‘uintArr’中的元素下标
//
// 如果没找到，则返回err
//func binaryMatch(matchVal uint8, uintArr []uint8) (index int, err error) {
//	var (
//		leftIndex   int
//		middleIndex int
//		rightIndex  int
//	)
//	leftIndex = 0
//	rightIndex = len(uintArr) - 1
//	for leftIndex <= rightIndex {
//		middleIndex = (leftIndex + rightIndex) / 2
//		// 如果要找的数比midVal大
//		if uintArr[middleIndex] > matchVal {
//			// 在arr数组的左边找
//			rightIndex = middleIndex - 1
//		} else if uintArr[middleIndex] < matchVal {
//			// 在arr数组的右边找
//			leftIndex = middleIndex + 1
//		} else if uintArr[middleIndex] == matchVal {
//			return middleIndex, nil
//		}
//	}
//	return 0, errors.New("index is nil")
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// mkDataDir 创建库存储目录
func mkDataDir(dataName string) (err error) {
	dataPath := filepath.Join(obtainConf().DataDir, dataName)
	if gnomon.File().PathExists(dataPath) {
		return ErrDatabaseExist
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}

// rmDataDir 删除库存储目录
func rmDataDir(dataName string) (err error) {
	dataPath := filepath.Join(obtainConf().DataDir, dataName)
	if gnomon.File().PathExists(dataPath) {
		return os.Remove(dataPath)
	}
	return nil
}

// mkFormResource 创建表资源
//
// dataID 数据库唯一id
//
// formID 表唯一id
func mkFormResource(dataID, formID string) (err error) {
	if err = mkFormDir(dataID, formID); nil != err {
		return
	}
	if err = mkFormDataFile(dataID, formID); nil != err {
		_ = rmFormDir(dataID, formID)
		return
	}
	return
}

// mkFormDir 创建表存储目录
//
// dataID 数据库唯一id
//
// formID 表唯一id
func mkFormDir(dataID, formID string) (err error) {
	dataPath := pathFormDir(dataID, formID)
	if gnomon.File().PathExists(dataPath) {
		return ErrFormExist
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}

// rmFormDir 删除表存储目录
//
// dataID 数据库唯一id
//
// formID 表唯一id
func rmFormDir(dataID, formID string) (err error) {
	formPath := pathFormDir(dataID, formID)
	if gnomon.File().PathExists(formPath) {
		return os.Remove(formPath)
	}
	return nil
}

// mkFormDataFile 创建表文件
//
// dataID 数据库唯一id
//
// formID 表唯一id
func mkFormDataFile(dataID, formID string) (err error) {
	_, err = os.Create(pathFormDataFile(dataID, formID))
	return
}

// pathFormDir 表目录
//
// dataID 数据库唯一id
//
// formID 表唯一id
func pathFormDir(dataID, formID string) string {
	return filepath.Join(obtainConf().DataDir, dataID, formID)
}

// pathFormIndexFile 表索引文件路径
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
func pathFormIndexFile(dataID, formID, indexID string) string {
	return strings.Join([]string{obtainConf().DataDir, string(filepath.Separator), dataID, string(filepath.Separator), formID, string(filepath.Separator), indexID, ".idx"}, "")
}

func pathFormDataFile(dataID, formID string) string {
	return filepath.Join(obtainConf().DataDir, dataID, formID, "form.dat")
	//return strings.Join([]string{dataDir, string(filepath.Separator), dataID, string(filepath.Separator), formID, string(filepath.Separator), strconv.Itoa(fileIndex), ".dat"}, "")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func valueType2index(value *reflect.Value) (key string, hashKey uint64, support bool) {
	support = true
	switch value.Kind() {
	default:
		return "", 0, false
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		i64 := value.Int()
		key = strconv.FormatInt(i64, 10)
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uintptr:
		ui64 := value.Uint()
		key = strconv.FormatUint(ui64, 10)
		if ui64 > 9223372036854775807 { // 9223372036854775808 = 1 << 63
			return "", 0, false
		}
		hashKey = ui64 + 9223372036854775807 + 1
	case reflect.Float32, reflect.Float64:
		i64 := gnomon.Scale().Float64toInt64(value.Float(), 4)
		key = strconv.FormatInt(i64, 10)
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case reflect.String:
		key = value.String()
		hashKey = hash(key)
	case reflect.Bool:
		if value.Bool() {
			key = value.String()
			hashKey = 1
		} else {
			key = value.String()
			hashKey = 2
		}
	}
	return
}

func value2hashKey(value *reflect.Value) (hashKey uint64, support bool) {
	support = true
	switch value.Kind() {
	default:
		return 0, false
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		i64 := value.Int()
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uintptr:
		ui64 := value.Uint()
		if ui64 > 9223372036854775808 { // 9223372036854775808 = 1 << 63
			return 0, false
		}
		hashKey = ui64 + 9223372036854775807 + 1
	case reflect.Float32, reflect.Float64:
		hashKey = uint64(gnomon.Scale().Float64toInt64(value.Float(), 4) + 9223372036854775807 + 1)
	case reflect.String:
		hashKey = hash(value.String())
	case reflect.Bool:
		if value.Bool() {
			hashKey = 1
		} else {
			hashKey = 2
		}
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
