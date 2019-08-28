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
	"errors"
	"github.com/ennoo/rivet/utils/env"
	"hash/crc32"
)

// binaryMatcher 二分查询辅助接口
type binaryMatcher interface {
	// childCount 获取子节点集合数量
	childCount() int
	// child 根据子节点集合下标获取树-度对象
	child(index int) nodal
}

// position 父子节点游标辅助接口
type position interface {
	// getPreNodal 获取父节点对象
	getPreNodal() nodal
}

// nodal 节点对象接口
type nodal interface {
	binaryMatcher
	position
	// put 插入数据
	//
	// originalKey 真实key，必须string类型
	//
	// key 索引key，可通过hash转换string生成
	//
	// value 存储对象
	put(originalKey Key, key uint32, value interface{}) error
	// get 获取数据，返回存储对象
	//
	// originalKey 真实key，必须string类型
	//
	// key 索引key，可通过hash转换string生成
	get(originalKey Key, key uint32) (interface{}, error)
	// existChild 根据下标判定是否存在子节点
	existChild(index uint8) bool
	// createChild 根据下标创建新的子节点
	createChild(index uint8) nodal
	// getFlexibleKey 下一级最左最小树所对应真实key
	getFlexibleKey() uint32
	// getDegreeIndex 获取节点所在树中度集合中的数组下标
	getDegreeIndex() uint8
	lock()
	unLock()
	rLock()
	rUnLock()
}

const (
	//cityCount = 16
	//mallCount    = 128
	//trolleyCount = 128
	//purseCount   = 128
	//boxCount     = 128

	levelMax uint8 = 3
	//degreeMax uint8 = 128
	// 最大存储数，超过次数一律做新值换算
	//lilyMax      uint32 = 4294967280
	cityDistance uint32 = 268435455
	// mallDistance level1间隔 ld1=(treeCount+1)/128=2097152 128^3
	mallDistance uint32 = 2097152
	// trolleyDistance level2间隔 ld2=(16513*127+1)/128=16384 128^2
	trolleyDistance uint32 = 16384
	// purseDistance level3间隔 ld3=(129*127+1)/128=128 128^1
	purseDistance uint32 = 128
	// boxDistance level4间隔 ld3=(1*127+1)/128=1 128^0
	boxDistance uint32 = 1

	dataPath = "DATA_PATH"
)

//const (
//	//cityCount    = 1
//	//mallCount    = 4
//	//trolleyCount = 4
//	//purseCount   = 4
//	//boxCount     = 4
//
//	levelMax uint8 = 3
//	cityDistance uint32 = 0
//	// mallDistance level1间隔 ld1=(treeCount+1)/128=2097152 128^3
//	mallDistance uint32 = 64
//	// trolleyDistance level2间隔 ld2=(16513*127+1)/128=16384 128^2
//	trolleyDistance uint32 = 16
//	// purseDistance level3间隔 ld3=(129*127+1)/128=128 128^1
//	purseDistance uint32 = 4
//	// boxDistance level4间隔 ld3=(1*127+1)/128=1 128^0
//	boxDistance uint32 = 1
//
//	dataPath = "DATA_PATH"
//)

var (
	dataDir string
)

type Key string

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

func init() {
	dataDir = env.GetEnvDefault(dataPath, "/Users/aberic/Documents/tmp/lily/t1")
}

func matchableData(matchVal uint8, matcher binaryMatcher) bool {
	_, err := binaryMatchData(matchVal, matcher)
	return nil == err
}

func binaryMatchData(matchIndex uint8, matcher binaryMatcher) (realIndex int, err error) {
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
