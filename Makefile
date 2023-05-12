default: test test-coverage test-integration

.PHONY: test
test:
	go test -v  ./...

# runs coverage tests and generates the coverage report
.PHONY: test-coverage
test-coverage:
	go test ./... -v -coverpkg=./...

# runs integration tests
.PHONY: test-integration
test-integration:
	go test ./... -tags=integration ./...