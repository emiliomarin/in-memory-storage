.PHONY: mod
## Install project dependencies
mod:
	@go mod tidy
	@go mod vendor

.PHONY: test
## Run tests. Usage: 'make test' Options: path=./some-path/... [and/or] func=TestFunctionName
test: ; $(info running tests…)
	@if [ -z $(path) ]; then \
		path='./...'; \
	else \
		path=$(path); \
	fi; \
	if [ -z $(func) ]; then \
		go test -v -failfast $$path; \
	else \
		go test -v -failfast -run $$func $$path; \
	fi;
