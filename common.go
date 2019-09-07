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
	"hash/crc32"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// distance 根据节点所在层级获取当前节点内部子节点之间的差
func distance(level uint8) uint32 {
	switch level {
	case 0:
		return mallDistance
	case 1:
		return trolleyDistance
	case 2:
		return purseDistance
	case 3:
		return boxDistance
	}
	return 0
}

// String hashes a string to a unique hashcode.
func hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// matchableData 'Nodal'内子节点数组二分查找是否存在指定值
//
// matchIndex 要查找的值
//
// matcher 'binaryMatcher'接口支持的获取‘Nodal’接口的内置方法对象
//
// bool 返回存在与否
func matchableData(matchIndex uint8, matcher Nodal) bool {
	_, err := binaryMatchData(matchIndex, matcher)
	return nil == err
}

// binaryMatchData 'Nodal'内子节点数组二分查找基本方法
//
// matchIndex 要查找的值
//
// matcher 'binaryMatcher'接口支持的获取‘Nodal’接口的内置方法对象
//
// realIndex 返回查找到的真实的元素下标，该下标是对应数组内的下标，并非树中节点数组原型的下标
//
// 如果没找到，则返回err
func binaryMatchData(matchIndex uint8, matcher Nodal) (realIndex int, err error) {
	var (
		leftIndex   int
		middleIndex int
		rightIndex  int
	)
	leftIndex = 0
	rightIndex = matcher.childCount() - 1
	for leftIndex <= rightIndex {
		middleIndex = (leftIndex + rightIndex) / 2
		// 如果要找的数比midVal大
		if matcher.child(middleIndex).getDegreeIndex() > matchIndex {
			// 在arr数组的左边找
			rightIndex = middleIndex - 1
		} else if matcher.child(middleIndex).getDegreeIndex() < matchIndex {
			// 在arr数组的右边找
			leftIndex = middleIndex + 1
		} else if matcher.child(middleIndex).getDegreeIndex() == matchIndex {
			return middleIndex, nil
		}
	}
	return 0, errors.New("catalog is nil")
}

// binaryMatch 数组内二分查找基本方法
//
// matchVal 要查找的值
//
// uintArr 在‘uintArr’数组中检索
//
// catalog 返回查找到的在数组‘uintArr’中的元素下标
//
// 如果没找到，则返回err
func binaryMatch(matchVal uint8, uintArr []uint8) (index int, err error) {
	var (
		leftIndex   int
		middleIndex int
		rightIndex  int
	)
	leftIndex = 0
	rightIndex = len(uintArr) - 1
	for leftIndex <= rightIndex {
		middleIndex = (leftIndex + rightIndex) / 2
		// 如果要找的数比midVal大
		if uintArr[middleIndex] > matchVal {
			// 在arr数组的左边找
			rightIndex = middleIndex - 1
		} else if uintArr[middleIndex] < matchVal {
			// 在arr数组的右边找
			leftIndex = middleIndex + 1
		} else if uintArr[middleIndex] == matchVal {
			return middleIndex, nil
		}
	}
	return 0, errors.New("catalog is nil")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// mkDataDir 创建库存储目录
func mkDataDir(dataName string) (err error) {
	dataPath := filepath.Join(dataDir, dataName)
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
		return ErrDatabaseExist
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}

// rmDataDir 删除库存储目录
func rmDataDir(dataName string) (err error) {
	dataPath := filepath.Join(dataDir, dataName)
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
		return os.Remove(dataPath)
	}
	return nil
}

// mkFormResourceDoc 创建表资源
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
func mkFormResourceSQL(dataID, formID, indexID string, fileIndex int) (err error) {
	if err = mkFormDir(dataID, formID); nil != err {
		return
	}
	if err = mkFormIndexDir(dataID, formID, indexID); nil != err {
		_ = rmFormDir(dataID, formID)
		return
	}
	if err = mkFormDataFile(dataID, formID, fileIndex); nil != err {
		_ = rmFormDir(dataID, formID)
		return
	}
	return
}

