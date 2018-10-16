#!/bin/bash -x

rm -f buildtmp
mkdir buildtmp

for target in linux windows darwin; do
	for i in $(ls *.go); do
		GOOS=$target go  build  -o buildtmp/${i%".go"}_$target $i
	done;
done 
