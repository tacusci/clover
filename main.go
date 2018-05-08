package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tacusci/clover/cltools"
)

func outputUsage() {
	println("Usage: " + os.Args[0] + " </TOOLFLAG>")
	fmt.Printf("\t/sdc (StorageDeviceChecker) - Tool for checking size of storage devices.\n")
	fmt.Printf("\t/rtc (RawToCompressed) - Tool for batch compressing raw images.\n")
}

func outputUsageAndClose() {
	outputUsage()
	os.Exit(1)
}

func main() {

	if len(os.Args) == 1 {
		outputUsageAndClose()
	}

	runTool(os.Args[1])
}

func runTool(toolFlag string) {
	//kind of hack to force flag parser to find tool argument flags correctly
	os.Args = os.Args[1:]
	if toolFlag == "/sdc" {
		locationPath := flag.String("location", "", "Location to write data to.")
		sizeToWrite := flag.Int("size", 0, "Size of total data to write.")
		skipFileIntegrityCheck := flag.Bool("skip-integrity-check", false, "Skip verifying output file integrity.")
		dontDeleteFiles := flag.Bool("no-delete", false, "Don't delete outputted files.")

		flag.Parse()

		cltools.RunSdc(*locationPath, *sizeToWrite, *skipFileIntegrityCheck, *dontDeleteFiles)

	} else if toolFlag == "/rtc" {
		sourceDirectory := flag.String("input-directory", "", "Location containing raw images to convert.")
		outputDirectory := flag.String("output-directory", "", "Location to save compressed images.")
		inputType := flag.String("input-type", "", "Extension of image type to convert.")
		outputType := flag.String("output-type", "", "Extension of image type to output to.")
		recursive := flag.Bool("recursive-search", false, "Scan all sub folders in root recursively.")

		flag.Parse()

		cltools.RunRtc(*sourceDirectory, *outputDirectory, *inputType, *outputType, *recursive)
	} else {
		outputUsageAndClose()
	}
}
