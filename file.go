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
	"os"
	"path/filepath"
)

var fileExistErr = errors.New("checkbook already exist")

// mkDataPath 创建库存储目录
//
// 存储格式 {dataDir}/checkbook/{dataName}/{shopperName}/{shopperName}.dat/idx...
func mkDataPath(dataName string) (err error) {
	dataPath := filepath.Join(dataDir, "checkbook", dataName)
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
		return fileExistErr
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}

// rmDataPath 删除库存储目录
func rmDataPath(dataName string) (err error) {
	dataPath := filepath.Join(dataDir, "checkbook", dataName)
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
		return fileExistErr
	}
	return os.Remove(dataPath)
}

// mkFormPath 创建库存储目录
//
// 存储格式 {dataDir}/checkbook/{dataName}/{shopperName}/{shopperName}.dat/idx...
func mkFormPath(dataName, formName string) (err error) {
	dataPath := filepath.Join(dataDir, "checkbook", dataName, formName)
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
		return fileExistErr
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}

// rmFormPath 删除库存储目录
func rmFormPath(dataName, formName string) (err error) {
	dataPath := filepath.Join(dataDir, "checkbook", dataName, formName)
	if exist, err := pathExist(dataPath); nil != err {
		return err
	} else if exist {
		return fileExistErr
	}
	return os.Remove(dataPath)
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
