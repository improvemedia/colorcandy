thrift:
	thrift --gen go --out . candy.thrift

build: thrift
	go get github.com/tools/godep
	godep go build -o colorcandy bin/main.go

run: build
	./colorcandy
