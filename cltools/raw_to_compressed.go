package cltools

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

	if bytesRead < 7 {
		log.Fatal("Unable to read enough bytes for the header...")
		return
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	println(isTiffImage(header))
}

func isTiffImage(header []byte) bool {
	if len(header) >= 4 {
		var magicNum uint8
		magicNum |= uint8(header[2]) | uint8(header[3])

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
