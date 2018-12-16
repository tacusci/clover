package cltools

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/tacusci/logging"
)

//RunSdc to run the storage device checker tool
func RunSdc(locationPath string, sizeToWrite int, skipFileIntegrityCheck bool, dontDeleteFiles bool) {
	if len(locationPath) == 0 || sizeToWrite == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	filesToWrite := sizeToWrite / 1024 / 1024

	var writtenFilesCount int

	var wg sync.WaitGroup

	for writtenFilesCount < filesToWrite {
		wg.Add(1)
		go writeFile(fmt.Sprintf("%s%s%d", locationPath, string(os.PathSeparator), writtenFilesCount), &wg)
		writtenFilesCount++
	}

	wg.Wait()
}

func writeFile(locationPath string, wg *sync.WaitGroup) {
	logging.Info(fmt.Sprintf("Written file %s", locationPath))
	wg.Done()
}
