BUILD_DIR=cmd
BINARY_NAME=blog

css:
	tailwindcss -c tailwind.config.js -i web/src/css/style.css -o web/dist/css/style.css --minify

build:
	go build -C ${BUILD_DIR} -v -o ${BINARY_NAME}

run: build css
	./${BUILD_DIR}/${BINARY_NAME} $(ARGS)
