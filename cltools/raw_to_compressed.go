package cltools

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type endian int

const (
	bigEndian    endian = 0
	littleEndian endian = 1
)

type tiffHeaderData struct {
	endianOrder endian
	magicNum    uint16
	tiffOffset  uint32
}

//RunRtc runs the raw to compressed image conversion tool
func RunRtc(locationpath string, intputType string, outputType string) {
	if len(locationpath) == 0 || len(intputType) == 0 || len(outputType) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if isDir, _ := isDirectory(locationpath); isDir {
		fileInfos, err := ioutil.ReadDir(locationpath)

		if err != nil {
			log.Fatal(err)
			return
		}

		for i := range fileInfos {
			fileInfo := fileInfos[i]
			if !fileInfo.IsDir() && strings.Contains(strings.ToLower(fileInfo.Name()), strings.ToLower(intputType)) {
				imageFile, err := os.Open(locationpath + "\\" + fileInfo.Name())
				if err != nil {
					log.Fatal(err)
					return
				}
				defer imageFile.Close()
				parseAllImageMeta(imageFile)
			}
		}
	}
}

func parseAllImageMeta(file *os.File) {
	readHeader(file)
}

func readHeader(file *os.File) {
	header := make([]byte, 8)
	file.Seek(0, 0)
	bytesRead, err := file.Read(header)

	if bytesRead < 8 {
		log.Fatal("Unable to read enough bytes for the header...")
		return
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = getTiffData(header)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func getEdianOrder(header []byte) endian {
	if len(header) >= 4 {
		var endianFlag uint16
		//add the bits to the 2 byte int and shove them to the left to make room for the other bits
		endianFlag |= uint16(header[0]) << 8
		endianFlag |= uint16(header[1])
		if endianFlag == 19789 {
			return bigEndian
		} else if endianFlag == 18761 {
			return littleEndian
		}
	}
	return bigEndian
}

func getTiffData(header []byte) (*tiffHeaderData, error) {
	endianOrder := getEdianOrder(header)
	tiffData := new(tiffHeaderData)
	if len(header) >= 8 {

		var magicNum uint16
		if endianOrder == bigEndian {
			magicNum |= uint16(header[2]) | uint16(header[3])
		} else {
			magicNum |= uint16(header[3]) | uint16(header[2])
		}

		tiffData.magicNum = magicNum
		tiffData.endianOrder = endianOrder

		var tiffOffset uint32
		if tiffData.endianOrder == bigEndian {
			tiffOffset |= uint32(header[4]) << 24
			tiffOffset |= uint32(header[5]) << 16
			tiffOffset |= uint32(header[6]) << 8
			tiffOffset |= uint32(header[7])
		}

		tiffData.tiffOffset = tiffOffset
	} else {
		return tiffData, errors.New("Header incorrect length")
	}
	return tiffData, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
