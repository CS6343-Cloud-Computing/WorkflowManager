#!/bin/bash

a=0
input="./out1.txt"
while IFS= read -r line
do
	docker run -d --name $line --env kafkaserver=10.176.128.170:9090 --env workflow=$line --env limit=200 lcds
done < "$input"