#!/bin/sh

for f in data/in/*; do
    out=data/out/$(basename $f)
    echo "Covert $f -> $out"
    go run main.go --input "$f" --output "$out"
done