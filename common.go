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

// mkFormResource 创建表资源
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

// duoInt64Map 三十二进制对应十进制映射
var duoInt64Map = map[int64]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
	16: "g", 17: "h", 18: "i", 19: "j",
	20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r",
	28: "s", 29: "t", 30: "u", 31: "v",
}

// dDuoInt64Map 64进制对应十进制映射
var dDuoInt64Map = map[int64]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
	16: "g", 17: "h", 18: "i", 19: "j",
	20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r",
	28: "s", 29: "t", 30: "u", 31: "v",
	33: "w", 34: "x", 35: "y", 36: "z",
	37: "A", 38: "B", 39: "C", 40: "D",
	41: "E", 42: "F", 43: "G", 44: "H",
	45: "I", 46: "J", 47: "K", 48: "L",
	49: "M", 50: "N", 51: "O", 52: "P",
	53: "Q", 54: "R", 55: "S", 56: "T",
	57: "U", 58: "V", 59: "W", 60: "X",
	61: "Y", 62: "Z", 63: "+", 64: "-",
}

// dDuoInt64Map 64进制对应十进制映射
var dDuoUint32Map = map[uint32]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
	16: "g", 17: "h", 18: "i", 19: "j",
	20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r",
	28: "s", 29: "t", 30: "u", 31: "v",
	32: "w", 33: "x", 34: "y", 35: "z",
	36: "A", 37: "B", 38: "C", 39: "D",
	40: "E", 41: "F", 42: "G", 43: "H",
	44: "I", 45: "J", 46: "K", 47: "L",
	48: "M", 49: "N", 50: "O", 51: "P",
	52: "Q", 53: "R", 54: "S", 55: "T",
	56: "U", 57: "V", 58: "W", 59: "X",
	60: "Y", 61: "Z", 62: "+", 63: "-",
}

// dDuoIntMap 64进制对应十进制映射
var dDuoIntMap = map[int]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
	16: "g", 17: "h", 18: "i", 19: "j",
	20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r",
	28: "s", 29: "t", 30: "u", 31: "v",
	32: "w", 33: "x", 34: "y", 35: "z",
	36: "A", 37: "B", 38: "C", 39: "D",
	40: "E", 41: "F", 42: "G", 43: "H",
	44: "I", 45: "J", 46: "K", 47: "L",
	48: "M", 49: "N", 50: "O", 51: "P",
	52: "Q", 53: "R", 54: "S", 55: "T",
	56: "U", 57: "V", 58: "W", 59: "X",
	60: "Y", 61: "Z", 62: "+", 63: "-",
}

// duoIntMap 三十二进制对应十进制映射
var duoIntMap = map[int]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
	16: "g", 17: "h", 18: "i", 19: "j",
	20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r",
	28: "s", 29: "t", 30: "u", 31: "v",
}

// duoIntMap 三十二进制对应十进制映射
var duoUint32Map = map[uint32]string{
	0: "0", 1: "1", 2: "2", 3: "3",
	4: "4", 5: "5", 6: "6", 7: "7",
	8: "8", 9: "9", 10: "a", 11: "b",
	12: "c", 13: "d", 14: "e", 15: "f",
	16: "g", 17: "h", 18: "i", 19: "j",
	20: "k", 21: "l", 22: "m", 23: "n",
	24: "o", 25: "p", 26: "q", 27: "r",
	28: "s", 29: "t", 30: "u", 31: "v",
}

// intDuoMap 十进制对应三十二进制映射
var intDuoMap = map[string]int{
	"0": 0, "1": 1, "2": 2, "3": 3,
	"4": 4, "5": 5, "6": 6, "7": 7,
	"8": 8, "9": 9, "a": 10, "b": 11,
	"c": 12, "d": 13, "e": 14, "f": 15,
	"g": 16, "h": 17, "i": 18, "j": 19,
	"k": 20, "l": 21, "m": 22, "n": 23,
	"o": 24, "p": 25, "q": 26, "r": 27,
	"s": 28, "t": 29, "u": 30, "v": 31,
}

// uint32DuoMap 十进制对应三十二进制映射
var uint32DuoMap = map[string]uint32{
	"0": 0, "1": 1, "2": 2, "3": 3,
	"4": 4, "5": 5, "6": 6, "7": 7,
	"8": 8, "9": 9, "a": 10, "b": 11,
	"c": 12, "d": 13, "e": 14, "f": 15,
	"g": 16, "h": 17, "i": 18, "j": 19,
	"k": 20, "l": 21, "m": 22, "n": 23,
	"o": 24, "p": 25, "q": 26, "r": 27,
	"s": 28, "t": 29, "u": 30, "v": 31,
}

