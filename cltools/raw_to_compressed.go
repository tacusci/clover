package cltools

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"strings"

	"github.com/tacusci/clover/utils"
)

//EXIF tag values
const (
	subfileTypeTag                   uint16 = 0x00fe
	oldSubfileTypeTag                uint16 = 0x00ff
	imageWidthTag                    uint16 = 0x0100
	imageHeightTag                   uint16 = 0x0101
	bitsPerSampleTag                 uint16 = 0x0102
	compressionTag                   uint16 = 0x0103
	photometricInterpretationTag     uint16 = 0x0106
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
	wbGRGBLevelsTag                  uint16 = 0x8602
	leafDataTag                      uint16 = 0x8606
	photoshopSettingsTag             uint16 = 0x8649
	exifOffsetTag                    uint16 = 0x8769
	iccProfileTag                    uint16 = 0x8773
	tiffFxExtensionsTag              uint16 = 0x8774
	multiProfilesTag                 uint16 = 0x8780
	sharedDataTag                    uint16 = 0x8781
	t88OptionsTag                    uint16 = 0x8782
	imageLayerTag                    uint16 = 0x87ac
	geoTiffDirectoryTag              uint16 = 0x87af
	geoTiffDoubleParamsTag           uint16 = 0x87b0
	geoTiffASCIIParamsTag            uint16 = 0x87b1
	jBIGOptionsTag                   uint16 = 0x87be
	exposureProgramTag               uint16 = 0x8822
	spectralSensitivityTag           uint16 = 0x8824
	gpsInfoTag                       uint16 = 0x8825
	isoTag                           uint16 = 0x8827
	optoElectricConvFactorTag        uint16 = 0x8828
	interlaceTag                     uint16 = 0x8829
	timeZoneOffsetTag                uint16 = 0x882a
	selfTimerModeTag                 uint16 = 0x882b
	sensitivityTypeTag               uint16 = 0x8830
	standardOutputSensitivityTag     uint16 = 0x8831
	recommendedExposureIndexTag      uint16 = 0x8832
	isoSpeedTag                      uint16 = 0x8833
	isoSpeedLatitudeyyyTag           uint16 = 0x8834
	isoSpeedLatitudezzzTag           uint16 = 0x8835
	faxRecvParamsTag                 uint16 = 0x885c
	faxSubAddressTag                 uint16 = 0x885d
	faxRecvTimeTag                   uint16 = 0x885e
	fedexEDRTag                      uint16 = 0x8871
	leafSubIDTag                     uint16 = 0x888a
	exifVersionTag                   uint16 = 0x9000
	dateTimeOriginalTag              uint16 = 0x9003
	createDateTag                    uint16 = 0x9004
	googlePlusUploadCodeTag          uint16 = 0x9009
	offsetTimeTag                    uint16 = 0x9010
	offsetTimeOriginalTag            uint16 = 0x9011
	offsetTimeDigitizedTag           uint16 = 0x9012
	componentsConfigurationTag       uint16 = 0x9101
	compressedBitsPerPixel           uint16 = 0x9102
	shutterSpeedValueTag             uint16 = 0x9201
	apertueValueTag                  uint16 = 0x9202
	brightnessValueTag               uint16 = 0x9203
	exposureCompensationTag          uint16 = 0x9204
	maxApertureValueTag              uint16 = 0x9205
	subjectDistanceTag               uint16 = 0x9206
	meteringModeTag                  uint16 = 0x9207
	lightSourceTag                   uint16 = 0x9208
	flashTag                         uint16 = 0x9209
	focalLengthTag                   uint16 = 0x920a
	flashEnergyTag                   uint16 = 0x920b
	spatialFrequencyResponseTag      uint16 = 0x920c
	noiseTag                         uint16 = 0x920d
	focalPlaneXResolutionTag         uint16 = 0x920e
	focalPlaneYResolutionTag         uint16 = 0x920f
	focalPlaneResolutionUnitTag      uint16 = 0x9210
	imageNumberTag                   uint16 = 0x9211
	securityClassificationTag        uint16 = 0x9212
	imageHistoryTag                  uint16 = 0x9213
	subjectAreaTag                   uint16 = 0x9214
	exposureIndexTag                 uint16 = 0x9215
	tiffEPStandardIDTag              uint16 = 0x9216
	sensingMethodTag                 uint16 = 0x9217
	cip3DataFileTag                  uint16 = 0x923a
	cip3SheetTag                     uint16 = 0x923b
	cip3SideTag                      uint16 = 0x923c
	stoNitsTag                       uint16 = 0x923f
	makerNoteAppleTag                uint16 = 0x927c
	makerNoteNikonTag                uint16 = 0x927c
	makerNoteCanonTag                uint16 = 0x927c
	makerNoteCasioTag                uint16 = 0x927c
	makerNoteCasio2Tag               uint16 = 0x927c
	makerNoteDJITag                  uint16 = 0x927c
	makerNoteFLIRTag                 uint16 = 0x927c
	makerNoteFujiFilmTag             uint16 = 0x927c
	makerNoteGETag                   uint16 = 0x927c
	makerNoteGE2Tag                  uint16 = 0x927c
	makerNoteHasselbladTag           uint16 = 0x927c
	makerNoteHPTag                   uint16 = 0x927c
	makerNoteHP2Tag                  uint16 = 0x927c
	makerNoteHP4Tag                  uint16 = 0x927c
	makerNoteHP6Tag                  uint16 = 0x927c
	makerNoteISLTag                  uint16 = 0x927c
	makerNoteJVCTag                  uint16 = 0x927c
	makerNoteJVCTextTag              uint16 = 0x927c
	makerNoteKodak1aTag              uint16 = 0x927c
	makerNoteKodak1bTag              uint16 = 0x927c
	makerNoteKodak2Tag               uint16 = 0x927c
	makerNoteKodak3Tag               uint16 = 0x927c
	makerNoteKodak4Tag               uint16 = 0x927c
	makerNoteKodak5Tag               uint16 = 0x927c
	makerNoteKodak6aTag              uint16 = 0x927c
	makerNoteKodak6bTag              uint16 = 0x927c
	makerNoteKodak7Tag               uint16 = 0x927c
	makerNoteKodak8aTag              uint16 = 0x927c
	makerNoteKodak8bTag              uint16 = 0x927c
	makerNoteKodak8cTag              uint16 = 0x927c
	makerNoteKodak9Tag               uint16 = 0x927c
	makerNoteKodak10Tag              uint16 = 0x927c
	makerNoteKodak11Tag              uint16 = 0x927c
	makerNoteKodakUnknownTag         uint16 = 0x927c
	makerNoteKyoceraTag              uint16 = 0x927c
	makerNoteMinoltaTag              uint16 = 0x927c
	makerNoteMinolta2Tag             uint16 = 0x927c
	makerNoteMinolta3Tag             uint16 = 0x927c
	makerNoteMotorolaTag             uint16 = 0x927c
	makerNoteNikon2Tag               uint16 = 0x927c
	makerNoteNikon3Tag               uint16 = 0x927c
	makerNoteNintendoTag             uint16 = 0x927c
	makerNoteOlympusTag              uint16 = 0x927c
	makerNoteOlympus2Tag             uint16 = 0x927c
	makerNoteLeicaTag                uint16 = 0x927c
	makerNoteLeica2Tag               uint16 = 0x927c
	makerNoteLeica3Tag               uint16 = 0x927c
	makerNoteLeica4Tag               uint16 = 0x927c
	makerNoteLeica5Tag               uint16 = 0x927c
	makerNoteLeica6Tag               uint16 = 0x927c
	makerNoteLeica7Tag               uint16 = 0x927c
	makerNoteLeica8Tag               uint16 = 0x927c
	makerNoteLeica9Tag               uint16 = 0x927c
	makerNotePanasonicTag            uint16 = 0x927c
	makerNotePentaxTag               uint16 = 0x927c
	makerNotePentax2Tag              uint16 = 0x927c
	makerNotePentax3Tag              uint16 = 0x927c
	makerNotePentax4Tag              uint16 = 0x927c
	makerNotePentax5Tag              uint16 = 0x927c
	makerNotePentax6Tag              uint16 = 0x927c
	makerNotePhaseOneTag             uint16 = 0x927c
	makerNoteReconyxTag              uint16 = 0x927c
	makerNoteReconyx2Tag             uint16 = 0x927c
	makerNoteRicohTag                uint16 = 0x927c
	makerNoteRicoh2Tag               uint16 = 0x927c
	makerNoteRicohTextTag            uint16 = 0x927c
	makerNoteSamsung1aTag            uint16 = 0x927c
	makerNoteSamsung1bTag            uint16 = 0x927c
	makerNoteSamsung2Tag             uint16 = 0x927c
	makerNoteSanyoTag                uint16 = 0x927c
	makerNoteSanyoC4Tag              uint16 = 0x927c
	makerNoteSanyoPatchTag           uint16 = 0x927c
	makerNoteSigmaTag                uint16 = 0x927c
	makerNoteSonyTag                 uint16 = 0x927c
	makerNoteSony2Tag                uint16 = 0x927c
	makerNoteSony3Tag                uint16 = 0x927c
	makerNoteSony4Tag                uint16 = 0x927c
	makerNoteSony5Tag                uint16 = 0x927c
	makerNoteSonyEricssonTag         uint16 = 0x927c
	makerNoteSonySRFTag              uint16 = 0x927c
	makerNoteUnknownTextTag          uint16 = 0x927c
	makerNoteUnknownBinaryTag        uint16 = 0x927c
	makerNoteUknownTag               uint16 = 0x927c
	userCommentTag                   uint16 = 0x9286
	subSecTimeTag                    uint16 = 0x9290
	subSecTimeOriginalTag            uint16 = 0x9291
	subSecTimeDigitizedTag           uint16 = 0x9292
	msDocumentTextTag                uint16 = 0x932f
	msPropertySetStorageTag          uint16 = 0x9330
	msDocumentTextPositionTag        uint16 = 0x9331
	imageSourceDataTag               uint16 = 0x935c
	ambientTempratureTag             uint16 = 0x9400
	humidityTag                      uint16 = 0x9401
	pressureTag                      uint16 = 0x9402
	waterDepthTag                    uint16 = 0x9403
	accelerationTag                  uint16 = 0x9404
	cameraElevationAngleTag          uint16 = 0x9405
	xpTitleTag                       uint16 = 0x9c9b
	xpCommentTag                     uint16 = 0x9c9c
	xpAuthorTag                      uint16 = 0x9c9d
	xpKeywordsTag                    uint16 = 0x9c9e
	xpSubjectTag                     uint16 = 0x9c9f
	flashpixVersionTag               uint16 = 0xa000
	colorSpaceTag                    uint16 = 0xa001
	exifImageWidthTag                uint16 = 0xa002
	exifImageHeightTag               uint16 = 0xa003
	relatedSoundFileTag              uint16 = 0xa004
	interopOffsetTag                 uint16 = 0xa005
	samsungRawPointersOffsetTag      uint16 = 0xa010
	samsungRawPointersLengthTag      uint16 = 0xa011
	samsungRawByteOrderTag           uint16 = 0xa101
	samsungRawUnknownTag             uint16 = 0xa102
	flashEnergy2Tag                  uint16 = 0xa20b
	spatialFrequencyResponse2Tag     uint16 = 0xa20c
	noise2Tag                        uint16 = 0xa20d
	focalPlaneXResolution2Tag        uint16 = 0xa20e
	focalPlaneYResolution2Tag        uint16 = 0xa20f
	focalPlaneResolutionUnit2Tag     uint16 = 0xa210
	imageNumber2Tag                  uint16 = 0xa211
	securityClassification2Tag       uint16 = 0xa212
	imageHistory2Tag                 uint16 = 0xa213
	subjectLocationTag               uint16 = 0xa214
	exposureIndex2Tag                uint16 = 0xa215
	tiffEPStandardID2Tag             uint16 = 0xa216
	sensingMethod2Tag                uint16 = 0xa217
	fileSourceTag                    uint16 = 0xa300
	sceneTypeTag                     uint16 = 0xa301
	cfaPatternTag                    uint16 = 0xa302
	customRenderedTag                uint16 = 0xa401
	exposureModeTag                  uint16 = 0xa402
	whiteBalanceTag                  uint16 = 0xa403
	digitalZoomRatioTag              uint16 = 0xa404
	focalLengthIn35mmFormatTag       uint16 = 0xa405
	sceneCaptureTypeTag              uint16 = 0xa406
	gainControlTag                   uint16 = 0xa407
	contrastTag                      uint16 = 0xa408
	saturationTag                    uint16 = 0xa409
	sharpnessTag                     uint16 = 0xa40a
	deviceSettingDescriptionTag      uint16 = 0xa40b
	subjectDistanceRangeTag          uint16 = 0xa40c
	imageUniqueIDTag                 uint16 = 0xa420
	ownerNameTag                     uint16 = 0xa430
	serialNumberTag                  uint16 = 0xa431
	lensInfoTag                      uint16 = 0xa432
	lensMakeTag                      uint16 = 0xa433
	lensModelTag                     uint16 = 0xa434
	lensSerialNumberTag              uint16 = 0xa435
	gdalMetadataTag                  uint16 = 0xa480
	gdalNoDataTag                    uint16 = 0xa481
	gammaTag                         uint16 = 0xa500
	expandSoftwareTag                uint16 = 0xafc0
	expandLensTag                    uint16 = 0xafc1
	expandFilmTag                    uint16 = 0xafc2
	expandFilterLensTag              uint16 = 0xafc3
	expandScannerTag                 uint16 = 0xafc4
	expandFlashLampTag               uint16 = 0xafc5
	pixelFormatTag                   uint16 = 0xbc01
	transformationTag                uint16 = 0xbc02
	uncompressedTag                  uint16 = 0xbc03
	imageTypeTag                     uint16 = 0xbc04
	imageWidth2Tag                   uint16 = 0xbc80
	imageHeight2Tag                  uint16 = 0xbc81
	widthResolutionTag               uint16 = 0xbc82
	heightResolutionTag              uint16 = 0xbc83
	imageOffsetTag                   uint16 = 0xbcc0
	imageByteCountTag                uint16 = 0xbcc1
	alphaOffsetTag                   uint16 = 0xbcc2
	alphaByteCountTag                uint16 = 0xbcc3
	imageDataDiscardTag              uint16 = 0xbcc4
	alphaDataDiscardTag              uint16 = 0xbcc5
	oceScanjobDescTag                uint16 = 0xc427
	oceApplicationSelectorTag        uint16 = 0xc428
	oceIDNumberTag                   uint16 = 0xc429
	oceImageLogicTag                 uint16 = 0xc42a
	annotationsTag                   uint16 = 0xc44f
	printImTag                       uint16 = 0xc4a5
	originalFileNameTag              uint16 = 0xc573
	usptoOriginalContentTypeTag      uint16 = 0xc580
	cr2cfaPatternTag                 uint16 = 0xc5e0
	dngVersionTag                    uint16 = 0xc612
	dngBackwardVersionTag            uint16 = 0xc613
	uniqueCameraModelTag             uint16 = 0xc614
	cfaPlaneColorTag                 uint16 = 0xc616
	cfaLayoutTag                     uint16 = 0xc617
	linearizationTableTag            uint16 = 0xc618
	blackLevelRepeatDimTag           uint16 = 0xc619
	blackLevelTag                    uint16 = 0xc61a
	blackLevelDeltaHTag              uint16 = 0xc61b
	blackLevelDeltaVTag              uint16 = 0xc61c
	whiteLevelTag                    uint16 = 0xc61d
	defaultScaleTag                  uint16 = 0xc61e
	defaultCropOriginTag             uint16 = 0xc61f
	defaultCropSizeTag               uint16 = 0xc620
	colorMatrix1Tag                  uint16 = 0xc621
	colorMatrix2Tag                  uint16 = 0xc622
	cameraCalibrationTag             uint16 = 0xc623
	cameraCalibration2Tag            uint16 = 0xc624
	reductionMatrix1Tag              uint16 = 0xc625
	reductionMatrix2Tag              uint16 = 0xc626
	analogBalanceTag                 uint16 = 0xc627
	asShotNeutralTag                 uint16 = 0xc628
	asShotWhiteXYTag                 uint16 = 0xc629
	baselineExposureTag              uint16 = 0xc62a
	baselineNoiseTag                 uint16 = 0xc62b
	baselineSharpnessTag             uint16 = 0xc62c
	bayerGreenSplitTag               uint16 = 0xc62d
	linearResponseLimitTag           uint16 = 0xc62e
	cameraSerialNumberTag            uint16 = 0xc62f
	dngLensInfoTag                   uint16 = 0xc630
	chromaBlurRadiusTag              uint16 = 0xc631
	antiAliasStrengthTag             uint16 = 0xc632
	shadowScaleTag                   uint16 = 0xc633
	sr2PrivateTag                    uint16 = 0xc634
	dngAdobeDataTag                  uint16 = 0xc634
	makerNotePentax22Tag             uint16 = 0xc634
	makerNotePentax52Tag             uint16 = 0xc634
	dngPrivateDataTag                uint16 = 0xc634
	makerNoteSafetyTag               uint16 = 0xc635
	rawImageSegmentationTag          uint16 = 0xc640
	calibrationIlluminant1Tag        uint16 = 0xc65a
	calibrationIlluminant2Tag        uint16 = 0xc65b
	bestQualityScaleTag              uint16 = 0xc65c
	rawDataUniqueIDTag               uint16 = 0xc65d
	aliasLayerMetadataTag            uint16 = 0xc660
	originalRawFileNameTag           uint16 = 0xc68b
	originalRawFileDataTag           uint16 = 0xc68c
	activeAreaTag                    uint16 = 0xc68d
	maskedAreasTag                   uint16 = 0xc68e
	asShotICCProfileTag              uint16 = 0xc68f
	asShotPreProfileMatrixTag        uint16 = 0xc690
	colorimetricReferenceTag         uint16 = 0xc6bf
	sRawTypeTag                      uint16 = 0xc6c5
	panasonicTitleTag                uint16 = 0xc6d2
	panasonicTitle2Tag               uint16 = 0xc6d3
	cameraCalibrationSigTag          uint16 = 0xc6f3
	profileCalibrationSigTag         uint16 = 0xc6f4
	profileIFDTag                    uint16 = 0xc6f5
	asShotProfileNameTag             uint16 = 0xc6f6
	noiseReductionAppliedTag         uint16 = 0xc6f7
	profileNameTag                   uint16 = 0xc6f8
	profileHueSatMapDimsTag          uint16 = 0xc6f9
	profileHueSatMapData1Tag         uint16 = 0xc6fa
	profileHueSatMapData2Tag         uint16 = 0xc6fb
	profileToneCurveTag              uint16 = 0xc6fc
	profileEmbedPolicyTag            uint16 = 0xc6fd
	profileCopyrightTag              uint16 = 0xc6fe
	forwardMatrix1Tag                uint16 = 0xc714
	forwardMatrix2Tag                uint16 = 0xc715
	previewApplicationNameTag        uint16 = 0xc716
	previewApplicationVersionTag     uint16 = 0xc717
	previewSettingsNameTag           uint16 = 0xc718
	previewSettingsDigestTag         uint16 = 0xc719
	previewColorSpaceTag             uint16 = 0xc71a
	previewDateTimeTag               uint16 = 0xc71b
	rawImageDigestTag                uint16 = 0xc71c
	originalRawFileDigestTag         uint16 = 0xc71d
	subTileBlockSizeTag              uint16 = 0xc71e
	rowInterleaveFactorTag           uint16 = 0xc71f
	profileLookTableDimsTag          uint16 = 0xc725
	profileLookTableDataTag          uint16 = 0xc726
	opcodeList1Tag                   uint16 = 0xc740
	opcodeList2Tag                   uint16 = 0xc741
	opcodeList3Tag                   uint16 = 0xc74e
	noiseProfileTag                  uint16 = 0xc761
	timeCodesTag                     uint16 = 0xc763
	frameRateTag                     uint16 = 0xc764
	tStopTag                         uint16 = 0xc772
	reelNameTag                      uint16 = 0xc789
	originalDefaultFinalSizeTag      uint16 = 0xc791
	originalBestQualitySizeTag       uint16 = 0xc792
	originalDefaultCropSizeTag       uint16 = 0xc793
	cameraLabelTag                   uint16 = 0xc7a1
	profileHueSatMapEncodingTag      uint16 = 0xc7a3
	profileLookTableEncodingTag      uint16 = 0xc7a4
	baselineExposureOffsetTag        uint16 = 0xc7a5
	defaultBlackRenderTag            uint16 = 0xc7a6
	newRawImageDigestTag             uint16 = 0xc7a7
	rawToPreviewGainTag              uint16 = 0xc7a8
	defaultUserCropTag               uint16 = 0xc7b5
	paddingTag                       uint16 = 0xea1c
	offsetSchemaTag                  uint16 = 0xea1d
	ownerName2Tag                    uint16 = 0xfde8
	serialNumber2Tag                 uint16 = 0xfde9
	lensTag                          uint16 = 0xfdea
	kdcIFDTag                        uint16 = 0xfe00
	rawFileTag                       uint16 = 0xfe4c
	converterTag                     uint16 = 0xfe4d
	whiteBalance2Tag                 uint16 = 0xfe4e
	exposureTag                      uint16 = 0xfe51
	shadowsTag                       uint16 = 0xfe52
	brightnessTag                    uint16 = 0xfe53
	contrast2Tag                     uint16 = 0xfe54
	saturation2Tag                   uint16 = 0xfe55
	sharpness2Tag                    uint16 = 0xfe56
	smoothnessTag                    uint16 = 0xfe57
	moireFilterTag                   uint16 = 0xfe58

	unsignedByteType     uint8 = 1  //is 1 byte in size
	asciiStringsType     uint8 = 2  //ASCII strings, always a 1 byte long pointer
	unsignedShortType    uint8 = 3  //is 2 bytes in size
	unsignedLongType     uint8 = 4  //is 4 bytes in size
	unsignedRationalType uint8 = 5  //is 8 bytes in size
	signedByteType       uint8 = 6  //is 1 bytes in size
	undefinedType        uint8 = 7  //is 1 byte in size?
	signedShortType      uint8 = 8  //is 2 bytes in size
	signedLongType       uint8 = 9  //is 4 bytes in size
	signedRationalType   uint8 = 10 //is 8 bytes in size
	singleFloatType      uint8 = 11 //is 4 bytes in size
	doubleFloatType      uint8 = 12 //is 8 bytes in size

	photometricInterpretationMinIsWhite       uint16 = 0
	photometricInterpretationMinIsBlack       uint16 = 1
	photometricInterpretationRGB              uint16 = 2
	photometricInterpretationPaletteColor     uint16 = 3
	photometricInterpretationTransparencyMask uint16 = 4
	photometricInterpretationSeperated        uint16 = 5
	photometricInterpretationYCBCR            uint16 = 6
	photometricInterpretationCILAB            uint16 = 8
	photometricInterpretationICCLAB           uint16 = 9
	photometricInterpretationITULAB           uint16 = 10
	photometricInterpretationLOGL             uint16 = 32844
	photometricInterpretationLOGLUV           uint16 = 32845

	compressionNone                uint16 = 1
	compressionCCITTRLE            uint16 = 2
	compressionCCITTFAX3           uint16 = 3
	compressionCCITTFAX4           uint16 = 4
	compressionLZW                 uint16 = 5
	compressionOJPEG               uint16 = 6
	compressionJPEG                uint16 = 7
	compressionADOBEDEFLATE        uint16 = 8
	compressionJBIGOnBlackAndWhite uint16 = 9
	compressionJBIGOnColor         uint16 = 10
)

