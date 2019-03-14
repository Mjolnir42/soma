all: validate

validate:
	@go build ./...
	@go vet .
	@go tool vet -shadow .
	@golint . | grep -v -E '(or be unexported|comment on exported)' || true
	@ineffassign .
	@gofmt -w .
	@unparam .
