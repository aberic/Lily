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
)

const defaultLily = "_default"

type Data struct {
	name   string
	lilies map[string]*lily
}

func NewData(name string) *Data {
	data := &Data{name: name, lilies: map[string]*lily{}}
	data.lilies[defaultLily] = newLily(defaultLily, "default data lily", data)
	return data
}

func (d *Data) createGroup(name, comment string) error {
	if nil == d {
		return errors.New("data had never been created")
	}
	d.lilies[name] = newLily(name, comment, d)
	return nil
}

func (d *Data) Put(key Key, value interface{}) error {
	if nil == d {
		return errors.New("data had never been created")
	}
	return d.lilies[defaultLily].put(key, hash(key), value)
}

func (d *Data) Get(key Key) (interface{}, error) {
	if nil == d {
		return nil, errors.New("data had never been created")
	}
	return d.lilies[defaultLily].get(key, hash(key))
}

func (d *Data) PutG(groupName string, key Key, value interface{}) error {
	if nil == d {
		return errors.New("data had never been created")
	}
	l := d.lilies[groupName]
	if nil == l || nil == l.cities {
		return errors.New(strings.Join([]string{"group is invalid with name ", groupName}, ""))
	}
	return l.put(key, hash(key), value)
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
