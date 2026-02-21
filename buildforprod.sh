#!/bin/bash
npx tailwindcss -i ./tailwind.css -o ./static/css/tailwind.css --minify
mkdir -p dist
cp -r static dist/
cp -r templates dist/
go build -o dist/app main.go
