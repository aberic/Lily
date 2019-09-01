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

package Lily

import (
	"errors"
	"hash/crc32"
	"math"
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
func hash(key Key) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// matchableData 'nodal'内子节点数组二分查找是否存在指定值
//
// matchIndex 要查找的值
//
// matcher 'binaryMatcher'接口支持的获取‘nodal’接口的内置方法对象
//
// bool 返回存在与否
func matchableData(matchIndex uint8, matcher nodal) bool {
	_, err := binaryMatchData(matchIndex, matcher)
	return nil == err
}

// binaryMatchData 'nodal'内子节点数组二分查找基本方法
//
// matchIndex 要查找的值
//
// matcher 'binaryMatcher'接口支持的获取‘nodal’接口的内置方法对象
//
// realIndex 返回查找到的真实的元素下标，该下标是对应数组内的下标，并非树中节点数组原型的下标
//
// 如果没找到，则返回err
func binaryMatchData(matchIndex uint8, matcher nodal) (realIndex int, err error) {
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
	return 0, errors.New("index is nil")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// mkDataDir 创建库存储目录
func mkDataDir(dataName string) (err error) {
	dataPath := filepath.Join(dataDir, dataName)
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
		return databaseExistErr
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

// mkFormResource 创建表资源
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
//
// customID put key index id
func mkFormResource(dataID, formID, indexID, customID string, fileIndex int) (err error) {
	if err = mkFormDir(dataID, formID); nil != err {
		return
	}
	if err = mkFormIndexDir(dataID, formID, indexID); nil != err {
		_ = rmFormDir(dataID, formID)
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
		return formExistErr
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
// index 所在表顶层数组中下标
func mkFormIndexFile(dataID, formID, indexID string, index uint8) (err error) {
	_, err = os.Create(pathFormIndexFile(dataID, formID, indexID, index))
	return
}

// mkFormIndexDir 创建表索引文件
//
// dataID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
//
// index 所在表顶层数组中下标
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
// index 所在表顶层数组中下标
func pathFormIndexFile(dataID, formID, indexID string, index uint8) string {
	return strings.Join([]string{pathFormIndexDir(dataID, formID, indexID), string(filepath.Separator), strconv.Itoa(int(index)), ".idx"}, "")
}

func pathFormDataFile(dataID, formID string, fileIndex int) string {
	return strings.Join([]string{dataDir, string(filepath.Separator), dataID, string(filepath.Separator), formID, string(filepath.Separator), strconv.Itoa(fileIndex), ".dat"}, "")
}

func pathFormNodeDir(dataID, formID string, index uint8) string {
	return filepath.Join(pathFormDir(dataID, formID), strconv.Itoa(int(index)))
}

func pathFormNodeFile(dataID, formID, fileName string, index uint8) string {
	return filepath.Join(dataDir, dataID, formID, strconv.Itoa(int(index)), fileName)
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

// hexIntMap 十六进制对应十进制映射
var hexInt64Map = map[int64]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
}

// hexIntMap 十六进制对应十进制映射
var hexIntMap = map[int]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
}

// intHexMap 十进制对应十六进制映射
var intHexMap = map[string]int{
	"0": 0, "1": 1, "2": 2, "3": 3,
	"4": 4, "5": 5, "6": 6, "7": 7,
	"8": 8, "9": 9, "a": 10, "b": 11,
	"c": 12, "d": 13, "e": 14, "f": 15,
}

// int32ToHexString int32转十六进制字符串
func int32ToHexString(i int) string {
	hexStrArr := make([]string, 8)
	for index := 7; index >= 0; index-- {
		if i >= 16 {
			hexStrArr[index] = hexIntMap[i%16]
			i /= 16
		} else if i >= 0 && i < 16 {
			hexStrArr[index] = hexIntMap[i]
			i = 0
		} else {
			hexStrArr[index] = hexIntMap[0]
		}
	}
	return strings.Join(hexStrArr, "")
}

// int64ToHexString int64转十六进制字符串
func int64ToHexString(i int64) string {
	hexStrArr := make([]string, 16)
	for index := 15; index >= 0; index-- {
		if i >= 16 {
			hexStrArr[index] = hexInt64Map[i%16]
			i /= 16
		} else if i >= 0 && i < 16 {
			hexStrArr[index] = hexInt64Map[i]
			i = 0
		} else {
			hexStrArr[index] = hexIntMap[0]
		}
	}
	return strings.Join(hexStrArr, "")
}

// int64ToHexString int字符串转int
func hexStringToInt64(hex string) (int64, error) {
	math.Pow(1, 2)
	hexLen := len(hex)
	intHex := 0
	for i := 0; i < hexLen; i++ {
		intHex += intHexMap[hex[i:i+1]] * int(math.Pow(16, float64(hexLen-i-1)))
	}
	return int64(intHex), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
