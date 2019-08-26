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

type database interface {
	put(originalKey Key, key uint32, value interface{}) error
	get(originalKey Key, key uint32) (interface{}, error)
	existChild(index uint8) bool
	createChild(index uint8)
}

const (
	//cityCount = 16
	//mallCount    = 128
	//trolleyCount = 128
	//purseCount   = 128
	//boxCount     = 128

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

// String hashes a string to a unique hashcode.
func hash(key Key) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func init() {
	dataDir = env.GetEnv(dataPath)
}

func matchable(matchVal uint8, uintArr []uint8) bool {
	_, err := binaryMatch(matchVal, uintArr)
	return nil == err
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
