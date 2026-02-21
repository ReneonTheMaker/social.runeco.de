#!/bin/bash
#!/bin/bash
set -e

echo "Building Tailwind..."
npx tailwindcss -i ./tailwind.css -o ./static/css/tailwind.css --minify

echo "Preparing dist..."
rm -rf dist
mkdir -p dist

cp -r static dist/
cp -r views dist/
cp config.ini dist/

echo "Building Go binary..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/app

echo "Done."
