#!/bin/bash

a=0
input="./out1.txt"
while IFS= read -r line
do
	if (( "$a"% 3 == 0 ))
	then
		docker run -d --name $line --env kafkaserver=10.176.128.170:9090 --env workflow=$line cds
	elif (( "$a"% 5 == 0 ))
	then
		docker run -d --name $line --env kafkaserver=10.176.128.170:9090 --env workflow=$line wds
	else
		docker run -d --name $line --env kafkaserver=10.176.128.170:9090 --env workflow=$line tds
	fi
    let "a += 1"
done < "$input"