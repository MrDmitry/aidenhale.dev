#!/bin/bash -ex
tailwindcss -c tailwind.config.js -i web/src/css/style.css -o web/dist/css/style.css --minify
go run cmd/main.go
