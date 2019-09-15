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
	"bufio"
	"container/heap"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/aberic/gnomon"
	"github.com/ennoo/rivet/utils/log"
	"github.com/vmihailenco/msgpack"
	"io"
	"os"
	"path/filepath"
	s "sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
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
	t.Log("asd = ", hash("asd"))
	t.Log("2147483648 = ", hash("2147483648"))
	t.Log("2147483648 = ", hash("2147483648"))
	t.Log("2147483650 = ", hash("2147483650"))
	t.Log("2147483650 = ", hash("2147483650"))
}

func TestChan(t *testing.T) {
	chanTest := make(chan int, 3)
	go func() {
		time.Sleep(2 * time.Second)
		chanTest <- 1
		chanTest <- 2
		chanTest <- 3
	}()
	for i := 0; i < 3; i++ {
		x := <-chanTest
		log.Self.Debug("TestChan", log.Int("x", x))
	}
}

func TestCond(t *testing.T) {
	var locker = new(sync.Mutex)
	var cond = sync.NewCond(locker)
	for i := 0; i < 5; i++ {
		go func(index int) {
			cond.L.Lock() //获取锁
			cond.Wait()   // 等待通知  暂时阻塞
			fmt.Println("index: ", index)
			cond.L.Unlock() //释放锁
		}(i)
	}
	//time.Sleep(time.Second * 1)
	//cond.Signal()
	//time.Sleep(time.Second * 1)
	//cond.Signal()
	time.Sleep(time.Second * 2)
	fmt.Println("broadcast")
	cond.Broadcast() // 下发广播给所有等待的goroutine
	time.Sleep(time.Second * 2)
}

func TestUint32toFullState(t *testing.T) {
	var index uint32
	index = 97890417
	intIndex := int(index)
	pos := 0
	for index > 1 {
		index /= 10
		pos++
	}
	backZero := 10 - pos
	backZeroStr := strconv.Itoa(intIndex)
	t.Log("backZeroStr =", backZeroStr)
	for i := 0; i < backZero; i++ {
		backZeroStr = strings.Join([]string{"0", backZeroStr}, "")
	}
	t.Log("backZeroStr =", backZeroStr)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	checkbookName = "test"
	shopperName   = "shop"
)

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
	_, err := l.CreateDatabase(checkbookName)
	if nil != err {
		t.Log(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
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
	_, err := l.CreateDatabase(checkbookName)
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
	_, err := l.CreateDatabase(checkbookName)
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
	Id        int
	Age       int
	Timestamp int64
}

func TestQuerySelector2(t *testing.T) {
	//gnomon.Log().Set(gnomon.Log().ErrorLevel(), false)
	gnomon.Log().Debug("TestQuerySelector2 Start")
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName)
	if nil != err {
		t.Error(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	if err = l.CreateIndex(checkbookName, shopperName, "Id"); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, IndexTimestamp); nil != err {
		t.Error(err)
	}
	if err = l.CreateIndex(checkbookName, shopperName, "Age"); nil != err {
		t.Error(err)
	}
	gnomon.Log().Debug("TestQuerySelector2 Put")
	for i := 1007; i > 0; i-- {
		go func(i int) {
			if _, err := l.Put(checkbookName, shopperName, strconv.Itoa(i), &TestValue{Id: i, Age: i + 19, Timestamp: time.Now().Local().UnixNano()}); nil != err {
				t.Log(err)
			}
		}(i)
	}
	time.Sleep(10 * time.Second)
	gnomon.Log().Debug("TestQuerySelector2 Get")
	for i := 1007; i > 1000; i-- {
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
	count, i, err = l.Select(checkbookName, shopperName, &Selector{})
	t.Log("select nil count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Conditions: []*condition{{Param: "Timestamp", Cond: "gt", Value: 1}}})
	t.Log("select time count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Conditions: []*condition{{Param: "Timestamp", Cond: "gt", Value: 1}}, Sort: &sort{Param: "Id", ASC: true}})
	t.Log("select time id true count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Conditions: []*condition{{Param: "Timestamp", Cond: "gt", Value: 1}}, Sort: &sort{Param: "Id", ASC: false}})
	t.Log("select time id false count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "Id", ASC: true}})
	t.Log("select id true count =", count, "i =", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: "Id", ASC: false}})
	t.Log("select id false count =", count, "i = ", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: IndexTimestamp, ASC: true}})
	t.Log("select time true count =", count, "i =", i, "err = ", err)
	count, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{Param: IndexTimestamp, ASC: false}})
	t.Log("select time false = ", i, "err = ", err)
}

func TestQuerySelector3(t *testing.T) {
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName)
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

