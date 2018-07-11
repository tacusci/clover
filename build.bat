@echo off
set VERSION=0.0.2a
set GOOS_ENVS=(windows linux darwin)
set GOARCHS=(386 amd64)
for %%i in %GOOS_ENVS% do (
	for %%j in %GOARCHS% do (
		env GOOS=%%i GOARCH=%%j go build -o "bin\clover-v%VERSION%-%%i-%%j\clover"
	)
)