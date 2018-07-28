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

func RunTee(ts bool, sdir string, odir string, itype string, showExportOutput bool, overwrite bool, recursive bool) {
	if len(sdir) == 0 || len(odir) == 0 || len(itype) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Clover - Running TIFF EXIF export tool...\n")

	var st time.Time
	if ts {
		st = time.Now()
	}

	err := createDirectoryIfNotExists(odir)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	supportedInputTypes := []string{".nef"}
	supportedOutputTypes := []string{".jpg", ".png"}

	inputTypePrefixToMatch, inputType, err := parseInputOutputTypes(itype, "", supportedInputTypes, supportedOutputTypes)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	doneSearchingChan := make(chan bool, 32)
	imagesToExportExifChan := make(chan img.TiffImage, 32)

	if isDir, err := isDirectory(sdir); isDir {
		//file searching wait group
		var fswg sync.WaitGroup
		//images to export EXIF wait group
		var ieewg sync.WaitGroup
		//add a wait for the initial single call of 'findImagesInDir'
		fswg.Add(1)
		go findImagesInDir(&fswg, &imagesToExportExifChan, &doneSearchingChan, sdir, inputTypePrefixToMatch, inputType, recursive)
		ieewg.Add(1)
		go exportRawImageEXIF(&ieewg, &imagesToExportExifChan, &doneSearchingChan, itype, showExportOutput, overwrite, recursive, sdir, odir)
		//main thread doesn't wait after firing these goroutines, so force it to
		//wait until the file searching thread has finished
		fswg.Wait()
		//then tell the image exif export goroutine that there's no more images coming to export exif's
		doneSearchingChan <- true
		//wait on the image exif export goroutine until it's finished with all images it's already been working on
		ieewg.Wait()
		//both worker goroutines have finished, main thread continues
	} else {
		if err != nil {
			logging.ErrorAndExit(err.Error())
		}
	}

	if ts {
		logging.Info(fmt.Sprintf("Time taken: %d ms", time.Since(st).Nanoseconds()/1000000))
	}
}

func exportRawImageEXIF(wg *sync.WaitGroup, iteec *chan img.TiffImage, dsc *chan bool, itype string, showExportOutput bool, overwrite bool, recursive bool, sdir string, odir string) {
	for {
		if !<-*dsc {
			ri := <-*iteec
			wg.Add(1)
			if ri != nil {
				exportRawEXIFExport(ri, itype, showExportOutput, overwrite, sdir, odir)
			}
			wg.Done()
		} else {
			wg.Done()
		}
	}
}

func exportRawEXIFExport(ti img.TiffImage, itype string, showExportOutput bool, overwrite bool, sdir string, odir string) {
	if ti == nil {
		return
	}

	if ti.GetRawImage().File == nil {
		return
	}

	defer ti.GetRawImage().File.Close()

	sb := strings.Builder{}
	sb.WriteString(strings.TrimRight(odir, string(os.PathSeparator)))

	var fileNameToAdd string
	fileNameToAdd = filepath.Base(ti.GetRawImage().File.Name())
	fileNameToAdd = strings.Replace(fileNameToAdd, itype, ".txt", 1)
	fileNameToAdd = strings.Replace(fileNameToAdd, strings.ToUpper(itype), ".txt", 1)

	sb.WriteString(fileNameToAdd)

	outputPath := utils.TranslatePath(sb.String())

	if _, err := os.Stat(outputPath); err == nil && !overwrite {
		if showExportOutput {
			logging.Error(" [FAILED] (Output EXIF export file already exists.)")
		}
		return
	}

	for i := 0; i < len(ti.GetRawImage().Ifds); i++ {
		logging.Info(string(ti.GetRawImage().Ifds[i].ImageMakeTag))
	}
}
