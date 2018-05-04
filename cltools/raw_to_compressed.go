package cltools

import (
	"flag"
	"fmt"
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
				} else {
					defer imageFile.Close()
					parseAllImageMeta(imageFile)
				}
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

	if err != nil {
		log.Fatal(err)
		return
	}

	for i := 0; i < bytesRead; i++ {
		fmt.Printf("%b", header[i])
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
