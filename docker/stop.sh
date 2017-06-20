#!/usr/bin/env bash

echo "Stopping..."
docker stop kafka
docker stop zookeeper

if [[ $1 == "--remove" ]] ; then
    echo "Removing..."
    docker rm kafka
    docker rm zookeeper
fi

if [[ $1 == "--rmi" ]] ; then
    echo "Removing..."
    docker rmi kafka -f
    docker rmi zookeeper -f
fi
