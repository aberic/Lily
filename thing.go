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
	"github.com/ennoo/rivet/utils/log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type thing struct {
	nodal       nodal // box 所属 purse
	originalKey Key
	value       interface{}
	lock        sync.RWMutex
}

func (t *thing) put(originalKey Key, key uint32, value interface{}) error {
	var (
		path string
	)
	l := t.nodal.getPreNodal().getPreNodal().getPreNodal().getPreNodal().getPreNodal().(*lily)
	path = filepath.Join(dataDir,
		l.data.name,
		l.name,
		strconv.Itoa(int(t.nodal.getPreNodal().getPreNodal().getPreNodal().getPreNodal().getDegreeIndex())),
		strconv.Itoa(int(t.nodal.getPreNodal().getPreNodal().getPreNodal().getDegreeIndex())),
		strconv.Itoa(int(t.nodal.getPreNodal().getPreNodal().getDegreeIndex())),
		strconv.Itoa(int(t.nodal.getPreNodal().getDegreeIndex())),
		strings.Join([]string{
			strconv.Itoa(int(t.nodal.getPreNodal().getDegreeIndex())),
			"_",
			strconv.Itoa(int(t.nodal.getDegreeIndex())),
			".dat"}, "",
		),
	)
	log.Self.Debug("box", log.Uint32("key", key), log.Reflect("value", value), log.String("path", path))
	t.originalKey = originalKey
	t.value = value
	return nil
}

func (t *thing) get(originalKey Key, key uint32) (interface{}, error) {
	if t.originalKey == originalKey {
		return t.value, nil
	}
	return nil, errors.New(strings.Join([]string{"had no value for key ", string(originalKey)}, ""))
}
