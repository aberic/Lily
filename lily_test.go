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
	"encoding/json"
	"github.com/aberic/gnomon"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	checkbookName = "test"
	shopperName   = "shop"
)

func TestLily_Restart(t *testing.T) {
	l := ObtainLily()
	l.Restart()
	for i := 807; i > 800; i-- {
		j, err := l.Get(checkbookName, shopperName, strconv.Itoa(i))
		if nil != err {
			t.Log(err)
		} else {
			t.Log("get ", i, " = ", j)
		}
	}
}

func TestLily_New(t *testing.T) {
	var nds []Nodal
	for i := 0; i < 100000; i++ {
		nd := &node{index: nil, nodes: []Nodal{}}
		nds = append(nds, nd)
	}
	t.Log("success", len(nds))
}

func TestLilyPutGet(t *testing.T) {
	l := ObtainLily()
	l.Start()
	hashKey, err := l.PutD("young", 101)
	t.Log("put 101 young | hashKey =", hashKey, "| err = ", err)
	i, err := l.GetD("young")
	t.Log("get 101 young =", i, "| err = ", err)
	hashKey, err = l.SetD("young", 102)
	t.Log("put 102 young | hashKey =", hashKey, "| err = ", err)
	i, err = l.GetD("young")
	t.Log("get 102 young =", i, "| err = ", err)
	hashKey, err = l.PutD("young", 103)
	t.Log("put 103 young | hashKey =", hashKey, "| err = ", err)
	i, err = l.GetD("young")
	t.Log("get 103 young =", i, "| err = ", err)
}

