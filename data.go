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
	"strings"
	"sync"
)

const (
	defaultLily         = "_default"
	defaultSequenceLily = "_default_sequence"
)

type Data struct {
	name   string
	lilies map[string]*lily
}

func NewData(name string, sequence bool) *Data {
	data := &Data{name: name, lilies: map[string]*lily{}}
	data.lilies[defaultLily] = newLily(defaultLily, "default data lily", data)
	if sequence {
		data.lilies[defaultSequenceLily] = newLily(defaultSequenceLily, "default data lily", data)
	}
	return data
}

func (d *Data) createGroup(name, comment string, sequence bool) error {
	if nil == d {
		return errors.New("data had never been created")
	}
	d.lilies[name] = newLily(name, comment, d)
	if sequence {
		sequenceName := sequenceName(name)
		d.lilies[sequenceName] = newLily(sequenceName, comment, d)
	}
	return nil
}

func (d *Data) Put(key Key, value interface{}) error {
	return d.PutG(defaultLily, key, value)
}

func (d *Data) Get(key Key) (interface{}, error) {
	return d.GetG(defaultLily, key)
}

func (d *Data) PutG(groupName string, key Key, value interface{}) error {
	if nil == d {
		return errors.New("data had never been created")
	}
	l := d.lilies[groupName]
	if nil == l || nil == l.cities {
		return errors.New(strings.Join([]string{"group is invalid with name ", groupName}, ""))
	}
	sequenceName := sequenceName(groupName)
	if nil == d.lilies[sequenceName] {
		return l.put(key, hash(key), value)
	} else {
		var (
			ls       *lily
			wg       sync.WaitGroup
			checkErr chan error
		)
		ls = d.lilies[sequenceName]
		checkErr = make(chan error, 2)
		wg.Add(2)
		go func(key Key, value interface{}) {
			defer wg.Done()
			err := l.put(key, hash(key), value)
			if nil != err {
				checkErr <- err
			}
		}(key, value)
		go func(key Key, value interface{}) {
			defer wg.Done()
			err := ls.put(key, hash(key), value)
			if nil != err {
				checkErr <- err
			}
		}(key, value)
		wg.Wait()
		return <-checkErr
	}
}

func (d *Data) GetG(groupName string, key Key) (interface{}, error) {
	if nil == d {
		return nil, errors.New("data had never been created")
	}
	l := d.lilies[groupName]
	if nil == l || nil == l.cities {
		return nil, errors.New(strings.Join([]string{"group is invalid with name ", groupName}, ""))
	}
	return l.get(key, hash(key))
}

func (d *Data) PutGInt(groupName string, key int, value interface{}) error {
	l := d.lilies[groupName]
	if nil == l || nil == l.cities {
		return errors.New(strings.Join([]string{"group is invalid with name ", groupName}, ""))
	}
	return l.put(Key(key), uint32(key), value)
}

func (d *Data) GetGInt(groupName string, key int) (interface{}, error) {
	l := d.lilies[groupName]
	if nil == l || nil == l.cities {
		return nil, errors.New(strings.Join([]string{"group is invalid with name ", groupName}, ""))
	}
	return l.get(Key(key), uint32(key))
}

func sequenceName(name string) string {
	return strings.Join([]string{name, "sequence"}, "_")
}
