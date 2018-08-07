#!/bin/bash

VERSION="0.0.2a"
GOOS_ENVS=(windows linux darwin freebsd)
GOARCHS=(386 amd64 arm)
WINEXT=""

echo "Clover build script - Building clover $VERSION for release..."

for i in ${GOOS_ENVS[@]};
do
    echo "Building clover-v$VERSION for $i"
    if [ $i == "windows" ]; then WINEXT=".exe"; else WINEXT=""; fi
    for j in ${GOARCHS[@]};
    do
        if [ $j == "arm" ] && [ $i != "linux" ]; then continue; fi
        OUTPUTFOLDER="clover-v$VERSION-$i-$j"
        OUTPUTPATH="bin/v$VERSION/$OUTPUTFOLDER/clover-v$VERSION$WINEXT"
        BUILDSTR="env GOOS=$i GOARCH=$j go build -o $OUTPUTPATH"
        $BUILDSTR
        
        echo "Built $OUTPUTPATH"

        zip "/$OUTPUTPATH" "/$OUTPUTPATH.zip"
    done
done