func TestPutGet(t *testing.T) {
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Log(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	if err := l.CreateKey(checkbookName, shopperName, "ID"); nil != err {
		t.Error(err)
	}
	if _, err = l.Put(checkbookName, shopperName, strconv.Itoa(198), 200); nil != err {
		t.Log(err)
	}
	if i, err := l.Get(checkbookName, shopperName, strconv.Itoa(198)); nil != err {
		t.Log(err)
	} else {
		t.Log("get 198 =", i, "err =", err)
	}
	if _, err = l.Set(checkbookName, shopperName, strconv.Itoa(198), 201); nil != err {
		t.Log(err)
	}
	if i, err := l.Get(checkbookName, shopperName, strconv.Itoa(198)); nil != err {
		t.Log(err)
	} else {
		t.Log("get 198 =", i, "err =", err)
	}
	if _, err := l.Put(checkbookName, shopperName, strconv.Itoa(198), 200); nil != err {
		t.Log(err)
	}
	if i, err := l.Get(checkbookName, shopperName, strconv.Itoa(198)); nil != err {
		t.Log(err)
	} else {
		t.Log("get 198 =", i, "err =", err)
	}
}

func TestPutGets(t *testing.T) {
	//gnomon.Log().Set(gnomon.Log().ErrorLevel(), false)
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Error(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	for i := 255; i > 0; i-- {
		_, _ = l.Put(checkbookName, shopperName, strconv.Itoa(i), i)
	}
	for i := 255; i > 0; i-- {
		j, err := l.Get(checkbookName, shopperName, strconv.Itoa(i))
		t.Log("get ", i, " = ", j, "err = ", err)
	}
}

func TestQuerySelector1(t *testing.T) {
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Error(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	for i := 1; i <= 10; i++ {
		_, _ = l.Put(checkbookName, shopperName, strconv.Itoa(i), i)
	}
	for i := 1; i <= 10; i++ {
		j, err := l.Get(checkbookName, shopperName, strconv.Itoa(i))
		t.Log("get ", i, " = ", j, "err = ", err)
	}
	_, i, err := l.Select(checkbookName, shopperName, &Selector{})
	t.Log("select = ", i, "err = ", err)
	_, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{}})
	t.Log("select = ", i, "err = ", err)
}

type TestValue struct {
	ID          int
	Age         int
	IsMarry     bool
	Timestamp   int64
	TestValueIn *TestValueIn
}

type TestValueIn struct {
	ID        int
	Age       int
	IsMarry   bool
	Timestamp int64
}

func TestQuerySelector2(t *testing.T) {
	//gnomon.Log().Set(gnomon.Log().ErrorLevel(), false)
	gnomon.Log().Debug("TestQuerySelector2 Start")
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Error(err)
	}
	if err = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc); nil != err {
		t.Error(err)
	}
	if err = l.CreateKey(checkbookName, shopperName, "TestValueIn.ID"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "TestValueIn.Timestamp"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "TestValueIn.Age"); nil != err {
		t.Error(err)
	}
	gnomon.Log().Debug("TestQuerySelector2 Put")
	var wg sync.WaitGroup
	for i := 1007; i > 0; i-- {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tv := &TestValue{
				ID:        i,
				Age:       i + 19,
				Timestamp: time.Now().Local().UnixNano(),
				TestValueIn: &TestValueIn{
					ID:        i,
					Age:       i + 19,
					Timestamp: time.Now().Local().UnixNano(),
				},
			}
			if _, err := l.Put(checkbookName, shopperName, strconv.Itoa(i), tv); nil != err {
				t.Log(err)
			}
		}(i)
	}
	wg.Wait()
	gnomon.Log().Debug("TestQuerySelector2 Get")
	for i := 807; i > 800; i-- {
		j, err := l.Get(checkbookName, shopperName, strconv.Itoa(i))
		if nil != err {
			t.Log(err)
		} else {
			t.Log("get ", i, " = ", j)
		}
	}
	gnomon.Log().Debug("TestQuerySelector2 Select")
	var (
		i     interface{}
		count int
	)

	count, i, err = l.Select(checkbookName, shopperName, &Selector{
		Skip:  1000,
		Limit: 5,
	})
	t.Log("select nil count =", count, "i = ", i, "err = ", err)

	count, i, err = l.Select(checkbookName, shopperName, &Selector{
		Conditions: []*condition{{Param: "TestValueIn.Timestamp", Cond: "gt", Value: 10}},
		Skip:       5,
		Limit:      5,
	})
	t.Log("select time count =", count, "i = ", i, "err = ", err)

	count, i, err = l.Select(checkbookName, shopperName, &Selector{
		Sort:  &sort{Param: "TestValueIn.Timestamp", ASC: true},
		Skip:  1000,
		Limit: 5,
	})
	t.Log("select time true count =", count, "i =", i, "err = ", err)

	count, i, err = l.Select(checkbookName, shopperName, &Selector{
		Sort:  &sort{Param: "TestValueIn.Timestamp", ASC: false},
		Skip:  1000,
		Limit: 5,
	})
	t.Log("select time false count =", count, "i = ", i, "err = ", err)

	count, i, err = l.Select(checkbookName, shopperName, &Selector{
		Conditions: []*condition{
			{Param: "TestValueIn.Age", Cond: "gt", Value: 911},
			{Param: "TestValueIn.ID", Cond: "gt", Value: 888},
			{Param: "TestValueIn.ID", Cond: "lt", Value: 898},
			{Param: "TestValueIn.ID", Cond: "dif", Value: 891},
		},
		Sort:  &sort{Param: "TestValueIn.ID", ASC: true},
		Skip:  1,
		Limit: 5,
	})
	t.Log("select id true count =", count, "i =", i, "err = ", err)

	count, i, err = l.Select(checkbookName, shopperName, &Selector{
		Conditions: []*condition{{Param: "TestValueIn.ID", Cond: "gt", Value: 999}},
		Sort:       &sort{Param: "TestValueIn.ID", ASC: false},
		Skip:       1,
		Limit:      5,
	})
	t.Log("select id false count =", count, "i = ", i, "err = ", err)
}

func TestQuerySelector3(t *testing.T) {
	//gnomon.Log().Set(gnomon.Log().ErrorLevel(), false)
	gnomon.Log().Debug("TestQuerySelector2 Start")
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Error(err)
	}
	if err = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc); nil != err {
		t.Error(err)
	}
	if err = l.CreateKey(checkbookName, shopperName, "TestValueIn.ID"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "TestValueIn.Timestamp"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "TestValueIn.Age"); nil != err {
		t.Error(err)
	}
	gnomon.Log().Debug("TestQuerySelector2 Put")
	var wg sync.WaitGroup
	for i := 1007; i > 0; i-- {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tv := &TestValue{
				ID:        i,
				Age:       i + 19,
				Timestamp: time.Now().Local().UnixNano(),
				TestValueIn: &TestValueIn{
					ID:        i,
					Age:       i + 19,
					Timestamp: time.Now().Local().UnixNano(),
				},
			}
			data, err := json.Marshal(tv)
			if nil != err {
				t.Error(err)
			}
			var dataMap map[string]interface{}
			if err = json.Unmarshal(data, &dataMap); nil != err {
				t.Error(err)
			}
			if _, err := l.Put(checkbookName, shopperName, strconv.Itoa(i), dataMap); nil != err {
				t.Log(err)
			}
		}(i)
	}
	wg.Wait()
	gnomon.Log().Debug("TestQuerySelector2 Get")
	for i := 807; i > 800; i-- {
		j, err := l.Get(checkbookName, shopperName, strconv.Itoa(i))
		if nil != err {
			t.Log(err)
		} else {
			t.Log("get ", i, " = ", j)
		}
	}
	gnomon.Log().Debug("TestQuerySelector2 Select")
	var (
		i     interface{}
		count int
	)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{})
	t.Log("select nil count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Conditions: []*condition{{Param: "TestValueIn.Timestamp", Cond: "gt", Value: 1}}})
	t.Log("select time count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "TestValueIn.Timestamp", ASC: true}})
	t.Log("select time true count =", count, "i =", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "TestValueIn.Timestamp", ASC: false}})
	t.Log("select time false count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "TestValueIn.ID", ASC: true}})
	t.Log("select id true count =", count, "i =", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "TestValueIn.ID", ASC: false}})
	t.Log("select id false count =", count, "i = ", i, "err = ", err)
}

