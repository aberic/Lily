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
	"github.com/ennoo/rivet/utils/log"
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
				log.Self.Debug("running", log.String("type", "moldIndex"))
				it := task.(*indexTask)
				// 写入5位key及16位md5后key及5位起始seek和4位持续seek
				_, err = f.file.WriteString(strings.Join([]string{task.getAppendContent(), uint32ToDDuoString(it.accept.seekStart), intToDDuoString(it.accept.seekLast)}, ""))
				log.Self.Debug("running", log.Error(err))
				task.getChanResult() <- &writeResult{
					seekStart: it.accept.seekStart,
					seekLast:  it.accept.seekLast,
					err:       err,
				}
			case moldForm:
				log.Self.Debug("running", log.String("type", "moldForm"))
				seekLast, err = f.file.WriteString(task.getAppendContent())
				log.Self.Debug("running", log.Error(err))
				task.getChanResult() <- &writeResult{
					seekStart: uint32(seekStart),
					seekLast:  seekLast,
					err:       err,
				}
			}
		case <-to.C:
			// todo 需要优雅关闭，暂时暴力操作
			log.Self.Debug("timeout")
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

func (s *storage) appendIndex(node Nodal, path, key string, wr *writeResult) *writeResult {
	log.Self.Debug("appendIndex", log.String("path", path))
	return s.writeIndex(node, path, key, wr)
}

func (s *storage) appendForm(form Form, path string, value interface{}) *writeResult {
	var (
		data []byte
		err  error
	)
	log.Self.Debug("appendForm", log.String("path", path))
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
	log.Self.Debug("index")
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
	log.Self.Debug("form")
	fd.tasks <- &formTask{key: appendStr, result: result}
	return <-result
}

func (s *storage) read(filePath string, seekStart uint32, seekLast int, rr chan *readResult) {
	log.Self.Debug("read", log.Uint32("seekStart", seekStart), log.Int("seekLast", seekLast))
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Self.Debug("read", log.Error(err))
		rr <- &readResult{err: err}
		return
	}
	_, err = f.Seek(int64(seekStart), io.SeekStart) //表示文件的起始位置，从第seekStart个字符往后读取
	if err != nil {
		log.Self.Debug("read", log.Error(err))
		rr <- &readResult{err: err}
		return
	}
	inputReader := bufio.NewReader(f)
	var bytes []byte
	if bytes, err = inputReader.Peek(seekLast); nil != err {
		log.Self.Debug("read", log.Error(err))
		rr <- &readResult{err: err}
		return
	}
	log.Self.Debug("read", log.String("Data", string(bytes)))
	var value interface{}
	if err = json.Unmarshal(bytes, &value); nil != err {
		log.Self.Debug("read", log.Error(err))
		rr <- &readResult{err: err}
		return
	}
	rr <- &readResult{err: err, value: value}
}

func (s *storage) useFiled(data Data, filePath string) (fd *filed, err error) {
	if fd = s.files[filePath]; nil == fd {
		defer data.unLock()
		data.lock()
		log.Self.Debug("useFiled", log.String("filePath", filePath))
		if fd = s.files[filePath]; nil == fd {
			var f *os.File
			if f, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644); nil != err {
				log.Self.Debug("useFiled", log.Error(err))
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
