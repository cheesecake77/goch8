
.PHONY: run lint

run:
	go run .

lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run --config=.golangci.yml
