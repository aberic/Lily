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
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"sync"
)

const (
	// level1Distance level1间隔 65536^3 = 281474976710656 | 测试 4^3 = 64
	level1Distance uint64 = 281474976710656
	// level2Distance level2间隔 65536^2 = 4294967296 | 测试 4^2 = 16
	level2Distance uint64 = 4294967296
	// level3Distance level3间隔 65536^1 = 65536 | 测试 4^1 = 4
	level3Distance uint64 = 65536
	// level4Distance level4间隔 65536^0 = 1 | 测试 4^0 = 1
	level4Distance uint64 = 1
)

var (
	Version      = "1.0" // 版本号
	confInstance *Conf
	onceConf     sync.Once
)

// YamlConf lily启动配置文件根项目
type YamlConf struct {
	Conf Conf `yaml:"conf"`
}

// YamlConf lily启动配置文件子项目
type Conf struct {
	Port                     string `yaml:"Port"`                     // Port 开放端口，便于其它应用访问
	RootDir                  string `yaml:"RootDir"`                  // RootDir Lily服务默认存储路径
	DataDir                  string `yaml:"DataDir"`                  // DataDir Lily服务数据默认存储路径
	LimitOpenFile            int    `yaml:"LimitOpenFile"`            // LimitOpenFile 限制打开文件描述符次数
	TLS                      bool   `yaml:"TLS"`                      // TLS 是否开启 TLS
	TLSServerKeyFile         string `yaml:"TLSServerKeyFile"`         // TLSServerKeyFile lily服务私钥
	TLSServerCertFile        string `yaml:"TLSServerCertFile"`        // TLSServerCertFile lily服务数字证书
	Limit                    bool   `yaml:"Limit"`                    // Limit 是否启用服务限流策略
	LimitMillisecond         int    `yaml:"LimitMillisecond"`         // LimitMillisecond 请求限定的时间段（毫秒）
	LimitCount               int    `yaml:"LimitCount"`               // LimitCount 请求限定的时间段内允许的请求次数
	LimitIntervalMillisecond int    `yaml:"LimitIntervalMillisecond"` // LimitIntervalMillisecond 请求允许的最小间隔时间（毫秒），0表示不限
	lilyFilePath             string // Lily重启引导文件地址
}

// ObtainConf 根据文件地址获取Config对象
func ObtainConf(filePath string) *Conf {
	onceConf.Do(func() {
		confInstance = &Conf{}
		if gnomon.String().IsNotEmpty(filePath) {
			if err := confInstance.yaml2Conf(filePath); nil != err {
				gnomon.Log().Panic("ObtainConf", gnomon.Log().Err(err))
			}
		}
		if _, err := confInstance.scanDefault(); nil != err {
			gnomon.Log().Panic("ObtainConf", gnomon.Log().Err(err))
		}
	})
	return confInstance
}

// obtainConf 根据文件地址获取Config对象
func obtainConf() *Conf {
	onceConf.Do(func() {
		confInstance = &Conf{}
		if _, err := confInstance.scanDefault(); nil != err {
			gnomon.Log().Panic("obtainConf", gnomon.Log().Err(err))
		}
	})
	return confInstance
}

// scanDefault 扫描填充默认值
func (c *Conf) scanDefault() (*Conf, error) {
	if gnomon.String().IsEmpty(c.Port) {
		c.Port = "19877"
	}
	if gnomon.String().IsEmpty(c.RootDir) {
		c.RootDir = "lilyDB"
	}
	if gnomon.String().IsEmpty(c.DataDir) {
		c.DataDir = filepath.Join(c.RootDir, "data")
	}
	c.lilyFilePath = filepath.Join(c.DataDir, "lily.sync")
	if c.LimitOpenFile < 1000 {
		c.LimitOpenFile = 10000
	}
	if c.TLS {
		if gnomon.String().IsEmpty(c.TLSServerKeyFile) || gnomon.String().IsEmpty(c.TLSServerCertFile) {
			return nil, errors.New("tls server key file or cert file is nil")
		}
	}
	if c.Limit {
		if c.LimitCount < 0 || c.LimitMillisecond < 0 {
			return nil, errors.New("limit count or millisecond can not be zero")
		}
	}
	return c, nil
}

// yaml2Conf YML转配置对象
func (c *Conf) yaml2Conf(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if nil != err {
		return err
	}
	ymlConf := YamlConf{}
	err = yaml.Unmarshal([]byte(data), &ymlConf)
	if err != nil {
		return err
	}
	c = &ymlConf.Conf
	return nil
}
