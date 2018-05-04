package cltools

import (
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

	if isTiffImage(header, getEdianOrder(header)) {
		log.Println("Image confirmed to be TIFF")
	} else {
		log.Fatal("Image not of type TIFF")
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

func isTiffImage(header []byte, endianOrder endian) bool {
	if len(header) >= 4 {
		var magicNum uint16
		if endianOrder == bigEndian {
			magicNum |= uint16(header[2]) | uint16(header[3])
		} else {
			magicNum |= uint16(header[3]) | uint16(header[2])
		}

		//a TIFF image's magic number is 42
		if magicNum == 42 {
			return true
		}
	}
	return false
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
