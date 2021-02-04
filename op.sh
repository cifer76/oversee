#!/usr/bin/env bash

remote=3.6.94.102

if [ $1 = "pub" ]; then
    rsync --exclude=.git/ --exclude=oversee --exclude=oversee.pem --exclude=*.sw[po] -avzhe "ssh -i ./oversee.pem" . ubuntu@$remote:/home/ubuntu/oversee/
elif [ $1 = 'serve' ]; then
	go build
    ./oversee
else
	echo 'unsupported command'
fi
