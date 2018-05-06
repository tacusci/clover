package cltools

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/tacusci/clover/utils"
)

type endian int

const (
	bigEndian    endian = 0
	littleEndian endian = 1
)

//EXIF tag values
const (
	subfileTypeTag                   uint16 = 0x00fe
	oldSubfileTypeTag                uint16 = 0x00ff
	imageWidthTag                    uint16 = 0x0100
	imageHeightTag                   uint16 = 0x0101
	bitsPerSampleTag                 uint16 = 0x0102
	compressionTag                   uint16 = 0x0103
	photoMetricInterpretationTag     uint16 = 0x0106
	thresholdingTag                  uint16 = 0x0107
	cellWidthTag                     uint16 = 0x0108
	cellLengthTag                    uint16 = 0x0109
	fillOrderTag                     uint16 = 0x010a
	documentNameTag                  uint16 = 0x010d
	imageDescriptionTag              uint16 = 0x010e
	makeTag                          uint16 = 0x010f
	modelTag                         uint16 = 0x0110
	stripOffsetsTag                  uint16 = 0x0111
	orientationTag                   uint16 = 0x0112
	samplesPerPixelTag               uint16 = 0x0115
	rowsPerStripTag                  uint16 = 0x0116
	stripByteCountsTag               uint16 = 0x0117
	minSampleValueTag                uint16 = 0x0118
	maxSampleValueTag                uint16 = 0x0119
	xResolutionTag                   uint16 = 0x011a
	yResolutionTag                   uint16 = 0x011b
	planarConfigurationTag           uint16 = 0x011c
	pageNameTag                      uint16 = 0x011d
	xPositionTag                     uint16 = 0x011e
	yPositionTag                     uint16 = 0x011f
	freeOffsetsTag                   uint16 = 0x0120
	freeByteCountsTag                uint16 = 0x0121
	grayResponseUnitTag              uint16 = 0x0122
	grayResponseCurveTag             uint16 = 0x0123
	t4OptionsTag                     uint16 = 0x0124
	t6OptionsTag                     uint16 = 0x0125
	resolutionUnitTag                uint16 = 0x0128
	pageNumberTag                    uint16 = 0x0129
	colorResponseUnitTag             uint16 = 0x012c
	transferFunctionTag              uint16 = 0x012d
	softwareTag                      uint16 = 0x0131
	modifyDateTag                    uint16 = 0x0132
	artistTag                        uint16 = 0x013b
	hostComputerTag                  uint16 = 0x013c
	predictorTag                     uint16 = 0x013d
	whitePointTag                    uint16 = 0x013e
	primaryChromaticitiesTag         uint16 = 0x013f
	colorMapTag                      uint16 = 0x0140
	halftoneHintsTag                 uint16 = 0x0141
	tileWidthTag                     uint16 = 0x0142
	tileLengthTag                    uint16 = 0x0143
	tileOffsetsTag                   uint16 = 0x0144
	tileByteCountsTag                uint16 = 0x0145
	badFaxLinesTag                   uint16 = 0x0146
	cleanFaxDataTag                  uint16 = 0x0147
	consecutiveBadFaxLinesTag        uint16 = 0x0148
	subIFDA100DataOffsetTag          uint16 = 0x014a
	inkSetTag                        uint16 = 0x014c
	inkNamesTag                      uint16 = 0x014d
	numberOfInksTag                  uint16 = 0x014e
	dotRangeTag                      uint16 = 0x0150
	targetPrinterTag                 uint16 = 0x0151
	extraSamplesTag                  uint16 = 0x0152
	sampleFormatTag                  uint16 = 0x0153
	sMinSampleValueTag               uint16 = 0x0154
	sMaxSampleValueTag               uint16 = 0x0155
	transferRangeTag                 uint16 = 0x0156
	clipPathTag                      uint16 = 0x0157
	xClipPathUnitsTag                uint16 = 0x0158
	yClipPathUnitsTag                uint16 = 0x0159
	indexedTag                       uint16 = 0x015a
	jpegTablesTag                    uint16 = 0x015b
	opiproxyTag                      uint16 = 0x015f
	globalParametersIFDTag           uint16 = 0x0190
	profileTypeTag                   uint16 = 0x0191
	faxProfileTag                    uint16 = 0x0192
	codingMethodsTag                 uint16 = 0x0193
	versionYearTag                   uint16 = 0x0194
	modeNumberTag                    uint16 = 0x0195
	decodeTag                        uint16 = 0x01b1
	defaultImageColorTag             uint16 = 0x01b2
	t82OptionsTag                    uint16 = 0x01b3
	jpegTables2Tag                   uint16 = 0x01b5
	jpegProcTag                      uint16 = 0x0200
	thumbnailOffsetTag               uint16 = 0x0201
	previewImageStartTag             uint16 = 0x0201
	jpegFromRawStartTag              uint16 = 0x0201
	otherImageStartTag               uint16 = 0x0201
	thumbnailLengthTag               uint16 = 0x0202
	previewImageLengthTag            uint16 = 0x0202
	jpegFromRawLengthTag             uint16 = 0x0202
	otherImageLengthTag              uint16 = 0x0202
	jpegRestartIntervalTag           uint16 = 0x0203
	jpegLosslessPredictorsTag        uint16 = 0x0205
	jpegPointTransformsTag           uint16 = 0x0206
	jpegQTablesTag                   uint16 = 0x0207
	jpegDCTablesTag                  uint16 = 0x0208
	jpegACTablesTag                  uint16 = 0x0209
	yCbCrCoefficientsTag             uint16 = 0x0211
	yCbCrSubSamplingTag              uint16 = 0x0212
	yCbCrPositioningTag              uint16 = 0x0213
	referenceBlackWhiteTag           uint16 = 0x0214
	stripRowCountsTag                uint16 = 0x022f
	applicationNotesTag              uint16 = 0x02bc
	usptoMiscellaneousTag            uint16 = 0x03e7
	relatedImageFileFormatTag        uint16 = 0x1000
	relatedImageWidthTag             uint16 = 0x1001
	relatedImageHeightTag            uint16 = 0x1002
	ratingTag                        uint16 = 0x4746
	xpDipXMLTag                      uint16 = 0x4747
	stichInfoTag                     uint16 = 0x4748
	ratingPercentTag                 uint16 = 0x4749
	sonyRawFileTypeTag               uint16 = 0x7000
	sonyToneCurveTag                 uint16 = 0x7010
	vignettingCorrectionTag          uint16 = 0x7031
	vignettingCorrParamsTag          uint16 = 0x7032
	chromaticAberrationCorrectionTag uint16 = 0x7034
	chromaticAberrationCorrParamsTag uint16 = 0x7035
	distortionCorrectionTag          uint16 = 0x7036
	distorionCorrParamsTag           uint16 = 0x7037
	imageIDTag                       uint16 = 0x800d
	wangTag1Tag                      uint16 = 0x80a3
	wangAnnotationTag                uint16 = 0x80a4
	wangTag3Tag                      uint16 = 0x80a5
	wangTag4Tag                      uint16 = 0x80a6
	imageReferencePointsTag          uint16 = 0x80b9
	regionXformTrackPointTag         uint16 = 0x80ba
	warpQuadrilateralTag             uint16 = 0x80bb
	affineTransformMatTag            uint16 = 0x80bc
	matteingTag                      uint16 = 0x80e3
	dataTypeTag                      uint16 = 0x80e4
	imageDepthTag                    uint16 = 0x80e5
	tileDepthTag                     uint16 = 0x80e6
	imageFullWidthTag                uint16 = 0x8214
	imageFullHeightTag               uint16 = 0x8215
	textureFormatTag                 uint16 = 0x8216
	wrapModesTag                     uint16 = 0x8217
	fovCotTag                        uint16 = 0x8218
	matrixWorldToScreen              uint16 = 0x8219
	matrixWorldToCamera              uint16 = 0x821a
	model2Tag                        uint16 = 0x827d
	cfaRepeatPatternDimTag           uint16 = 0x828d
	cfaPattern2Tag                   uint16 = 0x828e
	batteryLevelTag                  uint16 = 0x828f
	kodakIFDTag                      uint16 = 0x8290
	copyrightTag                     uint16 = 0x8298
	exposureTimeTag                  uint16 = 0x829a
	fNumberTag                       uint16 = 0x829d
	mdFileTag                        uint16 = 0x82a5
	mdScalePixelTag                  uint16 = 0x82a6
	mdColorTableTag                  uint16 = 0x82a7
	mdLabNameTag                     uint16 = 0x82a8
	mdSampleInfoTag                  uint16 = 0x82a9
	mdPrepDateTag                    uint16 = 0x82aa
	mdPrepTimeTag                    uint16 = 0x82ab
	mdFileUnitsTag                   uint16 = 0x82ac
	pixelScaleTag                    uint16 = 0x830e
	adventScaleTag                   uint16 = 0x8335
	adventRevisionTag                uint16 = 0x8336
	uic1TagTag                       uint16 = 0x835c
	uic2TagTag                       uint16 = 0x835d
	uic3TagTag                       uint16 = 0x835e
	uic4TagTag                       uint16 = 0x835f
	iptcNAATag                       uint16 = 0x83bb
	intergraphPacketDataTag          uint16 = 0x847e
	intergraphFlagRegistersTag       uint16 = 0x847f
	intergraphMatrixTag              uint16 = 0x8480
	ingrReservedTag                  uint16 = 0x8481
	modelTiePointTag                 uint16 = 0x8482
	siteTag                          uint16 = 0x84e0
	colorSequenceTag                 uint16 = 0x84e1
	it8HeaderTag                     uint16 = 0x84e2
	rasterPaddingTag                 uint16 = 0x84e3
	bitsPerRunLengthTag              uint16 = 0x84e4
	bitsPerExtendedRunLengthTag      uint16 = 0x84e5
	colorTableTag                    uint16 = 0x84e6
	imageColorIndicatorTag           uint16 = 0x84e7
	backgroundColorIndictorTag       uint16 = 0x84e8
	imageColorValueTag               uint16 = 0x84e9
	backgroundColorValueTag          uint16 = 0x84ea
	pixelIntensityRangeTag           uint16 = 0x84eb
	transparencyIndicatorTag         uint16 = 0x84ec
	colorCharacterizationTag         uint16 = 0x84ed
	hcUsageTag                       uint16 = 0x84ee
	trapIndicatorTag                 uint16 = 0x84ef
	cmykEquivalentTag                uint16 = 0x84f0
	semInfoTag                       uint16 = 0x8546
	afcpIPTCTag                      uint16 = 0x8568
	pixelMagicJBigOptionsTag         uint16 = 0x85b8
	jplCartoIFDTag                   uint16 = 0x85d7
	modelTransformTag                uint16 = 0x85d8
)

