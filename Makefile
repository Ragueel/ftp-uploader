OUTPUT ?= ftp-uploader


.PHONY: clean
clean:
	rm -rf builds
.PHONY: build
build: clean
	go build -o builds/$(OUTPUT) $(FLAGS) ./cmd/ftp-uploader/main.go