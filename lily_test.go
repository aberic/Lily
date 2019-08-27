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
	s "sort"
	"testing"
)

func TestCalculate(t *testing.T) {
	t.Log("2^31 = ", 1<<31)
}

func TestJson(t *testing.T) {
	var x = 8
	var y int
	data, err := json.Marshal(x)
	t.Log("x = ", x)
	t.Log("err = ", err)
	err = json.Unmarshal(data, &y)
	t.Log("y = ", y)
	t.Log("err = ", err)
}

func TestHashCode(t *testing.T) {
	t.Log("1 = ", hash("1"))
	t.Log("a = ", hash("a"))
	t.Log("asd = ", hash("asd"))
	t.Log("asd = ", hash(Key("asd")))
	t.Log("2147483648 = ", hash(Key("2147483648")))
	t.Log("2147483648 = ", hash(Key("2147483648")))
	t.Log("2147483650 = ", hash(Key("2147483650")))
	t.Log("2147483650 = ", hash(Key("2147483650")))
}

func TestPut(t *testing.T) {
	lilyName := "lily"
	data := NewData("data")
	_ = data.CreateLily(lilyName, "", true)
	for i := 1; i <= 255; i++ {
		//_ = tmpLily.InsertD(Key(strconv.Itoa(i)), i)
		_, _ = data.InsertGInt(lilyName, i, i)
	}
	_, _ = data.InsertGInt(lilyName, 1, 1)
}

func TestList(t *testing.T) {
}

func TestPutGet(t *testing.T) {
	lilyName := "lily"
	data := NewData("data")
	_ = data.CreateLily(lilyName, "", true)
	_, _ = data.InsertGInt(lilyName, 198, 200)
	i, err := data.QueryGInt(lilyName, 198)
	t.Log("get 198 = ", i, "err = ", err)
}

func TestPutGetInts(t *testing.T) {
	lilyName := "lily"
	data := NewData("data")
	_ = data.CreateLily(lilyName, "", true)
	for i := 1; i <= 255; i++ {
		_, _ = data.InsertGInt(lilyName, i, i+10)
	}
	for i := 1; i <= 255; i++ {
		j, err := data.QueryGInt(lilyName, i)
		t.Log("get ", i, " = ", j, "err = ", err)
	}
}

func TestPutGets(t *testing.T) {
	lilyName := "lily"
	data := NewData("data")
	_ = data.CreateLily(lilyName, "", true)
	for i := 1; i <= 255; i++ {
		_, _ = data.Insert(lilyName, Key(i), i)
	}
	for i := 1; i <= 255; i++ {
		j, err := data.Query(lilyName, Key(i))
		t.Log("get ", i, " = ", j, "err = ", err)
	}
}

func TestQuerySelector(t *testing.T) {
	lilyName := "lily"
	data := NewData("data")
	_ = data.CreateLily(lilyName, "", true)
	//for i := 1; i <= 10; i++ {
	//	_ = data.InsertGInt(lilyName, i, i+10)
	//}
	var err error
	_, err = data.InsertGInt(lilyName, 1000, 1000)
	t.Log("err = ", err)
	_, err = data.InsertGInt(lilyName, 100, 100)
	t.Log("err = ", err)
	_, err = data.InsertGInt(lilyName, 110000, 110000)
	t.Log("err = ", err)
	_, err = data.InsertGInt(lilyName, 1100, 1100)
	t.Log("err = ", err)
	_, err = data.InsertGInt(lilyName, 10000, 10000)
	t.Log("err = ", err)
	_, err = data.InsertGInt(lilyName, 1, 1)
	t.Log("err = ", err)
	_, err = data.InsertGInt(lilyName, 10, 10)
	t.Log("err = ", err)
	_, err = data.InsertGInt(lilyName, 110, 110)
	t.Log("err = ", err)
	i, err := data.QuerySelector(lilyName, &Selector{})
	t.Log("get ", i, " = ", i, "err = ", err)
	i, err = data.QuerySelector(lilyName, &Selector{Indexes: &indexes{IndexArr: []*index{{param: "_id", order: 1}}}})
	t.Log("get ", i, " = ", i, "err = ", err)
}

func TestPrint(t *testing.T) {
	for i := 0; i < 100000; i++ {
		//log.Self.Debug("print", log.Int("i = ", i))
	}
}

func TestBinaryFind(t *testing.T) {
	index, err := binaryMatch(150, []uint8{0, 8, 19, 49, 63, 80, 81, 98, 133, 150, 201, 250})
	t.Log("index = ", index, " | err = ", err)
}

func TestBinaryFind2(t *testing.T) {
	var (
		left   int
		middle int
		right  int
		is     []int
	)
	is = []int{0, 8, 19, 49, 163, 180, 281, 310, 333, 350, 401, 500}
	left = 0
	right = len(is) - 1
	query := 281
	for left <= right {
		middle = (left + right) / 2
		// 如果要找的数比midVal大
		if is[middle] > query {
			// 在arr数组的左边找
			right = middle - 1
		} else if is[middle] < query {
			// 在arr数组的右边找
			left = middle + 1
		} else if is[middle] == query {
			t.Log("找到下标", middle)
			break
		}
	}
}

func TestMatch2String(t *testing.T) {
	s := Selector{}
	t.Log("1 -", s.match2String(100))
	t.Log("2 -", s.match2String(true))
	t.Log("3 -", s.match2String(false))
	t.Log("4 -", s.match2String("hello"))
	t.Log("5 -", s.match2String(100.0101))
	t.Log("6 -", s.match2String(100.010))
}

type person struct {
	Name string
	Age  int
}
type personSlice []person

func (s personSlice) Len() int           { return len(s) }
func (s personSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s personSlice) Less(i, j int) bool { return s[i].Age < s[j].Age }

func TestSort(t *testing.T) {
	a := personSlice{
		{
			Name: "AAA",
			Age:  55,
		},
		{
			Name: "BBB",
			Age:  22,
		},
		{
			Name: "CCC",
			Age:  0,
		},
		{
			Name: "DDD",
			Age:  22,
		},
		{
			Name: "EEE",
			Age:  11,
		},
	}
	s.Stable(a)
	t.Log(a)
}
