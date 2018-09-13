# vim: set ft=make ffs=unix fenc=utf8:
# vim: set noet ts=4 sw=4 tw=72 list:
#
SOMAVER != cat `git rev-parse --show-toplevel`/VERSION
BRANCH != git rev-parse --symbolic-full-name --abbrev-ref HEAD
GITHASH != git rev-parse --short HEAD

# if GOPATH contains a symbolic link, then PWD and GOPATH can diverge
# leading to go package loader errors
GOPATH != /bin/realpath ${GOPATH}
.export GOPATH

all: install

install: install_freebsd install_linux

install_freebsd: generate
	@env GOOS=freebsd GOARCH=amd64 go install -ldflags "-X main.somaVersion=$(SOMAVER)-$(GITHASH)/$(BRANCH)" ./...

install_linux: generate
	@env GOOS=linux GOARCH=amd64 go install -ldflags "-X main.somaVersion=$(SOMAVER)-$(GITHASH)/$(BRANCH)" ./...

generate:
	@go generate ./cmd/...

sanitize: build check

sanitize-full: build check-full

check: vet ineffassign misspell

check-full: vet ineffassign misspell lint coroner

build:
	@echo "Building ...."
	@go build ./...

vet:
	@echo "Running 'go vet' ...."
	@go vet ./cmd/eye/
	@go vet ./cmd/somad/
	@go vet ./cmd/soma/
	@go vet ./cmd/somadbctl/
	@go vet ./lib/auth/
	@go vet ./lib/proto/
	@go vet ./internal/adm/
	@go vet ./internal/cmpl/
	@go vet ./internal/config/
	@go vet ./internal/db/
	@go vet ./internal/help/
	@go vet ./internal/msg/
	@go vet ./internal/perm/
	@go vet ./internal/rest/
	@go vet ./internal/soma/
	@go vet ./internal/stmt/
	@go vet ./internal/super/
	@go vet ./internal/tree/
	@echo "Running 'go vet -shadow' ...."
	@go tool vet -shadow ./cmd/eye/
	@go tool vet -shadow ./cmd/somad/
	@go tool vet -shadow ./cmd/soma/
	@go tool vet -shadow ./cmd/somadbctl/
	@go tool vet -shadow ./lib/auth/
	@go tool vet -shadow ./lib/proto/
	@go tool vet -shadow ./internal/adm/
	@go tool vet -shadow ./internal/cmpl/
	@go tool vet -shadow ./internal/config/
	@go tool vet -shadow ./internal/db/
	@go tool vet -shadow ./internal/help/
	@go tool vet -shadow ./internal/msg/
	@go tool vet -shadow ./internal/perm/
	@go tool vet -shadow ./internal/rest/
	@go tool vet -shadow ./internal/soma/
	@go tool vet -shadow ./internal/stmt/
	@go tool vet -shadow ./internal/super/
	@go tool vet -shadow ./internal/tree/

ineffassign:
	@echo "Running 'ineffassign' ...."
	@ineffassign ./cmd
	@ineffassign ./lib
	@ineffassign ./internal

misspell:
	@echo "Running 'misspell' ...."
	@misspell ./cmd
	@misspell ./lib
	@misspell ./internal
	@misspell ./docs

coroner:
	@echo "Running 'codecoroner' ...."
	@codecoroner funcs ./cmd/... ./lib/... ./internal/...

lint:
	@echo "Running 'golint' ...."
	@golint ./cmd/... | grep -v 'should have comment or be unexported'
	@golint ./internal/... | grep -v 'should have comment or be unexported'
	@golint ./lib/... | grep -v 'should have comment or be unexported'
