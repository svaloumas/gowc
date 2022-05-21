.PHONY: build-linux build-mac build-win

BUILDCMD_LINUX=env GOOS=linux GOARCH=amd64 go build -v
BUILDCMD_MAC=env GOOS=darwin GOARCH=amd64 go build -v
BUILDCMD_WIN=env GOOS=windows GOARCH=amd64 go build -v

build-linux:
	$(BUILDCMD_LINUX) -o gowc *.go

build-mac:
	$(BUILDCMD_MAC) -o gowc *.go

build-win:
	$(BUILDCMD_WIN) -o gowc.exe *.go
