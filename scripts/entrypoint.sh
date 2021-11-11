#!/bin/bash

set -e

# Setup requirements
export GOBIN=$(pwd)/bin
export PATH=$GOBIN:$PATH
make gen
sleep 1
CompileDaemon --build="go build -o bin/goer ./cmd/goer" --exclude-dir=".git" --exclude-dir="**/**/test" --exclude-dir="**/**/gomock*" --command=./bin/goer -verbose
