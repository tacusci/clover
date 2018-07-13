package cltools

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/tacusci/clover/img"
	"github.com/tacusci/clover/utils"
	"github.com/tacusci/logging"
)

func RunTee(timeStamp bool, locationpath string, outputDirectory string, inputType string, outputType string, showConversionOutput bool, overwrite bool, recursive bool, retainFolderStructure bool) {
	if len(locationpath) == 0 || len(inputType) == 0 || len(outputType) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Clover - Running Tiff EXIF Export tool...\n")

	var st time.Time
	if timeStamp {
		st = time.Now()
	}

	err := createDirectoryIfNotExists(outputDirectory)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	var convertedImageCount uint32
	supportedInputTypes := []string{".nef"}
	supportedOutputTypes := []string{".jpg", ".png"}

	inputTypePrefixToMatch, inputType, err := parseInputOutputTypes(inputType, outputType, supportedInputTypes, supportedOutputTypes)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	doneSearchingChan := make(chan bool, 32)
	imagesToParseChan := make(chan img.TiffImage, 32)

	if isDir, err := isDirectory(locationpath); isDir {
		//file searching wait group
		var fswg sync.WaitGroup
		//images to export EXIF group
		var icwg sync.WaitGroup
		//add a wait for the initial single call of 'findImagesInDir'
		fswg.Add(1)
		go findImagesInDir(&fswg, &imagesToParseChan, &doneSearchingChan, locationpath, inputTypePrefixToMatch, inputType, recursive)
		//add a wait for the call of 'convertRawImagesToCompressed'
		icwg.Add(1)
		go exportRawImagesExifsToFile(&icwg, &imagesToParseChan, &doneSearchingChan, inputType, outputType, showConversionOutput, overwrite, retainFolderStructure, locationpath, outputDirectory, &convertedImageCount)
		//main thread doesn't wait after firing these goroutines, so force it to
		//wait until the file searching thread has finished
		fswg.Wait()
		//then tell the image conversion goroutine that there's no more images coming to convert
		doneSearchingChan <- true
		//wait on the image conversion goroutine until it's finished converting all images it's already been working on
		icwg.Wait()
		//both worker goroutines have finished, main thread continues
	} else {
		if err != nil {
			logging.ErrorAndExit(err.Error())
		}
	}
	close(doneSearchingChan)
	close(imagesToParseChan)
	var plural string
	if convertedImageCount != 1 {
		plural = "s"
	} else {
		plural = ""
	}
	logging.Info(fmt.Sprintf("Successfully converted %d raw image%s", convertedImageCount, plural))
	if timeStamp {
		logging.Info(fmt.Sprintf("Time taken: %d ms", time.Since(st).Nanoseconds()/1000000))
	}
}

func exportRawImagesExifsToFile(wg *sync.WaitGroup, itee *chan img.TiffImage, dsc *chan bool, inputType string, outputType string, showConversionOutput bool, overwrite bool, retainFolderStructure bool, inputDirectory string, outputDirectory string, convertedImageCount *uint32) {
	for {
		if !<-*dsc {
			ri := <-*itee
			wg.Add(1)
			if ri != nil {
				convertToCompressed(ri, inputType, outputType, showConversionOutput, overwrite, retainFolderStructure, inputDirectory, outputDirectory, convertedImageCount)
			}
			wg.Done()
		} else {
			wg.Done()
		}
	}
}

func exportEXIFs(ti img.TiffImage, inputType string, outputType string, showConversionOutput bool, overwrite bool, retainFolderStructure bool, inputDirectory string, outputDirectory string, convertedImageCount *uint32) {
	if ti == nil {
		return
	}

	if ti.GetRawImage().File == nil {
		return
	}

	defer ti.GetRawImage().File.Close()

	sb := strings.Builder{}
	sb.WriteString(strings.TrimRight(outputDirectory, string(os.PathSeparator)))

	if retainFolderStructure {
		subDirToAdd := strings.Replace(ti.GetRawImage().File.Name(), inputDirectory, "", -1)
		subDirToAdd = strings.Replace(subDirToAdd, filepath.Base(ti.GetRawImage().File.Name()), "", -1)
		if subDirToAdd != string(os.PathSeparator) {
			sb.WriteString(string(os.PathSeparator))
		}
		sb.WriteString(subDirToAdd)
		if err := createDirectoryIfNotExists(sb.String()); err != nil {
			logging.Error(err.Error())
			return
		}
	}

	var fileNameToAdd string
	fileNameToAdd = filepath.Base(ti.GetRawImage().File.Name())
	fileNameToAdd = strings.Replace(fileNameToAdd, inputType, outputType, 1)
	fileNameToAdd = strings.Replace(fileNameToAdd, strings.ToUpper(inputType), strings.ToUpper(outputType), 1)

	if !retainFolderStructure {
		sb.WriteRune(os.PathSeparator)
	}

	sb.WriteString(fileNameToAdd)

	outputPath := utils.TranslatePath(sb.String())

	if showConversionOutput {
		fmt.Printf("Converting image %s to %s", ti.GetRawImage().File.Name(), outputType)
	}

	if _, err := os.Stat(outputPath); err == nil && !overwrite {
		if showConversionOutput {
			logging.Error(" [FAILED] (Output result file already exists.)")
		}
		return
	}

	ti.Load()

}
