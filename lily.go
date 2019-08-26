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

func (l *lily) PutInt(key int, value interface{}) error {
	if nil == l || nil == l.cities {
		return errors.New("db is invalid")
	}
	return l.put(Key(key), uint32(key), value)
}

func (l *lily) GetInt(key int) (interface{}, error) {
	return l.get(Key(key), uint32(key))
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
	//realKey := key / cityDistance
	realKey := uint32(0)
	l.createChild(realKey)
	return l.cities[realKey].put(originalKey, key-realKey*cityDistance, value)
}

func (l *lily) get(originalKey Key, key uint32) (interface{}, error) {
	//realKey := key / cityDistance
	realKey := uint32(0)
	if l.existChild(realKey) {
		return l.cities[realKey].get(originalKey, key-realKey*cityDistance)
	} else {
		return nil, errors.New(strings.Join([]string{"lily key", string(originalKey), "is nil"}, " "))
	}
}

func (l *lily) existChild(index uint32) bool {
	return nil != l.cities[index]
}

func (l *lily) createChild(index uint32) {
	if !l.existChild(index) {
		l.cities[index] = &city{
			key:   index*cityDistance + cityDistance,
			lily:  l,
			malls: make([]*mall, mallCount),
		}
	}
}

func newLily(name string, data *Data) *lily {
	lily := &lily{
		name:   name,
		data:   data,
		cities: make([]*city, cityCount),
	}
	return lily
}
