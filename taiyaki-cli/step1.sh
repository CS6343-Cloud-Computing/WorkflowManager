#! /bin/bash

for i in {0..5..1}
do
	if (( $i % 3 == 0 ))
	then
		go run *.go workflow submit cityStats.yml
	elif (( $i % 5 == 0 ))
	then
		go run *.go workflow submit wordDefinition.yml
	else
		go run *.go workflow submit textprocessing.yml
	fi
done
