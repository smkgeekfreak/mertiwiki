#!/bin/sh
cd ../
PROJPATH=$(pwd)
export PROJPATH 
echo "PROJPATH=$PROJPATH"
GOPATH=$PROJPATH
export GOPATH
echo "GOPATH=$GOPATH"
export PATH=$PATH:$GOPATH/bin
cd build
