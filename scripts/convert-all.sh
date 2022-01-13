#!/bin/sh

for f in data/in/*; do
    out=data/out/$(basename $f)
    echo "Covert $f -> $out"
    go run main.go --input "$f" --output "$out" --resize_width=300 --resize_height=300
done