type tiffHeaderData struct {
	endianOrder endian
	magicNum    uint16
	tiffOffset  uint32
}

//RunRtc runs the raw to compressed image conversion tool
func RunRtc(locationpath string, intputType string, outputType string) {
	if len(locationpath) == 0 || len(intputType) == 0 || len(outputType) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if isDir, _ := isDirectory(locationpath); isDir {
		fileInfos, err := ioutil.ReadDir(locationpath)

		if err != nil {
			log.Fatal(err)
			return
		}

		for i := range fileInfos {
			fileInfo := fileInfos[i]
			if !fileInfo.IsDir() && strings.Contains(strings.ToLower(fileInfo.Name()), strings.ToLower(intputType)) {
				filename := path.Join(locationpath, fileInfo.Name())
				filename = utils.TranslatePath(filename)
				imageFile, err := os.Open(filename)
				if err != nil {
					log.Fatal(err)
					return
				}
				defer imageFile.Close()
				err = parseAllImageMeta(imageFile)

				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}
	}
}

func parseAllImageMeta(file *os.File) error {
	header, err := readHeader(file)
	imageTiffHeaderData := *new(tiffHeaderData)

	if err != nil {
		//return here before the next read + check because we always want the root cause error to bubble back up
		return err
	}

	imageTiffHeaderData, err = getTiffHeader(header)

	ifd0Data := readIfd(file, imageTiffHeaderData.tiffOffset, imageTiffHeaderData.endianOrder)

	for i := range ifd0Data {
		fmt.Printf("%#x ", ifd0Data[i])
	}

	return nil
}

func readIfd(file *os.File, ifdOffset uint32, endianOrder endian) []byte {
	ifdTagCountBytes := make([]byte, 2)
	file.Seek(int64(ifdOffset), os.SEEK_SET)
	file.Read(ifdTagCountBytes)

	var ifdTagCount uint16
	if endianOrder == bigEndian {
		ifdTagCount |= uint16(ifdTagCountBytes[0]) << 8
		ifdTagCount |= uint16(ifdTagCountBytes[1])
	} else if endianOrder == littleEndian {
		ifdTagCount |= uint16(ifdTagCountBytes[0])
		ifdTagCount |= uint16(ifdTagCountBytes[1]) << 8
	}

	//each IFD tag length is 12 bytes
	ifdData := make([]byte, ifdTagCount*12)
	file.Seek(int64(ifdOffset+2), os.SEEK_SET)
	file.Read(ifdData)

	return ifdData
}

func readHeader(file *os.File) ([]byte, error) {
	header := make([]byte, 8)
	file.Seek(0, 0)
	bytesRead, err := file.Read(header)

	if bytesRead < 8 {
		return header, errors.New("Unable to read full header")
	}

	if err != nil {
		return header, err
	}
	return header, nil
}

func getEdianOrder(header []byte) endian {
	if len(header) >= 4 {
		var endianFlag uint16
		//add the bits to the 2 byte int and shove them to the left to make room for the other bits
		endianFlag |= uint16(header[0]) << 8
		endianFlag |= uint16(header[1])
		if endianFlag == 19789 {
			return bigEndian
		} else if endianFlag == 18761 {
			return littleEndian
		}
	}
	return bigEndian
}

func getTiffHeader(header []byte) (tiffHeaderData, error) {
	tiffData := new(tiffHeaderData)
	tiffData.endianOrder = getEdianOrder(header)

	if len(header) >= 8 {

		var magicNum uint16
		if tiffData.endianOrder == bigEndian {
			magicNum |= uint16(header[2]) | uint16(header[3])
		} else if tiffData.endianOrder == littleEndian {
			magicNum |= uint16(header[3]) | uint16(header[2])
		}

		tiffData.magicNum = magicNum

		var tiffOffset uint32
		if tiffData.endianOrder == bigEndian {
			tiffOffset |= uint32(header[4]) << 24
			tiffOffset |= uint32(header[5]) << 16
			tiffOffset |= uint32(header[6]) << 8
			tiffOffset |= uint32(header[7])
		} else if tiffData.endianOrder == littleEndian {
			tiffOffset |= uint32(header[4])
			tiffOffset |= uint32(header[5]) << 8
			tiffOffset |= uint32(header[6]) << 16
			tiffOffset |= uint32(header[7]) << 24
		}

		tiffData.tiffOffset = tiffOffset
	} else {
		return *tiffData, errors.New("Header incorrect length")
	}
	return *tiffData, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
