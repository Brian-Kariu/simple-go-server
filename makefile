.PHONY: build

build:
	go build -o ./bin/simple-go-server ./cmd/.

.PHONY: run

run: build
	./bin/simple-go-server