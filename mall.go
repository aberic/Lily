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
	key         uint32
	realKey     uint32
	flexibleKey uint32
	city        *city
	trolleys    []*trolley
}

func (m *mall) put(originalKey Key, key uint32, value interface{}) error {
	m.realKey = m.city.flexibleKey / trolleyDistance
	m.flexibleKey = m.city.flexibleKey - m.realKey*trolleyDistance
	//log.Self.Debug("city", log.Uint32("key", key), log.Uint32("realKey", realKey))
	m.createChild(m.realKey)
	return m.trolleys[m.realKey].put(originalKey, key, value)
}

func (m *mall) get(originalKey Key, key uint32) (interface{}, error) {
	realKey := key / mallDistance
	m.flexibleKey = key - realKey*mallDistance
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
		m.trolleys[index] = &trolley{
			key:    m.city.realKey*mallDistance + index*trolleyDistance + trolleyDistance,
			mall:   m,
			purses: make([]*purse, purseCount),
		}
	}
}
