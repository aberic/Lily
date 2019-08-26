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
	var data *Data
	var tmpLily *lily
	data = &Data{
		name:   "data",
		lilies: map[string]*lily{},
	}
	tmpLily = newLily("lily", data)
	data.lilies[tmpLily.name] = tmpLily
	for i := 1; i <= 255; i++ {
		//_ = tmpLily.Put(Key(strconv.Itoa(i)), i)
		_ = tmpLily.PutInt(i, i)
	}
	_ = tmpLily.PutInt(1, 1)
}

func TestPrint(t *testing.T) {
	for i := 0; i < 100000; i++ {
		//log.Self.Debug("print", log.Int("i = ", i))
	}
}
