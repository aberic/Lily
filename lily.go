package Lily

import (
	"errors"
	"strings"
)

// lily The Shopper
//
// hash array 模型 [00, 01, 02, 03, 04, 05, 06, 07, 08, 09, a, b, c, d, e, f]
//
// b+tree 模型 degree=128;level=4;nodes=[degree^level]/(degree-1)=2113665;
//
// node 内范围控制数量 key=127
//
// tree 内范围控制数量 treeCount=nodes*key=268435455
//
// hash array 内范围控制数量 t*16=4294967280
//
// level1间隔 ld1=(treeCount+1)/128=2097152
//
// level2间隔 ld2=(16513*127+1)/128=16384
//
// level3间隔 ld3=(129*127+1)/128=128
//
// level4间隔 ld3=(1*127+1)/128=1
type lily struct {
	data   *Data
	name   string
	cities []*city
}

func (l *lily) Put(key Key, value interface{}) error {
	if nil == l || nil == l.cities {
		return errors.New("db is invalid")
	}
	return l.put(key, hash(key), value)
}

func (l *lily) Get(key Key) (interface{}, error) {
	return l.get(key, hash(key))
}

func (l *lily) put(originalKey Key, key uint32, value interface{}) error {
	realKey := key / lilyDistance
	l.createChild(realKey)
	return l.cities[realKey].put(originalKey, key-realKey*lilyDistance, value)
}

func (l *lily) get(originalKey Key, key uint32) (interface{}, error) {
	realKey := key / lilyDistance
	if l.existChild(realKey) {
		return l.cities[key].get(originalKey, realKey)
	} else {
		return nil, errors.New(strings.Join([]string{"lily key", string(originalKey), "is nil"}, " "))
	}
}

func (l *lily) existChild(index uint32) bool {
	if nil == l.cities[index] {
		return false
	}
	return true
}

func (l *lily) createChild(index uint32) {
	if !l.existChild(index) {
		malls := make([]*mall, cityCount)
		for i := 0; i < 64; i++ {
			malls = append(malls, nil)
		}
		l.cities[index] = &city{
			key:   int(index),
			lily:  l,
			malls: malls,
		}
	}
}

func newLily(name string, data *Data) *lily {
	lily := &lily{
		name:   name,
		data:   data,
		cities: make([]*city, cityCount),
	}
	for i := 0; i < lilyCount; i++ {
		lily.cities = append(lily.cities, nil)
	}
	return lily
}
