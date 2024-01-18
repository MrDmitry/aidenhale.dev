BIN_DIR=bin

build: blog sitemap robots

dist:
	mkdir -p ./web/dist

dist/css: dist
	mkdir -p ./web/dist/css

css: dist/css
	tailwindcss -c ./tailwind.config.js -i ./web/src/css/style.css -o ./web/dist/css/style.css --minify

blog:
	go build -v -o ./${BIN_DIR}/blog ./cmd/blog/main.go

sitemap: dist
	go build -v -o ./${BIN_DIR}/sitemap_builder ./cmd/sitemap/main.go
	./${BIN_DIR}/sitemap_builder $(SITEMAP_ARGS) > ./web/dist/sitemap.xml

robots: dist
	go build -v -o ./${BIN_DIR}/robots_builder ./cmd/robots/main.go
	./${BIN_DIR}/robots_builder $(ROBOTS_ARGS) > ./web/dist/robots.txt

run: build css
	./${BIN_DIR}/blog $(BLOG_ARGS)

test:
	go test -v ./...

clean:
	rm -rf ./${BIN_DIR}
	rm -rf ./web/dist
