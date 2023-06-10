#!/bin/sh

echo "Starting proto to struct..."
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

# For Grom in proto, run once
# 1. Run this cd $GOPATH/src/github.com/infobloxopen/protoc-gen-gorm && make vendor && make install

protoc -I=. -I=$GOPATH/src --go_out=.. --gorm_out=engine=postgres:.. *.proto

echo "Completed proto to struct..."
