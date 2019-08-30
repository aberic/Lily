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
	"encoding/json"
	"errors"
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

var (
	moldTypeInvalidErr = errors.New("mold type invalid")
)

type task interface {
	getMold() int             // 获取当前任务类型
	getAppendContent() string // getAppendContent 获取追加入文件的内容
	getSeekStart() int64
	getSeekLast() int
	getChanResult() chan *writeResult
}

type writeResult struct {
	seekStart int64
	seekLast  int
	err       error
}

type indexTask struct {
	key    string
	result chan *writeResult
	accept chan *writeResult
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
				it := task.(*indexTask)
				wr := <-it.accept
				// 写入10位key及16位起始seek和16位终止seek
				_, err = f.file.WriteString(strings.Join([]string{task.getAppendContent(), int64ToHexString(wr.seekStart), intToHexString(wr.seekLast)}, ""))
				task.getChanResult() <- &writeResult{err: err}
				continue
			case moldForm:
				seekLast, err = f.file.WriteString(task.getAppendContent())
				task.getChanResult() <- &writeResult{
					seekStart: seekStart,
					seekLast:  seekLast,
					err:       err,
				}
			}
		case <-to.C:
			// todo 需要优雅关闭，暂时暴力操作
			_ = f.file.Close()
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

func (s *storage) appendIndex(node nodal, path, key string, wr chan *writeResult) *writeResult {
	return s.write(node, path, key, wr, moldIndex)
}

func (s *storage) appendForm(form Form, path string, value interface{}, wr chan *writeResult) *writeResult {
	if data, err := json.Marshal(value); nil != err {
		return &writeResult{err: err}
	} else {
		return s.write(form, path, string(data), wr, moldForm)
	}
}

func (s *storage) write(data data, filePath, appendStr string, wr chan *writeResult, mold uint8) *writeResult {
	var (
		fd  *filed
		err error
	)
	if fd, err = s.useFiled(data, filePath); nil != err {
		return &writeResult{err: err}
	}
	result := make(chan *writeResult, 1)
	switch mold {
	default:
		return &writeResult{err: moldTypeInvalidErr}
	case moldIndex:
		fd.tasks <- &indexTask{key: appendStr, result: result, accept: wr}
		return <-result
	case moldForm:
		fd.tasks <- &formTask{key: appendStr, result: result}
		wrOut := <-result
		if nil == wrOut.err {
			wr <- wrOut
		}
		return wrOut
	}
}

func (s *storage) useFiled(data data, filePath string) (fd *filed, err error) {
	if fd = s.files[filePath]; nil == fd {
		defer data.unLock()
		data.lock()
		if fd = s.files[filePath]; nil == fd {
			var f *os.File
			if f, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644); nil != err {
				return
			}
			fd = &filed{
				file:  f,
				tasks: make(chan task, 1000),
			}
			err = pool().submit(func() {
				fd.running()
			})
			return
		}
	}
	return
}