type tiffHeaderData struct {
	endianOrder utils.EndianOrder
	magicNum    uint16
	tiffOffset  uint32
}

//RunRtc runs the raw to compressed image conversion tool
func RunRtc(locationpath string, outputDirectory string, inputType string, outputType string, recursive bool) {
	if len(locationpath) == 0 || len(inputType) == 0 || len(outputType) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if isDir, _ := isDirectory(locationpath); isDir {
		var imagesFoundCount int
		imagesFoundCount = readAllImagesInDir(imagesFoundCount, locationpath, outputDirectory, inputType, outputType, recursive)
		log.Printf("Finished running... Found %d images\n", imagesFoundCount)
	}
}

func readAllImagesInDir(imagesFoundCount int, locationpath string, outputDirectory string, inputType string, outputType string, recursive bool) int {

	fileInfos, err := ioutil.ReadDir(locationpath)

	if err != nil {
		log.Fatal(err)
		return imagesFoundCount
	}

	for i := range fileInfos {
		fileInfo := fileInfos[i]

		if !fileInfo.IsDir() && strings.HasSuffix(strings.ToLower(fileInfo.Name()), strings.ToLower(inputType)) {
			filename := path.Join(locationpath, fileInfo.Name())
			filename = utils.TranslatePath(filename)
			imageFile, err := os.Open(filename)

			if err != nil {
				log.Fatal(err)
				return imagesFoundCount
			}

			log.Printf("Reading image %s", filename)
			err = parseAllImageData(imageFile)

			if err != nil {
				log.Fatal(err)
				//return imagesFoundCount
			}

			imagesFoundCount++
			imageFile.Close()
		} else if recursive && fileInfo.IsDir() {
			imagesFoundCount = readAllImagesInDir(imagesFoundCount, utils.TranslatePath(path.Join(locationpath, fileInfo.Name())), outputDirectory, inputType, outputType, recursive)
		}
	}
	return imagesFoundCount
}

