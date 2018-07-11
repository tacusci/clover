
$VERSION="0.0.2a"
$GOOS_ENVS=@("windows", "linux", "darwin", "freebsd")
$GOARCHS=@("386", "amd64", "arm")
$WINEXT=""

foreach($i in $GOOS_ENVS) {
    if ($i -eq "windows") { $WINEXT = ".exe" } else { $WINEXT = "" }
    foreach($j in $GOARCHS) {
        if ($j -eq "arm" -and $i -ne "linux") { continue }
        $OUTPUTFOLDER = "clover-v$VERSION-$i-$j"
        $OUTPUTPATH = "bin\v$VERSION\$OUTPUTFOLDER\clover$WINEXT"
        $BUILDSTR = "env GOOS=$i GOARCH=$j go build -o $OUTPUTPATH"
        Invoke-Expression "& $BUILDSTR"
        Compress-Archive -Path $OUTPUTPATH -DestinationPath "bin\v$VERSION\$OUTPUTFOLDER.zip"
    }
}