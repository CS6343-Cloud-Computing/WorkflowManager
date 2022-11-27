#! /bin/bash

for i in {0..5..1}
do
	go run *.go workflow submit graph2.yaml
done
