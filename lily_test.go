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
	data := NewData("data", true)
	_ = data.createGroup(lilyName, "", true)
	for i := 1; i <= 255; i++ {
		//_ = tmpLily.Put(Key(strconv.Itoa(i)), i)
		_ = data.PutGInt(lilyName, i, i)
	}
	_ = data.PutGInt(lilyName, 1, 1)
}

func TestList(t *testing.T) {
}

func TestPutGet(t *testing.T) {
	lilyName := "lily"
	data := NewData("data", true)
	_ = data.createGroup(lilyName, "", true)
	_ = data.PutGInt(lilyName, 198, 200)
	i, err := data.GetGInt(lilyName, 198)
	t.Log("get 198 = ", i, "err = ", err)
}

func TestPutGetInts(t *testing.T) {
	lilyName := "lily"
	data := NewData("data", true)
	_ = data.createGroup(lilyName, "", true)
	for i := 1; i <= 255; i++ {
		_ = data.PutGInt(lilyName, i, i+10)
	}
	for i := 1; i <= 255; i++ {
		j, err := data.GetGInt(lilyName, i)
		t.Log("get ", i, " = ", j, "err = ", err)
	}
}

func TestPutGets(t *testing.T) {
	lilyName := "lily"
	data := NewData("data", true)
	_ = data.createGroup(lilyName, "", true)
	for i := 1; i <= 255; i++ {
		_ = data.PutG(lilyName, Key(i), i)
	}
	for i := 1; i <= 255; i++ {
		j, err := data.GetG(lilyName, Key(i))
		t.Log("get ", i, " = ", j, "err = ", err)
	}
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
