BINARY := resume-cli

.PHONY: build test lint clean setup-env

build:
	go build -o $(BINARY) .

test:
	go test ./...

lint:
	test -z "$$(gofmt -l .)"
	go vet ./...

setup-env:
	@if [ -f .env ]; then \
		echo ".env already exists, skip"; \
	else \
		cp .env.example .env; \
		echo "created .env from .env.example"; \
	fi

clean:
	rm -f $(BINARY) $(BINARY).exe
