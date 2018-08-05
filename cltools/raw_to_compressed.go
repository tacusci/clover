package cltools

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/tacusci/logging"

	"github.com/tacusci/clover/img"
	"github.com/tacusci/clover/utils"
)

//RunRtc runs the raw to compressed image conversion tool
func RunRtc(timeStamp bool, locationpath string, outputDirectory string, inputType string, outputType string, showConversionOutput bool, overwrite bool, recursive bool, retainFolderStructure bool) {
	if len(locationpath) == 0 || len(inputType) == 0 || len(outputType) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Clover - Running Raw To Compressed tool...\n")

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
	imagesToConvertChan := make(chan img.TiffImage, 32)

	if isDir, err := isDirectory(locationpath); isDir {
		//file searching wait group
		var fswg sync.WaitGroup
		//images to convert wait group
		var icwg sync.WaitGroup
		//add a wait for the initial single call of 'findImagesInDir'
		fswg.Add(1)
		go findImagesInDir(&fswg, &imagesToConvertChan, &doneSearchingChan, locationpath, inputTypePrefixToMatch, inputType, recursive)
		//add a wait for the call of 'convertRawImagesToCompressed'
		icwg.Add(1)
		go convertRawImagesToCompressed(&icwg, &imagesToConvertChan, &doneSearchingChan, inputType, outputType, showConversionOutput, overwrite, retainFolderStructure, locationpath, outputDirectory, &convertedImageCount)
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
	close(imagesToConvertChan)
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

func findImagesInDir(wg *sync.WaitGroup, itcc *chan img.TiffImage, dsc *chan bool, locationPath string, inputTypePrefixToMatch string, inputType string, recursive bool) {
	defer wg.Done()
	files, err := ioutil.ReadDir(locationPath)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	for i := range files {
		file := files[i]
		if !file.IsDir() {
			if strings.HasSuffix(strings.ToLower(file.Name()), strings.ToLower(inputType)) {
				if inputTypePrefixToMatch != "*" {
					if !strings.Contains(file.Name(), inputTypePrefixToMatch) {
						continue
					}
				}
				image, err := os.Open(utils.TranslatePath(path.Join(locationPath, file.Name())))
				if err != nil {
					logging.Error(err.Error())
					continue
				}
				var ti img.TiffImage
				switch inputType {
				case ".nef":
					ti = &img.NefImage{
						img.RawImage{
							File: image,
						},
					}
				case ".cr2":
					ti = &img.Cr2Image{
						img.RawImage{
							File: image,
						},
					}
				}
				if ti != nil {
					*itcc <- ti
					*dsc <- false
				}
			}
		} else {
			if file.IsDir() && recursive {
				wg.Add(1)
				findImagesInDir(wg, itcc, dsc, utils.TranslatePath(path.Join(locationPath, file.Name())), inputTypePrefixToMatch, inputType, recursive)
			}
		}
	}
}

func convertRawImagesToCompressed(wg *sync.WaitGroup, itcc *chan img.TiffImage, dsc *chan bool, inputType string, outputType string, showConversionOutput bool, overwrite bool, retainFolderStructure bool, inputDirectory string, outputDirectory string, convertedImageCount *uint32) {
	for {
		if !<-*dsc {
			ri := <-*itcc
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

func convertToCompressed(ti img.TiffImage, inputType string, outputType string, showConversionOutput bool, overwrite bool, retainFolderStructure bool, inputDirectory string, outputDirectory string, convertedImageCount *uint32) {
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
		logging.InfoNoColor(fmt.Sprintf("Converting image %s to %s", ti.GetRawImage().File.Name(), outputType))
	}

	if _, err := os.Stat(outputPath); err == nil && !overwrite {
		if showConversionOutput {
			logging.Error(" [FAILED] (Output result file already exists.)")
		}
		return
	}

	var succussfullyConvertedImage bool
	var conversionError error
	switch strings.ToLower(outputType) {
	case ".jpg":
		conversionError = ti.ConvertToJPEG(outputPath)
	case ".png":
		conversionError = ti.ConvertToPNG(outputPath)
	default:
		if showConversionOutput {
			logging.Error(fmt.Sprintf("[FAILED] (Output type %s not recognised/supported.)", outputType))
		}
		succussfullyConvertedImage = false
	}
	if conversionError != nil {
		if showConversionOutput {
			logging.Error(fmt.Sprintf(" [FAILED] (%s)", conversionError.Error()))
		}
		succussfullyConvertedImage = false
	} else {
		if showConversionOutput {
			logging.Info(" [SUCCESS]")
		}
		succussfullyConvertedImage = true
	}

	if succussfullyConvertedImage {
		*convertedImageCount++
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if fileInfo == nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func createDirectoryIfNotExists(dir string) error {
	if isDir, err := isDirectory(dir); !isDir {
		if err != nil {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func parseInputOutputTypes(inputType string, outputType string, supportedInputTypes []string, supportOutputTypes []string) (string, string, error) {

	//if the input type is *.nef then don't filter on file name

	r := regexp.MustCompile("(\\w+|\\*)\\.(\\w+)")
	res := r.FindStringSubmatch(inputType)

	if len(res) == 0 {
		return "", "", fmt.Errorf("Input type %s format not recognised, make sure input type matches <*|filename>.<typeext>", inputType)
	}

	inputPrefix := res[1]
	inputType = "." + res[2]

	if !utils.SSliceContains(supportedInputTypes, inputType) {
		return "", "", fmt.Errorf("Input type %s not supported", inputType)
	}

	if !utils.SSliceContains(supportOutputTypes, outputType) {
		if len(outputType) > 0 {
			return "", "", fmt.Errorf("Output type %s not supported", outputType)
		} else {
			return inputPrefix, inputType, nil
		}
	}

	return inputPrefix, inputType, nil
}
