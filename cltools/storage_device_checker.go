package cltools

import (
	"bufio"
	"crypto/md5"
	"flag"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/tacusci/clover/utils"
)

//RunSdc to run the storage device checker tool
func RunSdc(locationPath string, sizeToWrite int, skipFileIntegrityCheck bool, dontDeleteFiles bool) {
	if len(locationPath) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if sizeToWrite == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if sizeToWrite > 0 {
		fileCount, totalWrittenBytes, timeElapsed := writeDataToLocation(locationPath, sizeToWrite)

		var passed = false

		if !skipFileIntegrityCheck {
			passed = verify(fileCount, locationPath)
		}
		tidy(dontDeleteFiles, fileCount, locationPath)
		outputSummary(sizeToWrite, totalWrittenBytes, locationPath, passed, skipFileIntegrityCheck, timeElapsed)
	}
}

func writeDataToLocation(location string, size int) (int, int, time.Duration) {
	//bytes in 1MB
	var byteChunkSize = 1024 * 1000
	var totalWrittenBytes int
	var fileCount = 1

	startTime := time.Now()

	yColor := color.New(color.FgYellow)
	rColor := color.New(color.FgRed).Add(color.Bold)

	yColor.Printf("Running StorageDeviceChecker tool -> Writing %v bytes to %v\n", size, location)

	for {
		if totalWrittenBytes <= size-byteChunkSize {
			filename := path.Join(location, "cloverdata"+strconv.Itoa(fileCount)+".bin")
			filename = utils.TranslatePath(filename)
			file, err := os.Create(filename)
			check(err)
			defer file.Close()
			bufferedWriter := bufio.NewWriter(file)
			rand.Seed(int64(fileCount))
			bytesToWrite := make([]byte, 1024*1000)
			for i := len(bytesToWrite) / 2; i < len(bytesToWrite); i++ {
				bytesToWrite[i] = byte(rand.Intn(254))
			}
			fileMd5 := md5.Sum(bytesToWrite)
			for i := 0; i < len(fileMd5); i++ {
				bytesToWrite[i] = fileMd5[i]
			}
			bytesWritten, err := bufferedWriter.Write(bytesToWrite)
			bufferedWriter.Flush()
			if err != nil {
				rColor.Println("Unable to write more data")
				return fileCount, totalWrittenBytes, time.Now().Sub(startTime)
			}
			totalWrittenBytes += bytesWritten
			fileCount++
		} else {
			break
		}
	}
	return fileCount, totalWrittenBytes, time.Now().Sub(startTime)
}

func verify(fileCount int, location string) bool {

	rColor := color.New(color.FgRed).Add(color.Bold)

	for i := 1; i < fileCount; i++ {
		file, err := os.Open(path.Join(location, "cloverdata"+strconv.Itoa(i)+".bin"))
		if err != nil {
			rColor.Println("Unable to open " + file.Name() + " for verification...")
			break
		}
		defer file.Close()
		checksumBytes := make([]byte, 16)
		_, err = file.Seek(0, 0)
		check(err)
		_, err = file.Read(checksumBytes)
		check(err)
		fullFileBytes := make([]byte, 1024*1000)
		file.Seek(0, 0)
		file.Read(fullFileBytes)
		for i := 0; i < 17; i++ {
			fullFileBytes[i] = 0
		}
		checksumFromFileContents := md5.Sum(fullFileBytes)
		//if checksumFromFileContents[:] != fullFileBytes[:] {
		if compareHashBytes(checksumFromFileContents[:], fullFileBytes[:]) {
			rColor.Println("Incorrect checksum for file -> %v", file.Name())
			return false
		}
	}
	return true
}

func outputSummary(sizeToWrite int, totalWrittenBytes int, location string, verificationPassed bool, skipFileIntegrityCheck bool, timeElapsed time.Duration) {
	yColor := color.New(color.FgYellow)
	yBoldColor := color.New(color.FgYellow).Add(color.Bold)
	rColor := color.New(color.FgRed)
	gColor := color.New(color.FgGreen)
	yBoldColor.Println("------------- Summary -------------")
	yColor.Printf("Run for %v seconds...\n", timeElapsed.Seconds())

	writtenPercentage := (sizeToWrite / totalWrittenBytes) * 100

	yColor.Printf("Managed to write %v/%v (%v%%) bytes to %v\n", totalWrittenBytes, sizeToWrite, writtenPercentage, location)

	if !skipFileIntegrityCheck {
		if verificationPassed {
			gColor.Println("File Integrity -> PASSED...")
		} else {
			rColor.Println("File Integrity -> FAILED...")
		}
	} else {
		yColor.Printf("File Integrity -> SKIPPED")
	}
}

func tidy(dontDeleteFiles bool, fileCount int, location string) {
	yColor := color.New(color.FgYellow)
	if !dontDeleteFiles {
		yColor.Printf("Deleting %v data files...\n", fileCount)
		for i := 1; i < fileCount; i++ {
			os.Remove((path.Join(location, "cloverdata"+strconv.Itoa(i)+".bin")))
		}
	} else {
		yColor.Printf("Skipping file delete...\n")
	}
}

func compareHashBytes(hashBytes1 []byte, hashBytes2 []byte) bool {
	if len(hashBytes1) != len(hashBytes2) {
		return false
	}
	for i := range hashBytes1 {
		if hashBytes1[i] != hashBytes2[i] {
			return false
		}
	}
	return true
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
