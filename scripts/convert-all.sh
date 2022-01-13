#!/bin/sh

for f in data/in/*; do
    out=data/out/$(basename $f)
    echo "Convert $f -> $out"
    go run main.go --input "$f" --output "$out"
done