func parseAllImageData(file *os.File) error {
	header, err := readHeader(file)
	imageTiffHeaderData := *new(tiffHeaderData)

	if err != nil {
		//return here before the next read + check because we always want the root cause error to bubble back up
		return err
	}

	imageTiffHeaderData, err = getTiffHeader(header)

	ifd0Data := readIfd(file, imageTiffHeaderData.tiffOffset, imageTiffHeaderData.endianOrder)

	//for each byte in the IFD0
	for i := range ifd0Data {
		if math.Mod(float64(i), float64(12)) == 0 {
			tagAsInt := utils.ConvertBytesToUInt16(ifd0Data[i], ifd0Data[i+1], imageTiffHeaderData.endianOrder)
			dataFormatAsInt := utils.ConvertBytesToUInt16(ifd0Data[i+2], ifd0Data[i+3], imageTiffHeaderData.endianOrder)
			numOfElementsAsInt := utils.ConvertBytesToUInt32(ifd0Data[i+4], ifd0Data[i+5], ifd0Data[i+6], ifd0Data[i+7], imageTiffHeaderData.endianOrder)
			dataValueOrDataOffsetAsInt := utils.ConvertBytesToUInt32(ifd0Data[i+8], ifd0Data[i+9], ifd0Data[i+10], ifd0Data[i+11], imageTiffHeaderData.endianOrder)

			if tagAsInt == subfileTypeTag {
				if uint8(dataFormatAsInt) == unsignedLongType {
					if numOfElementsAsInt == 1 {
						firstBitFlag := ifd0Data[i+11]  //if first bit is 1 then its reduced resolution
						secondBitFlag := ifd0Data[i+10] //if second bit is 1 then its a single page image of a multi-page image
						thirdBitFlag := ifd0Data[i+9]   //if the third bit is 1 then image defines transparency mask for another image in tiff file. The Photometric interpritation value must be 4
						fourthBitFlag := ifd0Data[i+8]  //if the forth bit is 1 then MRC imaging model

						if firstBitFlag == 1 {
							fmt.Printf("Image type is -> Reduced resolution image\n")
						} else if secondBitFlag == 1 {
							fmt.Printf("Image type is -> Single page of multipage image\n")
						} else if thirdBitFlag == 1 {
							fmt.Printf("Image type is -> Transparency mask image\n")
						} else if fourthBitFlag == 1 {
							fmt.Printf("MRC imaging model?\n")
						}
					}
				}
			} else if tagAsInt == imageWidthTag {
				if uint8(dataFormatAsInt) == unsignedLongType {
					fmt.Printf("Image width -> %d\n", dataValueOrDataOffsetAsInt)
				}
			} else if tagAsInt == imageHeightTag {
				if uint8(dataFormatAsInt) == unsignedLongType {
					fmt.Printf("Image height -> %d\n", dataValueOrDataOffsetAsInt)
				}
			} else if tagAsInt == bitsPerSampleTag {
				if uint8(dataFormatAsInt) == unsignedShortType {
					file.Seek(int64(dataValueOrDataOffsetAsInt), os.SEEK_SET)
					bitsPerSampleData := make([]byte, numOfElementsAsInt)
					file.Read(bitsPerSampleData)
					fmt.Printf("Bits per sample -> %d\n", bitsPerSampleData)
				}
			} else if tagAsInt == compressionTag {
				if uint8(dataFormatAsInt) == unsignedShortType {
					imageCompressionValue := utils.ConvertBytesToUInt16(ifd0Data[i+8], ifd0Data[i+9], imageTiffHeaderData.endianOrder)
					if imageCompressionValue == compressionNone {
						fmt.Printf("Compression -> None\n")
					}
				}
			} else if tagAsInt == photometricInterpretationTag {
				if uint8(dataFormatAsInt) == unsignedShortType {
					photometricInterpretationValue := utils.ConvertBytesToUInt16(ifd0Data[i+8], ifd0Data[i+9], imageTiffHeaderData.endianOrder)
					if photometricInterpretationValue == photometricInterpretationRGB {
						fmt.Printf("Photometric interpretation -> RGB\n")
					}
				}
			} else if tagAsInt == makeTag {
				if uint8(dataFormatAsInt) == asciiStringsType {
					imageMakeTagData := make([]byte, numOfElementsAsInt)
					file.Seek(int64(dataValueOrDataOffsetAsInt), os.SEEK_SET)
					file.Read(imageMakeTagData)
					fmt.Printf("Camera make -> %s\n", string(imageMakeTagData))
				}
			} else if tagAsInt == modelTag {
				if uint8(dataFormatAsInt) == asciiStringsType {
					imageModelTagData := make([]byte, numOfElementsAsInt)
					file.Seek(int64(dataValueOrDataOffsetAsInt), os.SEEK_SET)
					file.Read(imageModelTagData)
					fmt.Printf("Camera model -> %s\n", string(imageModelTagData))
				}
			} else if tagAsInt == stripOffsetsTag {
				if uint8(dataFormatAsInt) == unsignedLongType {
					fmt.Printf("Strip offsets -> %d\n", dataValueOrDataOffsetAsInt)
				}
			} else if tagAsInt == orientationTag {
				if uint8(dataFormatAsInt) == unsignedShortType {
					orientationTagData := utils.ConvertBytesToUInt16(ifd0Data[i+8], ifd0Data[i+9], imageTiffHeaderData.endianOrder)
					fmt.Printf("Orientation flag -> %d\n", orientationTagData)
				}
			} else if tagAsInt == samplesPerPixelTag {
				if uint8(dataFormatAsInt) == unsignedShortType {
					samplesPerPixelTagData := utils.ConvertBytesToUInt16(ifd0Data[i+8], ifd0Data[i+9], imageTiffHeaderData.endianOrder)
					fmt.Printf("Samples per pixel flag -> %d\n", samplesPerPixelTagData)
				}
			} else if tagAsInt == rowsPerStripTag {
				if uint8(dataFormatAsInt) == unsignedLongType {
					fmt.Printf("Rows per strip -> %d\n", dataValueOrDataOffsetAsInt)
				}
			} else if tagAsInt == stripByteCountsTag {
				if uint8(dataFormatAsInt) == unsignedLongType {
					fmt.Printf("Strip byte counts -> %d\n", dataValueOrDataOffsetAsInt)
				}
			} else if tagAsInt == xResolutionTag {
				if uint8(dataFormatAsInt) == unsignedRationalType {
					file.Seek(int64(dataValueOrDataOffsetAsInt), os.SEEK_SET)
					xResolutionTagData := make([]byte, 8)
					file.Read(xResolutionTagData)
					xResolutionTagNum := utils.ConvertBytesToUInt64(xResolutionTagData[0], xResolutionTagData[1], xResolutionTagData[2],
						xResolutionTagData[3], xResolutionTagData[4], xResolutionTagData[5],
						xResolutionTagData[6], xResolutionTagData[7], imageTiffHeaderData.endianOrder)
					fmt.Printf("X Resolution data -> %d\n", xResolutionTagNum)
				}
			} else if tagAsInt == yResolutionTag {
				if uint8(dataFormatAsInt) == unsignedRationalType {
					file.Seek(int64(dataValueOrDataOffsetAsInt), os.SEEK_SET)
					yResolutionTagData := make([]byte, 8)
					file.Read(yResolutionTagData)
					yResolutionTagNum := utils.ConvertBytesToUInt64(yResolutionTagData[0], yResolutionTagData[1], yResolutionTagData[2],
						yResolutionTagData[3], yResolutionTagData[4], yResolutionTagData[5],
						yResolutionTagData[6], yResolutionTagData[7], imageTiffHeaderData.endianOrder)
					fmt.Printf("Y Resolution data -> %d\n", yResolutionTagNum)
				}
			} else if tagAsInt == planarConfigurationTag {
				if uint8(dataFormatAsInt) == unsignedShortType {
					planarConfigurationTagData := utils.ConvertBytesToUInt16(ifd0Data[i+8], ifd0Data[i+9], imageTiffHeaderData.endianOrder)
					fmt.Printf("Planar configuration -> %d\n", planarConfigurationTagData)
				}
			} else if tagAsInt == resolutionUnitTag {
				if uint8(dataFormatAsInt) == unsignedShortType {
					resolutionUnitTagData := utils.ConvertBytesToUInt16(ifd0Data[i+8], ifd0Data[i+9], imageTiffHeaderData.endianOrder)
					fmt.Printf("Resolution unit -> %d\n", resolutionUnitTagData)
				}
			} else if tagAsInt == softwareTag {
				if uint8(dataFormatAsInt) == asciiStringsType {
					file.Seek(int64(dataValueOrDataOffsetAsInt), os.SEEK_SET)
					softwareTextData := make([]byte, numOfElementsAsInt)
					file.Read(softwareTextData)
					fmt.Printf("Software -> %s\n", string(softwareTextData))
				}
			}
		}
	}
	return nil
}

