build:
	mkdir -p bin && go build -o bin/ps4-20th-tool main.go

package:
	goxc -d=builds

test-brute: build
	TEST=1 bin/ps4-20th-tool brute

test-server:
	PORT=4010 go run test/server.go

.PHONY: build package test-brute test-server
