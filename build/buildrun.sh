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

source ./build.sh
checkReturnCode
echo "*******************"
echo "Starting MWS Server"
echo "*******************"
server -wikiEnv=$1
