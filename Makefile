.PHONY: build build-all build-macos build-macos-intel build-windows clean

build:
	go build -o video-arranger .

build-all: build-macos build-macos-intel build-windows

build-macos:
	GOOS=darwin GOARCH=arm64 go build -o dist/video-arranger-macos-arm64 .

build-macos-intel:
	GOOS=darwin GOARCH=amd64 go build -o dist/video-arranger-macos-amd64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o dist/video-arranger-windows-amd64.exe .

clean:
	rm -rf dist/ video-arranger