func TestQuerySelector4(t *testing.T) {
	//gnomon.Log().Set(gnomon.Log().ErrorLevel(), false)
	gnomon.Log().Debug("TestQuerySelector3 Start")
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Error(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	if err = l.CreateKey(checkbookName, shopperName, "ID"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "Timestamp"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "Age"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "IsMarry"); nil != err {
		t.Error(err)
	}
	gnomon.Log().Debug("TestQuerySelector3 Put")
	for i := 17; i > 0; i-- {
		go func(i int) {
			if _, err := l.Put(checkbookName, shopperName, strconv.Itoa(i), &TestValue{ID: i, Age: rand.Intn(17) + 1, IsMarry: i%2 == 0, Timestamp: time.Now().Local().UnixNano()}); nil != err {
				t.Log(err)
			}
		}(i)
	}
	time.Sleep(1 * time.Second)
	gnomon.Log().Debug("TestQuerySelector3 Get")
	for i := 17; i > 0; i-- {
		j, err := l.Get(checkbookName, shopperName, strconv.Itoa(i))
		if nil != err {
			t.Error(err)
		} else {
			t.Log("get ", i, " = ", j)
		}
	}
	gnomon.Log().Debug("TestQuerySelector2 Select")
	var (
		i     interface{}
		count int
	)
	//count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "IsMarry", ASC: true}})
	//t.Log("select IsMarry true count =", count, "i =", i, "err = ", err)
	//count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "IsMarry", ASC: false}})
	//t.Log("select IsMarry false count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{
		Conditions: []*condition{
			{Param: "Age", Cond: "gt", Value: 15},
		},
		Sort: &sort{Param: "Age", ASC: true},
	})
	t.Log("select Age true count =", count, "i =", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "Age", ASC: false}})
	t.Log("select Age false count =", count, "i = ", i, "err = ", err)
}

func TestQuerySelector5(t *testing.T) {
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Error(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	//for i := 1; i <= 10; i++ {
	//	_ = database.InsertInt(formName, i, i+10)
	//}
	_, err = l.Put(checkbookName, shopperName, "1000", 1000)
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "100", 100)
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "110000", 110000)
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "1100", 1100)
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "10000", 10000)
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "1", 1)
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "10", 10)
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "110", 110)
	t.Log("err = ", err)
	_, i, err := l.Select(checkbookName, shopperName, &Selector{})
	t.Log("select = ", i, "err = ", err)
	_, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{}})
	t.Log("select = ", i, "err = ", err)
}

func TestQuerySelector6(t *testing.T) {
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName, "数据库描述")
	if nil != err {
		t.Error(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	//for i := 1; i <= 10; i++ {
	//	_ = database.InsertInt(formName, i, i+10)
	//}
	_, err = l.Put(checkbookName, shopperName, "1000", &TestValue{ID: 1000, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "100", &TestValue{ID: 100, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "110000", &TestValue{ID: 110000, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "1100", &TestValue{ID: 1100, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "10000", &TestValue{ID: 10000, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "1", &TestValue{ID: 1, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "10", &TestValue{ID: 10, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "110", &TestValue{ID: 110, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, i, err := l.Select(checkbookName, shopperName, &Selector{})
	t.Log("select = ", i.([]interface{}), "err = ", err)
	_, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{}})
	t.Log("select = ", i.([]interface{}), "err = ", err)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