// mkFormResourceDoc 创建表资源
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
//
// customID put keyStructure catalog id
func mkFormResourceDoc(dataID, formID, customID string, fileIndex int) (err error) {
	if err = mkFormDir(dataID, formID); nil != err {
		return
	}
	if err = mkFormIndexDir(dataID, formID, customID); nil != err {
		_ = rmFormDir(dataID, formID)
		return
	}
	if err = mkFormDataFile(dataID, formID, fileIndex); nil != err {
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
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
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
	if exist, err := pathExist(formPath); nil != err {
		return err
	} else if exist {
		return os.Remove(formPath)
	}
	return nil
}

// mkFormIndexDir 创建表索引目录
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
func mkFormIndexDir(dataID, formID, indexID string) (err error) {
	return os.Mkdir(pathFormIndexDir(dataID, formID, indexID), os.ModePerm)
}

// rmFormIndexDir 删除表索引目录
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
func rmFormIndexDir(dataID, formID, indexID string) (err error) {
	indexPath := pathFormIndexDir(dataID, formID, indexID)
	if exist, err := pathExist(indexPath); nil != err {
		return err
	} else if exist {
		return os.Remove(indexPath)
	}
	return nil
}

// mkFormIndexDir 创建表索引文件
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
//
// catalog 所在表顶层数组中下标
func mkFormDataFile(dataID, formID string, fileIndex int) (err error) {
	_, err = os.Create(pathFormDataFile(dataID, formID, fileIndex))
	return
}

// pathFormDir 表目录
//
// dataID 数据库唯一id
//
// formID 表唯一id
func pathFormDir(dataID, formID string) string {
	return filepath.Join(dataDir, dataID, formID)
}

// pathFormIndexDir 表索引目录
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
func pathFormIndexDir(dataID, formID, indexID string) string {
	return filepath.Join(dataDir, dataID, formID, indexID)
}

// pathFormIndexFile 表索引文件路径
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
//
// catalog 所在表顶层数组中下标
func pathFormIndexFile(dataID, formID, indexID string, index uint8) string {
	return strings.Join([]string{pathFormIndexDir(dataID, formID, indexID), string(filepath.Separator), strconv.Itoa(int(index)), ".idx"}, "")
}

func pathFormDataFile(dataID, formID string, fileIndex int) string {
	return strings.Join([]string{dataDir, string(filepath.Separator), dataID, string(filepath.Separator), formID, string(filepath.Separator), strconv.Itoa(fileIndex), ".dat"}, "")
}

// pathExist 检查路径是否存在
func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// int16ToUint32Index int16转uint32
//
// 该方法仅对索引有效，在排序或范围查询时，如果有使用int16做索引的对象，都必须按照该方法转换并做内部检索和排序，如thing
func int16ToUint32Index(i16 int16) uint32 {
	ui16 := i16 + 32767 + 1 // 32768 = 1 << 15
	return uint32(ui16)
}

// int32ToUint32Index int32转uint32
//
// 该方法仅对索引有效，在排序或范围查询时，如果有使用int32做索引的对象，都必须按照该方法转换并做内部检索和排序，如thing
func int32ToUint32Index(i32 int32) uint32 {
	ui32 := i32 + 2147483647 + 1 // 2147483648 = 1 << 31
	return uint32(ui32)
}

// int64ToUint32Index int64转uint32
//
// 该方法仅对索引有效，在排序或范围查询时，如果有使用int64做索引的对象，都必须按照该方法转换并做内部检索和排序，如thing
func int64ToUint32Index(i64 int64) uint32 {
	ui64 := i64 + 9223372036854775807 + 1 // 9.223372036854776e18 || 9223372036854775808 = 1 << 63
	return uint64ToUint32Index(uint64(ui64))
}

// uint64ToUint32Index uint64转uint32
//
// 该方法仅对索引有效，在排序或范围查询时，如果有使用uint64做索引的对象，都必须按照该方法转换并做内部检索和排序，如thing
func uint64ToUint32Index(ui64 uint64) uint32 {
	return uint64ToUint32IndexDivide(ui64, 0)
}

// uint64ToUint32Index 毫秒时间戳转uint32
//
// divide 除算次数，首次为0
func uint64ToUint32IndexDivide(ui64 uint64, divide uint64) uint32 {
	if ui64 >= 4294967296 { // 4294967296 = 2 << 31
		if divide == 0 {
			divide++
			return uint64ToUint32IndexDivide(ui64, divide)
		}
		i64New := ui64 / (10 * divide)
		if i64New > 4294967296 {
			divide++
			return uint64ToUint32IndexDivide(ui64, divide)
		}
		return uint32(i64New)
	}
	return uint32(ui64)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
