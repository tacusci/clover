package cltools

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/tacusci/clover/utils"
)

type endian int

const (
	bigEndian    endian = 0
	littleEndian endian = 1
)

const (
	subfileTypeTag               uint16 = 0xfe
	oldSubfileTypeTag            uint16 = 0xff
	imageWidthTag                uint16 = 0x0100
	imageHeightTag               uint16 = 0x0101
	bitsPerSampleTag             uint16 = 0x0102
	compressionTag               uint16 = 0x0103
	photoMetricInterpretationTag uint16 = 0x0106
	thresholdingTag              uint16 = 0x0107
	cellWidthTag                 uint16 = 0x0108
	cellLengthTag                uint16 = 0x0109
	fillOrderTag                 uint16 = 0x010a
	documentNameTag              uint16 = 0x010d
	imageDescriptionTag          uint16 = 0x010e
	makeTag                      uint16 = 0x010f
	modelTag                     uint16 = 0x0110
	stripOffsetsTag              uint16 = 0x0111
	orientationTag               uint16 = 0x0112
	samplesPerPixelTag           uint16 = 0x0115
	rowsPerStripTag              uint16 = 0x0116
	stripByteCountsTag           uint16 = 0x0117
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
				filename := path.Join(locationpath, fileInfo.Name())
				filename = utils.TranslatePath(filename)
				imageFile, err := os.Open(filename)
				if err != nil {
					log.Fatal(err)
					return
				}
				defer imageFile.Close()
				err = parseAllImageMeta(imageFile)

				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}
	}
}

func parseAllImageMeta(file *os.File) error {
	header, err := readHeader(file)
	imageTiffHeaderData := *new(tiffHeaderData)

	if err != nil {
		//return here before the next read + check because we always want the root cause error to bubble back up
		return err
	}

	imageTiffHeaderData, err = getTiffHeader(header)

	ifd0Data := readIfd(file, imageTiffHeaderData.tiffOffset, imageTiffHeaderData.endianOrder)

	for i := range ifd0Data {
		fmt.Printf("%#x ", ifd0Data[i])
	}

	return nil
}

func readIfd(file *os.File, ifdOffset uint32, endianOrder endian) []byte {
	ifdTagCountBytes := make([]byte, 2)
	file.Seek(int64(ifdOffset), os.SEEK_SET)
	file.Read(ifdTagCountBytes)

	var ifdTagCount uint16
	if endianOrder == bigEndian {
		ifdTagCount |= uint16(ifdTagCountBytes[0]) << 8
		ifdTagCount |= uint16(ifdTagCountBytes[1])
	} else if endianOrder == littleEndian {
		ifdTagCount |= uint16(ifdTagCountBytes[0])
		ifdTagCount |= uint16(ifdTagCountBytes[1]) << 8
	}

	//each IFD tag length is 12 bytes
	ifdData := make([]byte, ifdTagCount*12)
	file.Seek(int64(ifdOffset+2), os.SEEK_SET)
	file.Read(ifdData)

	return ifdData
}

func readHeader(file *os.File) ([]byte, error) {
	header := make([]byte, 8)
	file.Seek(0, 0)
	bytesRead, err := file.Read(header)

	if bytesRead < 8 {
		return header, errors.New("Unable to read full header")
	}

	if err != nil {
		return header, err
	}
	return header, nil
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

func getTiffHeader(header []byte) (tiffHeaderData, error) {
	tiffData := new(tiffHeaderData)
	tiffData.endianOrder = getEdianOrder(header)

	if len(header) >= 8 {

		var magicNum uint16
		if tiffData.endianOrder == bigEndian {
			magicNum |= uint16(header[2]) | uint16(header[3])
		} else if tiffData.endianOrder == littleEndian {
			magicNum |= uint16(header[3]) | uint16(header[2])
		}

		tiffData.magicNum = magicNum

		var tiffOffset uint32
		if tiffData.endianOrder == bigEndian {
			tiffOffset |= uint32(header[4]) << 24
			tiffOffset |= uint32(header[5]) << 16
			tiffOffset |= uint32(header[6]) << 8
			tiffOffset |= uint32(header[7])
		} else if tiffData.endianOrder == littleEndian {
			tiffOffset |= uint32(header[4])
			tiffOffset |= uint32(header[5]) << 8
			tiffOffset |= uint32(header[6]) << 16
			tiffOffset |= uint32(header[7]) << 24
		}

		tiffData.tiffOffset = tiffOffset
	} else {
		return *tiffData, errors.New("Header incorrect length")
	}
	return *tiffData, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
