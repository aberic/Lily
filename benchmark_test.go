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
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func BenchmarkInsert(b *testing.B) {
	l := ObtainLily()
	l.Start()
	have := false
	for _, db := range l.GetDatabases() {
		if db.getName() == checkbookName {
			have = true
		}
	}
	if !have {
		if _, err := l.CreateDatabase(checkbookName, "数据库描述"); nil != err {
			b.Error(err)
		}
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	for i := 1; i <= b.N; i++ {
		go func(formName string, i int) {
			//_, _ = database.InsertInt(formName, i, i+10)
			_, _ = l.Put(checkbookName, formName, strconv.Itoa(i), &TestValue{ID: i, Age: rand.Intn(17) + 1, IsMarry: i%2 == 0, Timestamp: time.Now().Local().UnixNano()})
		}(shopperName, i)
		//_, _ = database.InsertInt(formName, i, i+10)
	}
}

func BenchmarkFileWrite1G1(b *testing.B) {
	f, err := os.OpenFile(filepath.Join(obtainConf().DataDir, "a.txt"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		b.Error(err)
	}
	strIn := ""
	appendStr := "document"
	for i := 0; i < 128; i++ {
		strIn = strings.Join([]string{strIn, appendStr}, "")
	}
	//t.Log("128 mem complete")
	//for i := 0; i < 1024; i++ {
	//	strIn = strings.Join([]string{strIn, strIn}, "")
	//}
	//t.Log("1024 mem complete")
	for i := 0; i < 1048576; i++ {
		_, _ = f.WriteString(strIn)
	}
	b.Log("1024 disk complete")
	err = f.Close()
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkFileWrite1G2(b *testing.B) {
	strIn := ""
	appendStr := "document"
	for i := 0; i < 128; i++ {
		strIn = strings.Join([]string{strIn, appendStr}, "")
	}
	//t.Log("128 mem complete")
	//for i := 0; i < 1024; i++ {
	//	strIn = strings.Join([]string{strIn, strIn}, "")
	//}
	//t.Log("1024 mem complete")
	for i := 0; i < 1048576; i++ {
		f, err := os.OpenFile(filepath.Join(obtainConf().DataDir, "b.txt"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			b.Error(err)
		}
		_, _ = f.WriteString(strIn)
		err = f.Close()
		if err != nil {
			b.Error(err)
		}
	}
	b.Log("1024 disk complete")
}
