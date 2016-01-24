all: test

clean:
	@rm -f master
	@rm -Rf pkg
	@rm -Rf .tmp

deps:
	@go get -d -t ./...

dist: deps
	@gox -output="pkg/dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		--osarch "darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64"

build: deps
	@go build

install: deps
	@go install

test: deps
	@go test -v

.PHONY: clean, deps, dist, build, install, test
