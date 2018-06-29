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

//ConvertBytesSliceToUInt16 takes a slice of two bytes and converts them to a uint16
func ConvertBytesSliceToUInt16(btc []byte, eo EndianOrder) uint16 {
	if len(btc) != 2 {
		return 0
	}
	return ConvertBytesToUInt16(btc[0], btc[1], eo)
}

//ConvertBytesToUInt16 takes two byte values and converts it to a uint16
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

//ConvertBytesSliceToUInt32 takes a slice of four bytes and converts them to a uint32
func ConvertBytesSliceToUInt32(btc []byte, eo EndianOrder) uint32 {
	if len(btc) != 4 {
		return 0
	}
	return ConvertBytesToUInt32(btc[0], btc[1], btc[2], btc[3], eo)
}

//ConvertBytesToUInt32 takes four byte values and converts them to a uint32
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

//ConvertBytesSliceToUInt64 takes a slice of eight bytes and converts them to a uint64
func ConvertBytesSliceToUInt64(btc []byte, eo EndianOrder) uint64 {
	if len(btc) != 8 {
		return 0
	}
	return ConvertBytesToUInt64(btc[0], btc[1], btc[2], btc[3], btc[4], btc[5], btc[6], btc[7], eo)
}

//ConvertBytesToUInt64 takes eight byte values and converts them to a uint64
func ConvertBytesToUInt64(byte1 byte, byte2 byte, byte3 byte, byte4 byte, byte5 byte, byte6 byte, byte7 byte, byte8 byte, endianOrder EndianOrder) uint64 {
	var resultInt uint64
	if endianOrder == BigEndian {
		resultInt |= uint64(byte1) << 56
		resultInt |= uint64(byte2) << 48
		resultInt |= uint64(byte3) << 40
		resultInt |= uint64(byte4) << 32
		resultInt |= uint64(byte5) << 24
		resultInt |= uint64(byte6) << 16
		resultInt |= uint64(byte7) << 8
		resultInt |= uint64(byte8)
	} else if endianOrder == LittleEndian {
		resultInt |= uint64(byte1)
		resultInt |= uint64(byte2) << 8
		resultInt |= uint64(byte3) << 16
		resultInt |= uint64(byte4) << 24
		resultInt |= uint64(byte5) << 32
		resultInt |= uint64(byte6) << 40
		resultInt |= uint64(byte7) << 48
		resultInt |= uint64(byte8) << 56
	}
	return resultInt
}

func ConvertBytesSliceToFloat32(btc []byte, eo EndianOrder) float32 {
	if len(btc) != 4 {
		return 0.0
	}
	return ConvertBytesToFloat32(btc[0], btc[1], btc[2], btc[3], eo)
}

func ConvertBytesToFloat32(byte1 byte, byte2 byte, byte3 byte, byte4 byte, endianOrder EndianOrder) float32 {
	return 0.0
}
