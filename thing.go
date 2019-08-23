package Lily

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type thing struct {
	box         *box
	originalKey Key
	value       interface{}
	lock        sync.RWMutex
}

func (t *thing) put(originalKey Key, key uint32, value interface{}) error {
	var (
		path string
	)
	path = filepath.Join(dataDir,
		t.box.purse.trolley.mall.city.lily.data.name,
		t.box.purse.trolley.mall.city.lily.name,
		strconv.Itoa(t.box.purse.trolley.mall.city.key),
		strconv.Itoa(t.box.purse.trolley.mall.key),
		strconv.Itoa(t.box.purse.trolley.key),
		strings.Join([]string{
			strconv.Itoa(t.box.purse.key),
			"_",
			strconv.Itoa(t.box.key),
			".dat"}, "",
		),
	)
	//log.Self.Debug("box", log.Uint32("key", key), log.Reflect("value", value), log.String("path", path))
	return errors.New(path)
}

func (t *thing) get(originalKey Key, key uint32) (interface{}, error) {
	return nil, nil
}
