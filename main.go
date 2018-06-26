package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tacusci/clover/cltools"
	"github.com/tacusci/logging"
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

func setLoggingLevel() {
	debugLevel := flag.Bool("debug", false, "Set logging to debug")
	flag.Parse()

	loggingLevel := logging.InfoLevel

	if *debugLevel {
		logging.SetLevel(logging.DebugLevel)
		return
	}
	logging.SetLevel(loggingLevel)
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
		locationPath := flag.String("l", "", "Location to write data to.")
		sizeToWrite := flag.Int("s", 0, "Size of total data to write.")
		skipFileIntegrityCheck := flag.Bool("sic", false, "Skip verifying output file integrity.")
		dontDeleteFiles := flag.Bool("nd", false, "Don't delete outputted files.")
		setLoggingLevel()

		flag.Parse()

		cltools.RunSdc(*locationPath, *sizeToWrite, *skipFileIntegrityCheck, *dontDeleteFiles)

	} else if toolFlag == "/rtc" {
		sourceDirectory := flag.String("id", "", "Location containing raw images to convert.")
		outputDirectory := flag.String("od", "", "Location to save compressed images.")
		inputType := flag.String("it", "", "Extension of image type to convert.")
		outputType := flag.String("ot", "", "Extension of image type to output to.")
		recursive := flag.Bool("rs", false, "Scan all sub folders in root recursively.")
		setLoggingLevel()
		logging.OutputDateTime, logging.OutputPath, logging.OutputLogLevelFlag, logging.OutputArrowSuffix = false, false, false, false
		flag.Parse()

		cltools.RunRtc(*sourceDirectory, *outputDirectory, *inputType, *outputType, *recursive)
	} else {
		outputUsageAndClose()
	}
}
