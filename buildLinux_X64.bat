mkdir uSafeServerExe
copy .\settingServer.json .\uSafeServerExe
copy .\server.pfx .\uSafeServerExe
copy .\restartServer.sh .\uSafeServerExe
for /R %%f in (.\uSafeServerExe\*.sh) do DOS2UNIX.EXE "%%f"


SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
SET VERSION=1.0
SET YEAR=%DATE:~2,2%
SET MONTH=%DATE:~5,2%
SET DAY=%DATE:~8,2%
SET BUILDVERSION=%VERSION%.%YEAR%%MONTH%%DAY%
go build -ldflags "-s -w -X 'main.BuildVersion=%BUILDVERSION%'" -o uSafeServerExe\ID_uSafe_Linux_X64


