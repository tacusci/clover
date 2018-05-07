package utils

import (
	"runtime"
	"strings"
)

type EndianOrder int8

const (
	BigEndian    EndianOrder = 0
	LittleEndian EndianOrder = 1
)

//TranslatePath to translate file location paths cross OS
func TranslatePath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, ":/", ":\\", -1)
		path = strings.Replace(path, "/", "\\", -1)
	} else if runtime.GOOS == "darwin" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	return path
}

//ConvertBytesToUInt16 takes two byte values and converts it to a two byte long uint
func ConvertBytesToUInt16(byte1 byte, byte2 byte, endianOrder EndianOrder) uint16 {
	var resultInt uint16
	if endianOrder == BigEndian {
		resultInt |= uint16(byte1) << 8
		resultInt |= uint16(byte2)
	} else if endianOrder == LittleEndian {
		resultInt |= uint16(byte1)
		resultInt |= uint16(byte2) << 8
	}
	return resultInt
}

//ConvertBytesToUInt32 takes four byte values and converts them to a 4 byte long uint
func ConvertBytesToUInt32(byte1 byte, byte2 byte, byte3 byte, byte4 byte, endianOrder EndianOrder) uint32 {
	var resultInt uint32
	if endianOrder == BigEndian {
		resultInt |= uint32(byte1) << 24
		resultInt |= uint32(byte2) << 16
		resultInt |= uint32(byte3) << 8
		resultInt |= uint32(byte4)
	} else if endianOrder == LittleEndian {
		resultInt |= uint32(byte1)
		resultInt |= uint32(byte2) << 8
		resultInt |= uint32(byte3) << 16
		resultInt |= uint32(byte4) << 24
	}
	return resultInt
}
