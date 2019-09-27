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
	"github.com/aberic/gnomon"
	"github.com/vmihailenco/msgpack"
	"io"
	"os"
	"strings"
	"sync"
)

// writeResult 数据存储结果
type writeResult struct {
	seekStartIndex int64  // 索引最终存储在文件中的起始位置
	seekStart      uint32 // 16位起始seek
	seekLast       int    // 8位持续seek
	err            error
}

// readResult 数据读取结果
type readResult struct {
	value interface{}
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
	md5Key := gnomon.CryptoHash().MD516(ib.getKey()) // hash(keyStructure) 会发生碰撞，因此这里存储md5结果进行反向验证
	// 写入11位key及16位md5后key
	appendStr := strings.Join([]string{gnomon.String().PrefixSupplementZero(gnomon.Scale().Uint64ToDDuoString(ib.getHashKey()), 11), md5Key}, "")
	//gnomon.Log().Debug("storeIndex",
	//	gnomon.Log().Field("appendStr", appendStr),
	//	gnomon.Log().Field("formIndexFilePath", ib.getFormIndexFilePath()),
	//	gnomon.Log().Field("seekStartIndex", ib.getLink().getSeekStartIndex()))
	defer func() {
		if nil != file {
			<-s.limitOpenFileChan
			_ = file.Close()
		}
	}()
	// 将获取到的索引存储位置传入。如果为0，则表示没有存储过；如果不为0，则覆盖旧的存储记录
	if file, err = s.openFile(ib.getFormIndexFilePath(), os.O_CREATE|os.O_RDWR); nil != err {
		gnomon.Log().Error("storeIndex", gnomon.Log().Err(err))
		return &writeResult{err: err}
	}
	var seekEnd int64
	//gnomon.Log().Debug("running", gnomon.Log().Field("type", "moldIndex"), gnomon.Log().Field("seekStartIndex", it.link.getSeekStartIndex()))
	if ib.getLink().getSeekStartIndex() == -1 {
		if seekEnd, err = file.Seek(0, io.SeekEnd); nil != err {
			gnomon.Log().Error("storeIndex", gnomon.Log().Err(err))
			return &writeResult{err: err}
		}
		//gnomon.Log().Debug("running", gnomon.Log().Field("it.link.seekStartIndex == -1", seekEnd))
	} else {
		if seekEnd, err = file.Seek(ib.getLink().getSeekStartIndex(), io.SeekStart); nil != err { // 寻址到原索引起始位置
			gnomon.Log().Error("storeIndex", gnomon.Log().Err(err))
			return &writeResult{err: err}
		}
		//gnomon.Log().Debug("running", gnomon.Log().Field("seekStartIndex", it.link.getSeekStartIndex()), gnomon.Log().Field("it.link.seekStartIndex != -1", seekEnd))
	}
	// 写入11位key及16位md5后key及5位起始seek和4位持续seek
	if _, err = file.WriteString(strings.Join([]string{appendStr,
		gnomon.String().PrefixSupplementZero(gnomon.Scale().Uint32ToDDuoString(wf.seekStart), 5),
		gnomon.String().PrefixSupplementZero(gnomon.Scale().IntToDDuoString(wf.seekLast), 4)}, "")); nil != err {
		//gnomon.Log().Error("running", gnomon.Log().Field("seekStartIndex", seekEnd), gnomon.Log().Err(err))
		return &writeResult{err: err}
	}
	//gnomon.Log().Debug("storeIndex", gnomon.Log().Field("ib.getKey()", ib.getKey()), gnomon.Log().Field("md516Key", md516Key), gnomon.Log().Field("seekStartIndex", wf.seekStartIndex))
	ib.getLink().setSeekStartIndex(seekEnd)
	ib.getLink().setMD5Key(md5Key)
	ib.getLink().setSeekStart(wf.seekStart)
	ib.getLink().setSeekLast(wf.seekLast)
	//gnomon.Log().Debug("running", gnomon.Log().Field("it.link.seekStartIndex", seekEnd), gnomon.Log().Err(err))
	return &writeResult{
		seekStartIndex: seekEnd,
		seekStart:      wf.seekStart,
		seekLast:       wf.seekLast,
		err:            err}
}

func (s *storage) storeData(path string, value interface{}) *writeResult {
	var (
		file      *os.File
		seekStart int64
		seekLast  int
		data      []byte
		err       error
	)
	if data, err = msgpack.Marshal(value); nil != err {
		return &writeResult{err: err}
	}
	defer func() {
		if nil != file {
			<-s.limitOpenFileChan
			_ = file.Close()
		}
	}()
	if file, err = s.openFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND); nil != err {
		gnomon.Log().Error("storeData", gnomon.Log().Err(err))
		return &writeResult{err: err}
	}
	seekStart, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		gnomon.Log().Debug("storeData", gnomon.Log().Err(err))
		return &writeResult{err: err}
	}
	if seekLast, err = file.Write(data); nil != err {
		gnomon.Log().Debug("storeData", gnomon.Log().Err(err))
		return &writeResult{err: err}
	}
	return &writeResult{
		seekStart: uint32(seekStart),
		seekLast:  seekLast,
		err:       err,
	}
}

func (s *storage) read(filePath string, seekStart uint32, seekLast int, rr chan *readResult) {
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
	//gnomon.Log().Debug("read", gnomon.Log().Field("filePath", filePath), gnomon.Log().Field("seekStart", seekStart), gnomon.Log().Field("seekLast", seekLast))
	file, err = s.openFile(filePath, os.O_RDONLY)
	if err != nil {
		gnomon.Log().Error("read", gnomon.Log().Err(err))
		rr <- &readResult{err: err}
		return
	}
	_, err = file.Seek(int64(seekStart), io.SeekStart) //表示文件的起始位置，从第seekStart个字符往后读取
	if err != nil {
		gnomon.Log().Error("read", gnomon.Log().Err(err))
		rr <- &readResult{err: err}
		return
	}
	inputReader := bufio.NewReader(file)
	var bytes []byte
	if bytes, err = inputReader.Peek(seekLast); nil != err {
		gnomon.Log().Error("read", gnomon.Log().Err(err))
		rr <- &readResult{err: err}
		return
	}
	var value interface{}
	if err = msgpack.Unmarshal(bytes, &value); nil != err {
		gnomon.Log().Error("read", gnomon.Log().Err(err))
		rr <- &readResult{err: err}
		return
	}
	rr <- &readResult{err: err, value: value}
}

func (s *storage) openFile(filePath string, flag int) (*os.File, error) {
	s.limitOpenFileChan <- 1
	return os.OpenFile(filePath, flag, 0644)
}
