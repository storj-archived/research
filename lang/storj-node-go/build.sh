#!/bin/bash
echo "Building dockerfile for storj-node-go"
docker build -t storj-node-go -f ./dockerfiles/Dockerfile .
echo "Container built"

