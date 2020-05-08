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
	"bufio"
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/log"
	"github.com/vmihailenco/msgpack"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	// ErrValueType value type error
	ErrValueType = errors.New("value type error")
	// ErrValueInvalid value is invalid
	ErrValueInvalid = errors.New("value is invalid")
)

type valueData struct {
	K string      // key
	I bool        // 是否有效
	V interface{} // 存储数据
}

// writeResult 数据存储结果
type writeResult struct {
	seekStartIndex int64 // 索引最终存储在文件中的起始位置
	seekStart      int64 // 16位起始seek
	seekLast       int   // 8位持续seek
	err            error
}

// readResult 数据读取结果
type readResult struct {
	key   string      // key
	value interface{} // value
	err   error
}

var (
	stg         *storage
	onceStorage sync.Once
)

func store() *storage {
	onceStorage.Do(func() {
		if nil == stg {
			stg = &storage{limitOpenFileChan: make(chan int, obtainConf().LimitOpenFile)}
		}
	})
	return stg
}

type storage struct {
	limitOpenFileChan chan int // limitOpenFileChan 限制打开文件描述符次数
}

func (s *storage) storeIndex(ib IndexBack, wf *writeResult) *writeResult {
	var (
		file *os.File
		err  error
	)
	defer ib.getLocker().unLock()
	ib.getLocker().lock()
	md5Key := gnomon.HashMD516(ib.getKey()) // hash(keyStructure) 会发生碰撞，因此这里存储md5结果进行反向验证
	// 写入11位key及16位md5后key
	appendStr := strings.Join([]string{gnomon.StringPrefixSupplementZero(gnomon.ScaleUint64ToDDuoString(ib.getHashKey()), 11), md5Key}, "")
	//log.Debug("storeIndex",
	//	log.Field("appendStr", appendStr),
	//	log.Field("formIndexFilePath", ib.getFormIndexFilePath()),
	//	log.Field("seekStartIndex", ib.getLink().getSeekStartIndex()))
	defer func() {
		if nil != file {
			<-s.limitOpenFileChan
			_ = file.Close()
		}
	}()
	// 将获取到的索引存储位置传入。如果为0，则表示没有存储过；如果不为0，则覆盖旧的存储记录
	if file, err = s.openFile(ib.getFormIndexFilePath(), os.O_CREATE|os.O_RDWR); nil != err {
		log.Error("storeIndex", log.Err(err))
		return &writeResult{err: err}
	}
	var seekEnd int64
	//log.Debug("running", log.Field("type", "moldIndex"), log.Field("seekStartIndex", it.link.getSeekStartIndex()))
	if ib.getLink().getSeekStartIndex() == -1 {
		if seekEnd, err = file.Seek(0, io.SeekEnd); nil != err {
			log.Error("storeIndex", log.Err(err))
			return &writeResult{err: err}
		}
		//log.Debug("running", log.Field("it.link.seekStartIndex == -1", seekEnd))
	} else {
		if seekEnd, err = file.Seek(ib.getLink().getSeekStartIndex(), io.SeekStart); nil != err { // 寻址到原索引起始位置
			log.Error("storeIndex", log.Err(err))
			return &writeResult{err: err}
		}
		//log.Debug("running", log.Field("seekStartIndex", it.link.getSeekStartIndex()), log.Field("it.link.seekStartIndex != -1", seekEnd))
	}
	// 写入11位key及16位md5后key及5位起始seek和4位持续seek
	if _, err = file.WriteString(strings.Join([]string{appendStr,
		gnomon.StringPrefixSupplementZero(gnomon.ScaleInt64ToDDuoString(wf.seekStart), 11),
		gnomon.StringPrefixSupplementZero(gnomon.ScaleIntToDDuoString(wf.seekLast), 4)}, "")); nil != err {
		//log.Error("running", log.Field("seekStartIndex", seekEnd), log.Err(err))
		return &writeResult{err: err}
	}
	//log.Debug("storeIndex", log.Field("ib.getKey()", ib.getKey()), log.Field("md516Key", md516Key), log.Field("seekStartIndex", wf.seekStartIndex))
	ib.getLink().setSeekStartIndex(seekEnd)
	ib.getLink().setMD5Key(md5Key)
	ib.getLink().setSeekStart(wf.seekStart)
	ib.getLink().setSeekLast(wf.seekLast)
	//log.Debug("running", log.Field("it.link.seekStartIndex", seekEnd), log.Err(err))
	return &writeResult{
		seekStartIndex: seekEnd,
		seekStart:      wf.seekStart,
		seekLast:       wf.seekLast,
		err:            err}
}

// storeData 存储具体内容
//
// path 存储文件路径
//
// value 存储具体内容
//
// valid 存储有效性，如无效则表示改记录不可用，即删除
func (s *storage) storeData(key, path string, value interface{}, valid bool) *writeResult {
	var (
		file      *os.File
		seekStart int64
		seekLast  int
		data      []byte
		err       error
	)
	// 存储数据外包装数据属性
	vd := &valueData{K: key, I: valid, V: value}
	if data, err = msgpack.Marshal(vd); nil != err {
		return &writeResult{err: err}
	}
	defer func() {
		if nil != file {
			<-s.limitOpenFileChan
			_ = file.Close()
		}
	}()
	if file, err = s.openFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND); nil != err {
		log.Error("storeData", log.Err(err))
		return &writeResult{err: err}
	}
	seekStart, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Debug("storeData", log.Err(err))
		return &writeResult{err: err}
	}
	if seekLast, err = file.Write(data); nil != err {
		log.Debug("storeData", log.Err(err))
		return &writeResult{err: err}
	}
	return &writeResult{
		seekStart: seekStart,
		seekLast:  seekLast,
		err:       err,
	}
}

func (s *storage) read(filePath string, seekStart int64, seekLast int) *readResult {
	var (
		file *os.File
		err  error
	)
	defer func() {
		if nil != file {
			<-s.limitOpenFileChan
			_ = file.Close()
		}
	}()
	//log.Debug("read", log.Field("filePath", filePath), log.Field("seekStart", seekStart), log.Field("seekLast", seekLast))
	file, err = s.openFile(filePath, os.O_RDONLY)
	if err != nil {
		//log.Error("read", log.Err(err))
		return &readResult{err: err}
	}
	_, err = file.Seek(seekStart, io.SeekStart) //表示文件的起始位置，从第seekStart个字符往后读取
	if err != nil {
		//log.Error("read", log.Err(err))
		return &readResult{err: err}
	}
	inputReader := bufio.NewReader(file)
	var bytes []byte
	if bytes, err = inputReader.Peek(seekLast); nil != err {
		//log.Error("read", log.Err(err))
		return &readResult{err: err}
	}
	var value interface{}
	if err = msgpack.Unmarshal(bytes, &value); nil != err {
		//log.Error("read", log.Err(err))
		return &readResult{err: err}
	}
	switch value.(type) {
	default:
		return &readResult{err: ErrValueType}
	case map[string]interface{}:
		valueMap := value.(map[string]interface{})
		if valueMap["I"].(bool) {
			return &readResult{key: valueMap["K"].(string), value: valueMap["V"], err: err}
		}
		return &readResult{err: ErrValueInvalid}
	}
}

func (s *storage) openFile(filePath string, flag int) (*os.File, error) {
	s.limitOpenFileChan <- 1
	return os.OpenFile(filePath, flag, 0644)
}
