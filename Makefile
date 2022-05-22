.PHONY: build-linux build-mac build-win test report fmt clean

BUILDCMD_LINUX=env GOOS=linux GOARCH=amd64 go build -v
BUILDCMD_MAC=env GOOS=darwin GOARCH=amd64 go build -v
BUILDCMD_WIN=env GOOS=windows GOARCH=amd64 go build -v

build-linux:
	$(BUILDCMD_LINUX) -o gowc cmd/gowc/*.go

build-mac:
	$(BUILDCMD_MAC) -o gowc cmd/gowc/*.go

build-win:
	$(BUILDCMD_WIN) -o gowc.exe cmd\gowc\*.go

test:
	go test `go list ./...` -v -cover -count=1

report:
	go test `go list ./...` -cover -covermode=count -coverprofile=coverage.out

fmt:
	! go fmt ./... 2>&1 | tee /dev/tty | read

clean:
	go clean ./...