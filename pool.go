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
	"github.com/panjf2000/ants"
	"sync"
)

var (
	dp   *dataPool
	once sync.Once
)

type dataPool struct {
	pool    *ants.Pool
	inserts chan insert
}

func pool() *dataPool {
	once.Do(func() {
		if nil == dp {
			p, _ := ants.NewPool(100)
			dp = &dataPool{
				pool:    p,
				inserts: make(chan insert, 1000),
			}
		}
	})
	return dp
}

type insert struct {
	data        nodal
	originalKey Key
	key         uint32
	value       interface{}
}

func (d *dataPool) insert(data nodal, originalKey Key, key uint32, value interface{}) error {
	return data.put(originalKey, key, value)
}

func (d *dataPool) submitInsert() error {
	for {
		select {
		case insert := <-d.inserts:
			insert.data.put(insert.originalKey, insert.key, insert.value)
		}
	}
}

// tune 动态变更协程池数量
func (d *dataPool) tune(poolSize int) {
	d.pool.Tune(poolSize)
}

func (d *dataPool) submit(task func()) error {
	return d.pool.Submit(func() {
		task()
	})
}

func (d *dataPool) submitInsert1(data nodal, originalKey Key, key uint32, value interface{}) error {
	return d.pool.Submit(func() {
		_ = data.put(originalKey, key, value)
	})
}