func readIfd(file *os.File, ifdOffset uint32, endianOrder utils.EndianOrder) []byte {
	ifdTagCountBytes := make([]byte, 2)
	file.Seek(int64(ifdOffset), os.SEEK_SET)
	file.Read(ifdTagCountBytes)

	ifdTagCount := utils.ConvertBytesToUInt16(ifdTagCountBytes[0], ifdTagCountBytes[1], endianOrder)

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

func getEdianOrder(header []byte) utils.EndianOrder {
	if len(header) >= 4 {
		var endianFlag uint16
		//add the bits to the 2 byte int and shove them to the left to make room for the other bits
		endianFlag |= uint16(header[0]) << 8
		endianFlag |= uint16(header[1])
		if endianFlag == 0x4d4d {
			return utils.BigEndian
		} else if endianFlag == 0x4949 {
			return utils.LittleEndian
		}
	}
	return utils.BigEndian
}

func getTiffHeader(header []byte) (tiffHeaderData, error) {
	tiffData := new(tiffHeaderData)
	tiffData.endianOrder = getEdianOrder(header)

	if len(header) >= 8 {

		var magicNum uint16
		if tiffData.endianOrder == utils.BigEndian {
			magicNum |= uint16(header[2]) | uint16(header[3])
		} else if tiffData.endianOrder == utils.LittleEndian {
			magicNum |= uint16(header[3]) | uint16(header[2])
		}

		tiffData.magicNum = magicNum
		tiffData.tiffOffset = utils.ConvertBytesToUInt32(header[4], header[5], header[6], header[7], tiffData.endianOrder)
	} else {
		return *tiffData, errors.New("Header incorrect length")
	}
	return *tiffData, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
