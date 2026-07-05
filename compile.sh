#!/usr/bin/env sh

MINIFIED_DIR=minified

# Clean and create the minified directory
rm -rf $MINIFIED_DIR
mkdir $MINIFIED_DIR

# Minify JS & CSS
npx esbuild client/main.js --bundle --minify --outfile="$MINIFIED_DIR/main.js" || cp client/main.js $MINIFIED_DIR/main.js
npx esbuild client/style.css --minify --outfile="$MINIFIED_DIR/style.css" || cp client/style.css $MINIFIED_DIR/style.css

# Minify HTML
npx html-minifier-terser client/index.html -o $MINIFIED_DIR/index.html || cp client/index.html $MINIFIED_DIR/index.html

# Minify SVG
npx svgo client/icon.svg -o $MINIFIED_DIR/icon.svg --multipass || cp client/icon.svg $MINIFIED_DIR/icon.svg

src_size=$(du -sb client | cut -f1)
dst_size=$(du -sb $MINIFIED_DIR | cut -f1)
saved=$((src_size - dst_size))

echo "Minification saved $saved bytes ($src_size B -> $dst_size B)"

# Compile
go build -o fils

echo "Compilation finished"