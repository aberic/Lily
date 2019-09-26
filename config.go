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
	// level1Distance level1间隔 65536^3 = 281474976710656 | 测试 4^3 = 64
	level1Distance int64 = 281474976710656
	// level2Distance level2间隔 65536^2 = 4294967296 | 测试 4^2 = 16
	level2Distance int64 = 4294967296
	// level3Distance level3间隔 65536^1 = 65536 | 测试 4^1 = 4
	level3Distance int64 = 65536
	// level4Distance level4间隔 65536^0 = 1 | 测试 4^0 = 1
	level4Distance int64 = 1
)

var (
	rootDir       string // Lily服务默认存储路径
	dataDir       string // Lily服务默认存储路径
	lilyFilePath  string // Lily重启引导文件地址
	limitOpenFile int    // 限制打开文件描述符次数
)

func init() {
	var err error
	rootDir = gnomon.Env().GetD("DATA_PATH", "test/t1")
	dataDir = filepath.Join(rootDir, "Data")
	lilyFilePath = filepath.Join(dataDir, "lily.sync")
	if limitOpenFile, err = strconv.Atoi(gnomon.Env().GetD("LIMIT_COUNT", "10000")); nil == err {
		limitOpenFile = 10000
	}
}
