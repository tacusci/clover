
$VERSION="0.0.2a"
$GOOS_ENVS=@("windows", "linux", "darwin")
$GOARCHS=@("386", "amd64")
$WINEXT=""

foreach($i in $GOOS_ENVS) {
    if ($i -eq "windows") { $WINEXT = ".exe" } else { $WINEXT = "" }
    foreach($j in $GOARCHS) {
        $OUTPUTFOLDER = "clover-v$VERSION-$i-$j"
        $OUTPUTPATH = "bin\$OUTPUTFOLDER\clover$WINEXT"
        $BUILDSTR = "env GOOS=$i GOARCH=$j go build -o $OUTPUTPATH"
        Invoke-Expression "& $BUILDSTR"
        Compress-Archive -Path $OUTPUTPATH -DestinationPath "bin\$OUTPUTFOLDER.zip"
    }
}