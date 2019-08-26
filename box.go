package Lily

import (
	"errors"
	"strings"
)

// box 包装盒一类
//
// 这里面存放的是具体物品
type box struct {
	key    uint32
	purse  *purse // box 所属 purse
	things map[uint32]*thing
}

func (b *box) put(originalKey Key, key uint32, value interface{}) error {
	b.createChild(originalKey, key, value)
	//log.Self.Debug("box", log.Uint32("key", key), log.Reflect("value", value))
	return b.things[key].put(originalKey, key, value)
}

func (b *box) get(originalKey Key, key uint32) (interface{}, error) {
	if b.existChild(key) {
		return b.things[key].get(originalKey, key)
	} else {
		return nil, errors.New(strings.Join([]string{"box key", string(originalKey), "is nil"}, " "))
	}
}

func (b *box) existChild(key uint32) bool {
	return nil != b.things[key]
}

func (b *box) createChild(originalKey Key, key uint32, value interface{}) {
	if !b.existChild(key) {
		b.things[key] = &thing{
			box:         b,
			originalKey: originalKey,
			value:       value,
		}
	}
}
