package cltools

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tacusci/logging"
)

//RunSdc to run the storage device checker tool
func RunSdc(locationPath string, sizeToWrite int, skipFileIntegrityCheck bool, dontDeleteFiles bool) {
	if len(locationPath) == 0 || sizeToWrite == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//maxiumum number of running goroutines
	sem := make(chan bool, 1000)

	filesToWrite := sizeToWrite / 1024 / 1024

	var writtenFilesCount int

	for writtenFilesCount < filesToWrite {
		sem <- true
		go writeFile(fmt.Sprintf("%s%scloverdata%d.bin", locationPath, string(os.PathSeparator), writtenFilesCount), writtenFilesCount, sem)
		writtenFilesCount++
		time.Sleep(2)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
		time.Sleep(2)
	}
}

func writeFile(locationPath string, fileCount int, sem chan bool) {
	defer func() { <-sem }()

	file, err := os.Create(locationPath)

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	rand.Seed(int64(fileCount))

	//bytesToWrite := make([]byte, 1024*1000)
	bytesToWrite := make([]byte, 1024*1000)

	for i := 0; i < len(bytesToWrite)/2; i++ {
		bytesToWrite[i] = 0
	}

	for i := len(bytesToWrite) / 2; i < len(bytesToWrite); i++ {
		bytesToWrite[i] = byte(rand.Intn(254))
	}

	_, err = file.Write(bytesToWrite)

	if err != nil {
		logging.Error(err.Error())
	}

	file.Close()
}
