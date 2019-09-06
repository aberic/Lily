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
	"encoding/json"
	"github.com/aberic/gnomon"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	moldIndex = iota
	moldForm
)

// task 执行存储任务对象
//
// 兼容数据存储和索引存储
type task interface {
	getMold() int             // 获取当前任务类型
	getAppendContent() string // getAppendContent 获取追加入文件的内容
	getSeekStart() int64
	getSeekLast() int
	getChanResult() chan *writeResult
}

// writeResult 数据存储结果
type writeResult struct {
	seekStart uint32 // 16位起始seek
	seekLast  int    // 8位持续seek
	err       error
}

type readResult struct {
	value interface{}
	err   error
}

type indexTask struct {
	key    string
	result chan *writeResult
	accept *writeResult
}

func (i *indexTask) getMold() int                     { return moldIndex }
func (i *indexTask) getAppendContent() string         { return i.key }
func (i *indexTask) getChanResult() chan *writeResult { return i.result }
func (i *indexTask) getSeekStart() int64              { return 0 }
func (i *indexTask) getSeekLast() int                 { return 0 }

type formTask struct {
	key       string
	seekStart int64
	seekLast  int
	result    chan *writeResult
}

func (f *formTask) getMold() int                     { return moldForm }
func (f *formTask) getAppendContent() string         { return f.key }
func (f *formTask) getChanResult() chan *writeResult { return f.result }
func (f *formTask) getSeekStart() int64              { return f.seekStart }
func (f *formTask) getSeekLast() int                 { return f.seekLast }

type filed struct {
	file  *os.File
	tasks chan task
}

func (f *filed) running() {
	to := time.NewTimer(time.Second)
	for {
		select {
		case task := <-f.tasks:
			var (
				seekStart int64
				seekLast  int
				err       error
			)
			to.Reset(time.Second)
			seekStart, err = f.file.Seek(0, io.SeekEnd)
			if err != nil {
				task.getChanResult() <- &writeResult{err: err}
				continue
			}
			switch task.getMold() {
			case moldIndex:
				gnomon.Log().Debug("running", gnomon.LogField("type", "moldIndex"))
				it := task.(*indexTask)
				// 写入5位key及16位md5后key及5位起始seek和4位持续seek
				_, err = f.file.WriteString(strings.Join([]string{task.getAppendContent(),
					gnomon.String().PrefixSupplementZero(gnomon.Scale().Uint32ToDDuoString(it.accept.seekStart), 5),
					gnomon.String().PrefixSupplementZero(gnomon.Scale().IntToDDuoString(it.accept.seekLast), 4)}, ""))
				gnomon.Log().Debug("running", gnomon.LogErr(err))
				task.getChanResult() <- &writeResult{
					seekStart: it.accept.seekStart,
					seekLast:  it.accept.seekLast,
					err:       err,
				}
			case moldForm:
				gnomon.Log().Debug("running", gnomon.LogField("type", "moldForm"))
				seekLast, err = f.file.WriteString(task.getAppendContent())
				gnomon.Log().Debug("running", gnomon.LogErr(err))
				task.getChanResult() <- &writeResult{
					seekStart: uint32(seekStart),
					seekLast:  seekLast,
					err:       err,
				}
			}
		case <-to.C:
			// todo 需要优雅关闭，暂时暴力操作
			gnomon.Log().Debug("timeout")
			_ = f.file.Close()
			return
		}
	}
}

var (
	stg         *storage
	onceStorage sync.Once
)

func store() *storage {
	onceStorage.Do(func() {
		if nil == stg {
			stg = &storage{
				files: map[string]*filed{},
			}
		}
	})
	return stg
}

type storage struct {
	files map[string]*filed
}

func (s *storage) appendIndex(ib IndexBack, key string, wr *writeResult) *writeResult {
	gnomon.Log().Debug("appendIndex", gnomon.LogField("path", ib.getFormIndexFilePath()))
	return s.writeIndex(ib.getNodal(), ib.getFormIndexFilePath(), key, wr)
}

func (s *storage) appendForm(form Form, path string, value interface{}) *writeResult {
	var (
		data []byte
		err  error
	)
	gnomon.Log().Debug("appendForm", gnomon.LogField("path", path))
	if data, err = json.Marshal(value); nil != err {
		return &writeResult{err: err}
	}
	return s.writeForm(form, path, string(data))
}

func (s *storage) writeIndex(data Data, filePath, appendStr string, wr *writeResult) *writeResult {
	var (
		fd  *filed
		err error
	)
	if fd, err = s.useFiled(data, filePath); nil != err {
		return &writeResult{err: err}
	}
	result := make(chan *writeResult, 1)
	gnomon.Log().Debug("index")
	fd.tasks <- &indexTask{key: appendStr, result: result, accept: wr}
	return <-result
}

func (s *storage) writeForm(data Data, filePath, appendStr string) *writeResult {
	var (
		fd  *filed
		err error
	)
	if fd, err = s.useFiled(data, filePath); nil != err {
		return &writeResult{err: err}
	}
	result := make(chan *writeResult, 1)
	gnomon.Log().Debug("form")
	fd.tasks <- &formTask{key: appendStr, result: result}
	return <-result
}

func (s *storage) read(filePath string, seekStart uint32, seekLast int, rr chan *readResult) {
	gnomon.Log().Debug("read", gnomon.LogField("seekStart", seekStart), gnomon.LogField("seekLast", seekLast))
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		gnomon.Log().Debug("read", gnomon.LogErr(err))
		rr <- &readResult{err: err}
		return
	}
	_, err = f.Seek(int64(seekStart), io.SeekStart) //表示文件的起始位置，从第seekStart个字符往后读取
	if err != nil {
		gnomon.Log().Debug("read", gnomon.LogErr(err))
		rr <- &readResult{err: err}
		return
	}
	inputReader := bufio.NewReader(f)
	var bytes []byte
	if bytes, err = inputReader.Peek(seekLast); nil != err {
		gnomon.Log().Debug("read", gnomon.LogErr(err))
		rr <- &readResult{err: err}
		return
	}
	gnomon.Log().Debug("read", gnomon.LogField("Data", string(bytes)))
	var value interface{}
	if err = json.Unmarshal(bytes, &value); nil != err {
		gnomon.Log().Debug("read", gnomon.LogErr(err))
		rr <- &readResult{err: err}
		return
	}
	rr <- &readResult{err: err, value: value}
}

func (s *storage) useFiled(data Data, filePath string) (fd *filed, err error) {
	if fd = s.files[filePath]; nil == fd {
		defer data.unLock()
		data.lock()
		gnomon.Log().Debug("useFiled", gnomon.LogField("filePath", filePath))
		if fd = s.files[filePath]; nil == fd {
			var f *os.File
			if f, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644); nil != err {
				gnomon.Log().Debug("useFiled", gnomon.LogErr(err))
				return
			}
			fd = &filed{
				file:  f,
				tasks: make(chan task, 1000),
			}
			err = pool().submit(func() {
				fd.running()
			})
		}
	}
	return
}
