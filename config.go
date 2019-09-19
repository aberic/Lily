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
	"github.com/aberic/gnomon"
	"path/filepath"
	"strconv"
)

const (
	hashCount = 16
	//nodalCount = 128
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
)

//const (
//	//hashCount    = 1
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
	rootDir       string // Lily服务默认存储路径
	dataDir       string // Lily服务默认存储路径
	limitOpenFile int    // 限制打开文件描述符次数
)

func init() {
	var err error
	rootDir = gnomon.Env().GetD("DATA_PATH", "test/t1")
	dataDir = filepath.Join(rootDir, "Data")
	if limitOpenFile, err = strconv.Atoi(gnomon.Env().GetD("LIMIT_COUNT", "10000")); nil == err {
		limitOpenFile = 10000
	}
}