// int64DDuoMap 十进制对应64进制映射
var int64DDuoMap = map[string]int64{
	"0": 0, "1": 1, "2": 2, "3": 3,
	"4": 4, "5": 5, "6": 6, "7": 7,
	"8": 8, "9": 9, "a": 10, "b": 11,
	"c": 12, "d": 13, "e": 14, "f": 15,
	"g": 16, "h": 17, "i": 18, "j": 19,
	"k": 20, "l": 21, "m": 22, "n": 23,
	"o": 24, "p": 25, "q": 26, "r": 27,
	"s": 28, "t": 29, "u": 30, "v": 31,
	"w": 32, "x": 33, "y": 34, "z": 35,
	"A": 36, "B": 37, "C": 38, "D": 39,
	"E": 40, "F": 41, "G": 42, "H": 43,
	"I": 44, "J": 45, "K": 46, "L": 47,
	"M": 48, "N": 49, "O": 50, "P": 51,
	"Q": 52, "R": 53, "S": 54, "T": 55,
	"U": 56, "V": 57, "W": 58, "X": 59,
	"Y": 60, "Z": 61, "+": 62, "-": 63,
}

// uint32DDuoMap 十进制对应64进制映射
var uint32DDuoMap = map[string]uint32{
	"0": 0, "1": 1, "2": 2, "3": 3,
	"4": 4, "5": 5, "6": 6, "7": 7,
	"8": 8, "9": 9, "a": 10, "b": 11,
	"c": 12, "d": 13, "e": 14, "f": 15,
	"g": 16, "h": 17, "i": 18, "j": 19,
	"k": 20, "l": 21, "m": 22, "n": 23,
	"o": 24, "p": 25, "q": 26, "r": 27,
	"s": 28, "t": 29, "u": 30, "v": 31,
	"w": 32, "x": 33, "y": 34, "z": 35,
	"A": 36, "B": 37, "C": 38, "D": 39,
	"E": 40, "F": 41, "G": 42, "H": 43,
	"I": 44, "J": 45, "K": 46, "L": 47,
	"M": 48, "N": 49, "O": 50, "P": 51,
	"Q": 52, "R": 53, "S": 54, "T": 55,
	"U": 56, "V": 57, "W": 58, "X": 59,
	"Y": 60, "Z": 61, "+": 62, "-": 63,
}

