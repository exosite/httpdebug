all: linux osx windows

linux: build/linux-amd64/httpdebug

osx: build/osx-amd64/httpdebug

windows: build/win-amd64/httpdebug.exe

# Linux Build
build/linux-amd64/httpdebug: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ github.com/exosite/httpdebug
# OS X Build
build/osx-amd64/httpdebug: main.go
	GOOS=darwin GOARCH=amd64 go build -o $@ github.com/exosite/httpdebug
# Windows Build
build/win-amd64/httpdebug.exe: main.go
	GOOS=windows GOARCH=amd64 go build -o $@ github.com/exosite/httpdebug

clean:
	rm -f build/linux-amd64/httpdebug
	rm -f build/osx-amd64/httpdebug
	rm -f build/win-amd64/httpdebug.exe
	rm -f *~

.PHONY: all clean linux osx windows
