#!/bin/sh
if [ $# -eq 0 ]; then
	echo ' Failed - Not enough arguments'
	exit 1
fi
if  [  -z $1  ] ; then
	echo ' Failed - Not enough arguments'
	exit 1
fi

echo 'Setting environment to ' $1 
source ./setpath.sh
echo "Set gopath = $PROJPATH"
echo "*******************"
echo "Starting MWS Server"
echo "*******************"
cd ../src/mws/server

server -wikiEnv=$1
