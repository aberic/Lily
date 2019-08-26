package Lily

import (
	"github.com/ennoo/rivet/utils/env"
	"hash/crc32"
)

type database interface {
	put(originalKey Key, key uint32, value interface{}) error
	get(originalKey Key, key uint32) (interface{}, error)
	existChild(index uint32) bool
	createChild(index uint32)
}

//const (
//	cityCount = 16
//	mallCount    = 128
//	trolleyCount = 128
//	purseCount   = 128
//	boxCount     = 128
//
//	// 最大存储数，超过次数一律做新值换算
//	//lilyMax      uint32 = 4294967280
//	cityDistance uint32 = 268435455
//	// mallDistance level1间隔 ld1=(treeCount+1)/128=2097152 128^3
//	mallDistance uint32 = 2097152
//	// trolleyDistance level2间隔 ld2=(16513*127+1)/128=16384 128^2
//	trolleyDistance uint32 = 16384
//	// purseDistance level3间隔 ld3=(129*127+1)/128=128 128^1
//	purseDistance uint32 = 128
//	// boxDistance level4间隔 ld3=(1*127+1)/128=1 128^0
//	boxDistance uint32 = 1
//
//	dataPath = "DATA_PATH"
//)

const (
	cityCount    = 1
	mallCount    = 4
	trolleyCount = 4
	purseCount   = 4
	boxCount     = 4

	cityDistance uint32 = 0
	// mallDistance level1间隔 ld1=(treeCount+1)/128=2097152 128^3
	mallDistance uint32 = 64
	// trolleyDistance level2间隔 ld2=(16513*127+1)/128=16384 128^2
	trolleyDistance uint32 = 16
	// purseDistance level3间隔 ld3=(129*127+1)/128=128 128^1
	purseDistance uint32 = 4
	// boxDistance level4间隔 ld3=(1*127+1)/128=1 128^0
	boxDistance uint32 = 1

	dataPath = "DATA_PATH"
)

var (
	dataDir string
)

type Key string

// String hashes a string to a unique hashcode.
func hash(key Key) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func init() {
	dataDir = env.GetEnv(dataPath)
}
