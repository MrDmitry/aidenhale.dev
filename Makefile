BIN_DIR=bin
BINARY_NAME=blog

css:
	tailwindcss -c ./tailwind.config.js -i ./web/src/css/style.css -o ./web/dist/css/style.css --minify

build:
	go build -v -o ./${BIN_DIR}/${BINARY_NAME} ./cmd/main.go

run: build css
	./${BIN_DIR}/${BINARY_NAME} $(ARGS)

test:
	go test -v ./...

clean:
	rm -rf ./${BIN_DIR}
	rm -rf ./web/dist
