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
	"testing"
	"time"
)

func BenchmarkInsert(b *testing.B) {
	lilyName := "lily"
	data := NewData("data")
	_ = data.CreateLily(lilyName, "", true)
	now := time.Now().UnixNano()
	for i := 1; i <= b.N; i++ {
		go func(lilyName string, i int) {
			//_, _ = data.InsertInt(lilyName, i, i+10)
			_, _ = data.Insert(lilyName, Key(i), i+10)
		}(lilyName, i)
		//_, _ = data.InsertInt(lilyName, i, i+10)
	}
	b.Log("time =", (time.Now().UnixNano()-now)/1e6)
}
