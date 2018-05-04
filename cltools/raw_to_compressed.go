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
		println("It is a directory")
		fileInfos, err := ioutil.ReadDir(locationpath)

		if err != nil {
			log.Fatal(err)
		} else {
			for i := range fileInfos {
				fileInfo := fileInfos[i]
				if !fileInfo.IsDir() && strings.Contains(strings.ToLower(fileInfo.Name()), strings.ToLower(intputType)) {
					imageFile, err := ioutil.ReadFile(locationpath + "\\" + fileInfo.Name())
					if err != nil {
						log.Fatal(err)
					} else {
						for i := 0; i < 8; i++ {
							fmt.Printf("%b", imageFile[i])
						}
					}
				}
			}
		}
	} else {
		println("It is not a directory")
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
