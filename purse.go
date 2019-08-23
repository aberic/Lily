package Lily

import (
	"errors"
	"strings"
)

// purse 手提袋
//
// 这里面能存放很多个包装盒
//
// box 包装盒集合
type purse struct {
	key         int
	flexibleKey uint32
	trolley     *trolley // purse 所属 trolley
	box         []*box
}

func (p *purse) put(originalKey Key, key uint32, value interface{}) error {
	p.flexibleKey = p.trolley.flexibleKey - purseDistance*uint32(p.key)
	realKey := p.flexibleKey / boxDistance
	//log.Self.Debug("purse", log.Uint32("key", key), log.Uint32("realKey", realKey))
	p.createChild(realKey)
	return p.box[realKey].put(originalKey, key, value)
}

func (p *purse) get(originalKey Key, key uint32) (interface{}, error) {
	p.flexibleKey = p.trolley.flexibleKey - purseDistance*uint32(p.key)
	realKey := p.flexibleKey / boxDistance
	if p.existChild(realKey) {
		return p.box[realKey].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"purse key", string(originalKey), "is nil"}, " "))
	}
}

func (p *purse) existChild(index uint32) bool {
	if nil == p.box[index] {
		return false
	}
	return true
}

func (p *purse) createChild(index uint32) {
	if !p.existChild(index) {
		p.box[index] = &box{
			key:    int(index),
			purse:  p,
			things: map[uint32]*thing{},
		}
	}
}