// int2DDuoMap 十进制对应64进制映射
var intDDuoMap = map[string]int{
	"0": 0, "1": 1, "2": 2, "3": 3,
	"4": 4, "5": 5, "6": 6, "7": 7,
	"8": 8, "9": 9, "a": 10, "b": 11,
	"c": 12, "d": 13, "e": 14, "f": 15,
	"g": 16, "h": 17, "i": 18, "j": 19,
	"k": 20, "l": 21, "m": 22, "n": 23,
	"o": 24, "p": 25, "q": 26, "r": 27,
	"s": 28, "t": 29, "u": 30, "v": 31,
	"w": 32, "x": 33, "y": 34, "z": 35,
	"A": 36, "B": 37, "C": 38, "D": 39,
	"E": 40, "F": 41, "G": 42, "H": 43,
	"I": 44, "J": 45, "K": 46, "L": 47,
	"M": 48, "N": 49, "O": 50, "P": 51,
	"Q": 52, "R": 53, "S": 54, "T": 55,
	"U": 56, "V": 57, "W": 58, "X": 59,
	"Y": 60, "Z": 61, "+": 62, "-": 63,
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
//
// 1073741824
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
func hexStringToInt32(hex string) (int32, error) {
	hexLen := len(hex)
	intHex := 0
	for i := 0; i < hexLen; i++ {
		intHex += intHexMap[hex[i:i+1]] * int(math.Pow(16, float64(hexLen-i-1)))
	}
	return int32(intHex), nil
}

// int64ToHexString int字符串转int
func hexStringToInt64(hex string) (int64, error) {
	hexLen := len(hex)
	intHex := 0
	for i := 0; i < hexLen; i++ {
		intHex += intHexMap[hex[i:i+1]] * int(math.Pow(16, float64(hexLen-i-1)))
	}
	return int64(intHex), nil
}

// intToDuoString int32转32进制字符串
func intToDuoString(i int) string {
	dupStrArr := make([]string, 6)
	for index := 5; index >= 0; index-- {
		if i >= 32 {
			dupStrArr[index] = duoIntMap[i%32]
			i /= 32
		} else if i >= 0 && i < 32 {
			dupStrArr[index] = duoIntMap[i]
			i = 0
		} else {
			dupStrArr[index] = duoIntMap[0]
		}
	}
	return strings.Join(dupStrArr, "")
}

// uint32ToDuoString uint32转32进制字符串
func uint32ToDuoString(i uint32) string {
	dupStrArr := make([]string, 6)
	for index := 5; index >= 0; index-- {
		if i >= 32 {
			dupStrArr[index] = duoUint32Map[i%32]
			i /= 32
		} else if i < 32 {
			dupStrArr[index] = duoUint32Map[i]
			i = 0
		} else {
			dupStrArr[index] = duoUint32Map[0]
		}
	}
	return strings.Join(dupStrArr, "")
}

// intToDDuoString uint32转32进制字符串
func intToDDuoString(i int) string {
	dupStrArr := make([]string, 4)
	for index := 3; index >= 0; index-- {
		if i >= 64 {
			dupStrArr[index] = dDuoIntMap[i%64]
			i /= 64
		} else if i < 64 {
			dupStrArr[index] = dDuoIntMap[i]
			i = 0
		} else {
			dupStrArr[index] = dDuoIntMap[0]
		}
	}
	return strings.Join(dupStrArr, "")
}

// uint32ToDDuoString uint32转32进制字符串
func uint32ToDDuoString(i uint32) string {
	dupStrArr := make([]string, 5)
	for index := 4; index >= 0; index-- {
		if i >= 64 {
			dupStrArr[index] = dDuoUint32Map[i%64]
			i /= 64
		} else if i < 64 {
			dupStrArr[index] = dDuoUint32Map[i]
			i = 0
		} else {
			dupStrArr[index] = dDuoUint32Map[0]
		}
	}
	return strings.Join(dupStrArr, "")
}

// int64ToDDuoString uint32转32进制字符串
func int64ToDDuoString(i int64) string {
	dupStrArr := make([]string, 5)
	for index := 4; index >= 0; index-- {
		if i >= 64 {
			dupStrArr[index] = dDuoInt64Map[i%64]
			i /= 64
		} else if i < 64 {
			dupStrArr[index] = dDuoInt64Map[i]
			i = 0
		} else {
			dupStrArr[index] = dDuoInt64Map[0]
		}
	}
	return strings.Join(dupStrArr, "")
}

// int64ToDuoString int64转32进制字符串
//
// 1073741824
func int64ToDuoString(i int64) string {
	duoStrArr := make([]string, 5)
	for index := 5; index >= 0; index-- {
		if i >= 32 {
			duoStrArr[index] = duoInt64Map[i%32]
			i /= 32
		} else if i >= 0 && i < 32 {
			duoStrArr[index] = duoInt64Map[i]
			i = 0
		} else {
			duoStrArr[index] = duoInt64Map[0]
		}
	}
	return strings.Join(duoStrArr, "")
}

// duoStringToInt32 int字符串转int
func duoStringToInt32(duo string) (int32, error) {
	duoLen := len(duo)
	intDuo := 0
	for i := 0; i < duoLen; i++ {
		intDuo += intDuoMap[duo[i:i+1]] * int(math.Pow(32, float64(duoLen-i-1)))
	}
	return int32(intDuo), nil
}

// duoStringToUint32 int字符串转int
func duoStringToUint32(duo string) (uint32, error) {
	duoLen := len(duo)
	var intDuo uint32
	intDuo = 0
	for i := 0; i < duoLen; i++ {
		intDuo += uint32DuoMap[duo[i:i+1]] * uint32(math.Pow(32, float64(duoLen-i-1)))
	}
	return intDuo, nil
}

// dDuoStringToInt int字符串转int
func dDuoStringToInt(duo string) (int, error) {
	dDuoLen := len(duo)
	intDuo := 0
	for i := 0; i < dDuoLen; i++ {
		intDuo += intDDuoMap[duo[i:i+1]] * int(math.Pow(64, float64(dDuoLen-i-1)))
	}
	return intDuo, nil
}

// dDuoStringToUint32 int字符串转int
func dDuoStringToUint32(duo string) (uint32, error) {
	dDuoLen := len(duo)
	var intDuo uint32
	intDuo = 0
	for i := 0; i < dDuoLen; i++ {
		intDuo += uint32DDuoMap[duo[i:i+1]] * uint32(math.Pow(64, float64(dDuoLen-i-1)))
	}
	return intDuo, nil
}

// dDuoStringToUint32 int字符串转int
func dDuoStringToint64(duo string) (int64, error) {
	dDuoLen := len(duo)
	var intDuo int64
	intDuo = 0
	for i := 0; i < dDuoLen; i++ {
		intDuo += int64DDuoMap[duo[i:i+1]] * int64(math.Pow(64, float64(dDuoLen-i-1)))
	}
	return intDuo, nil
}

// duoStringToInt64 int字符串转int
func duoStringToInt64(duo string) (int64, error) {
	duoLen := len(duo)
	intDuo := 0
	for i := 0; i < duoLen; i++ {
		intDuo += intDuoMap[duo[i:i+1]] * int(math.Pow(32, float64(duoLen-i-1)))
	}
	return int64(intDuo), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
