.PHONY: build run clean showcase

build:
	go build -o bin/showcase showcase/main.go
	@for dir in examples/*/; do \
		example=$$(basename $$dir); \
		go build -o bin/$$example $$dir/main.go; \
	done

run:
	go run showcase/main.go

clean:
	rm -rf bin/

showcase:
	go run showcase/main.go