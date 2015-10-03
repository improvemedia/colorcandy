thrift:
	thrift --gen go --out . thrift/candy.thrift

build: thrift
	go get github.com/tools/godep
	godep go build -o colorcandy bin/main.go

run: build
	./colorcandy

.PHONY: thrift
