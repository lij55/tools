#!/bin/bash -x

for target in linux windows darwin; do
	for i in $(ls *.go); do
		GOOS=$target go  build  -o ${i%".go"}_$target $i
	done;
done 