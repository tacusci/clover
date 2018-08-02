package cltools

import (
	"bytes"
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

	inputTypePrefixToMatch, itype, err := parseInputOutputTypes(itype, "", supportedInputTypes, supportedOutputTypes)
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
		go findImagesInDir(&fswg, &imagesToExportExifChan, &doneSearchingChan, sdir, inputTypePrefixToMatch, itype, recursive)
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
	sb.WriteRune(os.PathSeparator)

	var fileNameToAdd string
	fileNameToAdd = filepath.Base(ti.GetRawImage().File.Name())
	fileNameToAdd = strings.Replace(fileNameToAdd, itype, ".txt", 1)
	fileNameToAdd = strings.Replace(fileNameToAdd, strings.ToUpper(itype), ".txt", 1)

	sb.WriteString(fileNameToAdd)

	outputPath := utils.TranslatePath(sb.String())

	if showExportOutput {
		fmt.Printf("Exporting image %s EXIFs", ti.GetRawImage().File.Name())
	}

	if _, err := os.Stat(outputPath); err == nil && !overwrite {
		if showExportOutput {
			logging.Error(" [FAILED] (Output result file already exists.)")
		}
		return
	}

	err := ti.Load()
	if err != nil {
		logging.Error(fmt.Sprintf(" [FAILED] (%s)", err.Error()))
		return
	}

	sb.Reset()

	for index, ifd := range ti.GetRawImage().Ifds {
		sb.WriteString(fmt.Sprintf("--------- START IFD%d START ---------\n", index))

		if ifd.BitsPerSample != nil && len(ifd.BitsPerSample) > 0 {
			sb.WriteString(fmt.Sprintf("Bits per sample -> %b\n", ifd.BitsPerSample))
		}

		if ifd.ImageModelTag != nil && len(ifd.ImageModelTag) > 0 {
			sb.WriteString(tidiedStringForOutput("Camera model", ifd.ImageModelTag))
		}

		if ifd.ImageMakeTag != nil && len(ifd.ImageMakeTag) > 0 {
			sb.WriteString(tidiedStringForOutput("Camera make", ifd.ImageMakeTag))
		}

		sb.WriteString(fmt.Sprintf("--------- END IFD%d END  ---------\n\n", index))

		if ifd.GpsIFD != nil {
			sb.WriteString("--------- START GPS IFD ---------\n")

			if ifd.GpsIFD.GPSVersionID != nil {
				sb.WriteString(fmt.Sprintf("GPS Version -> %d\n", ifd.GpsIFD.GPSVersionID))
			}

			sb.WriteString(fmt.Sprintf("GPS Time -> %d\n", ifd.GpsIFD.GPSTimeStamp))

			if len(ifd.GpsIFD.GPSSatellites) > 0 {
				sb.WriteString(tidiedStringForOutput("GPS Satellites", []byte(ifd.GpsIFD.GPSSatellites)))
			}

			sb.WriteString("--------- END GPS IFD ---------\n\n")
		}
	}

	ofile, err := os.Create(outputPath)
	defer ofile.Close()
	if err != nil {
		if showExportOutput {
			logging.Error(fmt.Sprintf(" [FAILED] (%s)", err.Error()))
		}
		return
	}
	_, err = ofile.WriteString(sb.String())
	ofile.Sync()
	if err != nil {
		if showExportOutput {
			logging.Error(fmt.Sprintf(" [FAILED] (%s)", err.Error()))
		}
	} else {
		if showExportOutput {
			logging.Info(" [SUCCESS]")
		}
	}
}

func tidiedStringForOutput(dt string, b []byte) string {
	return fmt.Sprintf("%s -> %s\n", dt, bytes.Trim(b, "\x00"))
}
