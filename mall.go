package Lily

import (
	"errors"
	"strings"
)

// mall 商城
//
// 这里面能存放很多购物车，是商场提供给shopper的最大载物设备
//
// trolleys 购物车集合
//
// b+tree 模型 degree=128;level=4;nodes=[d^l]/(d-1)=2113665;
//
// node 内范围控制数量 key=127
//
// tree 内范围控制数量 n*k=268435455
type mall struct {
	key         int
	flexibleKey uint32
	city        *city
	trolleys    []*trolley
}

func (m *mall) put(originalKey Key, key uint32, value interface{}) error {
	m.flexibleKey = key - mallDistance*uint32(m.key)
	realKey := m.flexibleKey / trolleyDistance
	//log.Self.Debug("mall", log.Uint32("key", key), log.Uint32("realKey", realKey))
	m.createChild(realKey)
	return m.trolleys[realKey].put(originalKey, key, value)
}

func (m *mall) get(originalKey Key, key uint32) (interface{}, error) {
	m.flexibleKey = key - mallDistance*uint32(m.key)
	realKey := m.flexibleKey / trolleyDistance
	if m.existChild(realKey) {
		return m.trolleys[realKey].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"mall key", string(originalKey), "is nil"}, " "))
	}
}

func (m *mall) existChild(index uint32) bool {
	if nil == m.trolleys[index] {
		return false
	}
	return true
}

func (m *mall) createChild(index uint32) {
	if !m.existChild(index) {
		purses := make([]*purse, purseCount)
		for i := 0; i < purseCount; i++ {
			purses = append(purses, nil)
		}
		m.trolleys[index] = &trolley{
			key:    int(index),
			mall:   m,
			purses: purses,
		}
	}
}
