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

var (
	databaseName = "checkbook"
	formName     = "shopper"
)

func BenchmarkInsert(b *testing.B) {
	l := ObtainLily()
	l.Start()
	data, err := l.CreateDatabase(databaseName)
	if nil != err {
		b.Error(err)
	}
	_ = data.createForm(formName, "", true)
	now := time.Now().UnixNano()
	for i := 1; i <= b.N; i++ {
		go func(formName string, i int) {
			//_, _ = checkbook.InsertInt(formName, i, i+10)
			_, _ = l.Insert(databaseName, formName, Key(i), i+10)
		}(formName, i)
		//_, _ = checkbook.InsertInt(formName, i, i+10)
	}
	b.Log("time =", (time.Now().UnixNano()-now)/1e6)
}
