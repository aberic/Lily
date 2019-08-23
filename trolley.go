package Lily

import (
	"errors"
	"strings"
)

// trolley 购物车
//
// 这里面能存放很多手提袋
//
// purses 手提袋集合
//
// b+tree 模型 degree=128;level=4;nodes=[d^l]/(d-1)=2113665;
//
// node 内范围控制数量 key=127
//
// tree 内范围控制数量 n*k=268435455
type trolley struct {
	key         int
	flexibleKey uint32
	mall        *mall
	purses      []*purse
}

func (t *trolley) put(originalKey Key, key uint32, value interface{}) error {
	t.flexibleKey = t.mall.flexibleKey - trolleyDistance*uint32(t.key)
	realKey := t.flexibleKey / purseDistance
	//log.Self.Debug("trolley", log.Uint32("key", key), log.Uint32("realKey", realKey))
	t.createChild(realKey)
	return t.purses[realKey].put(originalKey, key, value)
}

func (t *trolley) get(originalKey Key, key uint32) (interface{}, error) {
	t.flexibleKey = t.mall.flexibleKey - trolleyDistance*uint32(t.key)
	realKey := t.flexibleKey / purseDistance
	if t.existChild(realKey) {
		return t.purses[realKey].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"trolley key", string(originalKey), "is nil"}, " "))
	}
}

func (t *trolley) existChild(index uint32) bool {
	if nil == t.purses[index] {
		return false
	}
	return true
}

func (t *trolley) createChild(index uint32) {
	if !t.existChild(index) {
		box := make([]*box, boxCount)
		for i := 0; i < boxCount; i++ {
			box = append(box, nil)
		}
		t.purses[index] = &purse{
			key:     int(index),
			trolley: t,
			box:     box,
		}
	}
}