func TestQuerySelector4(t *testing.T) {
	l := ObtainLily()
	l.Start()
	_, err := l.CreateDatabase(checkbookName)
	if nil != err {
		t.Error(err)
	}
	_ = l.CreateForm(checkbookName, shopperName, "", FormTypeDoc)
	//for i := 1; i <= 10; i++ {
	//	_ = database.InsertInt(formName, i, i+10)
	//}
	_, err = l.Put(checkbookName, shopperName, "1000", &TestValue{Id: 1000, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "100", &TestValue{Id: 100, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "110000", &TestValue{Id: 110000, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "1100", &TestValue{Id: 1100, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "10000", &TestValue{Id: 10000, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "1", &TestValue{Id: 1, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "10", &TestValue{Id: 10, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, err = l.Put(checkbookName, shopperName, "110", &TestValue{Id: 110, Timestamp: time.Now().Local().Unix()})
	t.Log("err = ", err)
	_, i, err := l.Select(checkbookName, shopperName, &Selector{})
	t.Log("select = ", i.([]interface{}), "err = ", err)
	_, i, err = l.Select(checkbookName, shopperName, &Selector{Sort: &sort{}})
	t.Log("select = ", i.([]interface{}), "err = ", err)
}

func TestPrint(t *testing.T) {
	for i := 0; i < 100000; i++ {
		//log.Self.Debug("print", log.Int("i = ", i))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestMatch2String(t *testing.T) {
	s := Selector{}
	t.Log("1 -", s.match2String(100))
	t.Log("2 -", s.match2String(true))
	t.Log("3 -", s.match2String(false))
	t.Log("4 -", s.match2String("hello"))
	t.Log("5 -", s.match2String(100.0101))
	t.Log("6 -", s.match2String(100.010))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type submiter struct {
}

func (s *submiter) doing(key uint32, value interface{}) {
	fmt.Println("keyStructure =", key, "value =", value)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TestHeap(t *testing.T) {
	h := &IntHeap{100, 16, 4, 8, 70, 2, 36, 22, 5, 12}

	fmt.Println("\nHeap:")
	heap.Init(h)

	fmt.Printf("最小值: %d\n", (*h)[0])

	//for(Pop)依次输出最小值,则相当于执行了HeapSort
	fmt.Println("\nHeap sort:")
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	//增加一个新值,然后输出看看
	fmt.Println("\nPush(h, 3),然后输出堆看看:")
	heap.Push(h, 3)
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	fmt.Println("\n使用sort.Sort排序:")
	h2 := IntHeap{100, 16, 4, 8, 70, 2, 36, 22, 5, 12}
	s.Sort(h2)
	for _, v := range h2 {
		fmt.Printf("%d ", v)
	}
	fmt.Println("\n二次打印:")
	for _, v := range h2 {
		fmt.Printf("%d ", v)
	}
	fmt.Println("")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestList(t *testing.T) {
	// 生成队列
	l := list.New()
	// 入队, 压栈
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)
	//for l.Len() > 0 {
	//	fmt.Printf("%d ", l.Front())
	//}
	// 出队
	i1 := l.Front()
	l.Remove(i1)
	fmt.Println(i1.Value, "len = ", l.Len())
	// 出栈
	i4 := l.Back()
	l.Remove(i4)
	fmt.Println(i4.Value, "len = ", l.Len())
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var lock sync.RWMutex

func TestRLock1(t *testing.T) {
	go func() {
		defer lock.RUnlock()
		lock.RLock()
		for i := 1; i <= 10; i++ {
			log.Self.Debug("持有", log.Int("i = ", i))
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(1 * time.Second)
	go func() {
		defer lock.Unlock()
		lock.Lock()
		log.Self.Debug("获取锁")
	}()
	time.Sleep(12 * time.Second)
}

func TestRLock2(t *testing.T) {
	go func() {
		defer lock.RUnlock()
		lock.RLock()
		for i := 1; i <= 10; i++ {
			log.Self.Debug("1持有", log.Int("i = ", i))
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		defer lock.RUnlock()
		lock.RLock()
		for i := 1; i <= 10; i++ {
			log.Self.Debug("2持有", log.Int("i = ", i))
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(1 * time.Second)
	go func() {
		defer lock.Unlock()
		lock.Lock()
		log.Self.Debug("获取锁")
	}()
	time.Sleep(12 * time.Second)
}

func TestLock(t *testing.T) {
	go func() {
		defer lock.Unlock()
		lock.Lock()
		for i := 1; i <= 10; i++ {
			log.Self.Debug("持有", log.Int("i = ", i))
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(1 * time.Second)
	go func() {
		defer lock.RUnlock()
		lock.RLock()
		log.Self.Debug("获取锁")
	}()
	time.Sleep(12 * time.Second)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestFileRandomWrite1(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(dataDir, "a.txt"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Error(err)
	}
	seek1, err := f.WriteString("document")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek1 = ", seek1)
	seek2, err := f.Seek(-3, io.SeekCurrent) //表示文件的起始位置，从第二个字符往后写入。
	if err != nil {
		t.Error(err)
	}
	t.Log("seek2 = ", seek2)
	seek3, err := f.WriteString("$$$$")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek3 = ", seek3)
	err = f.Close()
	if err != nil {
		t.Error(err)
	}
}
func TestFileRandomWrite2(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(dataDir, "a.txt"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Error(err)
	}
	seekCurrent, err := f.Seek(0, io.SeekCurrent)
	t.Log("seekCurrent = ", seekCurrent)
	seekEnd, err := f.Seek(0, io.SeekEnd)
	t.Log("seekEnd = ", seekEnd)
	seek, err := f.Seek(100, io.SeekStart)
	t.Log("seek = ", seek)
	seek1, err := f.WriteString("document")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek1 = ", seek1)
	seek1Current, err := f.Seek(0, io.SeekCurrent)
	t.Log("seek1Current = ", seek1Current)
	seek1End, err := f.Seek(0, io.SeekEnd)
	t.Log("seek1End = ", seek1End)
	//seek2, err := f.Seek(-3, io.SeekCurrent) //表示文件的起始位置，从第二个字符往后写入。
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log("seek2 = ", seek2)
	//seek3, err := f.WriteString("$$$$")
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log("seek3 = ", seek3)
	//err = f.Close()
	//if err != nil {
	//	t.Error(err)
	//}
}

func TestFileAppendWrite(t *testing.T) {
	var (
		f       *os.File
		seekNow int64
		err     error
	)
	f, err = os.OpenFile(filepath.Join(dataDir, "b.txt"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		t.Error(err)
	}
	seekNow, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		t.Error(err)
	}
	t.Log("seekNow = ", seekNow)
	seek1, err := f.WriteString("document")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek1 = ", seek1)

	seekNow, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		t.Error(err)
	}
	t.Log("seekNow = ", seekNow)
	seek2, err := f.WriteString("$$$$")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek2 = ", seek2)

	seekNow, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		t.Error(err)
	}
	t.Log("seekNow = ", seekNow)
	seek3, err := f.WriteString("document")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek3 = ", seek3)
	err = f.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestFileWrite1G(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(dataDir, "a.txt"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		t.Error(err)
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
	t.Log("1024 disk complete")
	err = f.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestFileRandomRead(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(dataDir, "b.txt"), os.O_RDONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	seek, err := f.Seek(1048577, io.SeekStart) //表示文件的起始位置，从第二个字符往后写入。
	if err != nil {
		t.Error(err)
	}
	t.Log("seek = ", seek)
	inputReader := bufio.NewReader(f)
	bs, err := inputReader.Peek(8)
	t.Log("bs string = ", string(bs))
}

func TestFileWriteInt(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(dataDir, "c.txt"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		t.Error(err)
	}
	seek1, err := f.WriteString("document")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek1 = ", seek1)
	seek2, err := f.WriteString("$$$$")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek2 = ", seek2)
	seek3, err := f.WriteString("document")
	if err != nil {
		t.Error(err)
	}
	t.Log("seek3 = ", seek3)
	err = f.Close()
	if err != nil {
		t.Error(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestMap(t *testing.T) {
	var mp = make(map[string]int)
	t.Log("len =", len(mp))
	mp["yes"] = 1
	t.Log("len =", len(mp))
	mp["no"] = 2
	t.Log("len =", len(mp))
	mp["ok"] = 3
	t.Log("len =", len(mp))
	delete(mp, "no")
	t.Log("len =", len(mp))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var is []interface{}

func TestShell(t *testing.T) {
	for i := 10; i > 0; i-- {
		tv := &TestValue{Id: i, Timestamp: time.Now().Local().UnixNano()}
		b, err := msgpack.Marshal(tv)
		if err != nil {
			panic(err)
		}
		var inter interface{}
		err = msgpack.Unmarshal(b, &inter)
		if err != nil {
			panic(err)
		}
		is = append(is, inter)
	}
	t.Log(is)
	selector := &Selector{Sort: &sort{Param: "Id", ASC: true}}
	t.Log(selector.shellAsc(is))
	t.Log(selector.shellDesc(is))
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
