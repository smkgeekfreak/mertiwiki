#!/bin/sh
function processGoPackage {
 if [ $# -eq 0 ]; then
	echo ' Failed - Not enough arguments'
	exit 1
 fi
 if  [  -z $1  ] ; then
	echo ' Failed - Not enough arguments'
	exit 1
 fi
 echo "******************************************"
 echo '* PROCESSING PACKAGE '$1'    *'
 PACKAGE=$1
 cd $PROJPATH/src/$PACKAGE
 echo "******************************************"
 echo "Cleaning $PROJPATH/src/$PACKAGE library"
 echo "******************************************"
 go clean -x -i
}

source ./setpath.sh
echo "Set gopath = $PROJPATH"
processGoPackage 'mws/util' 'get'
processGoPackage 'mws/dto'
processGoPackage 'mws/db'
processGoPackage 'mws/model'
processGoPackage 'mws/mockmodel'
processGoPackage 'mws/resttest'
processGoPackage 'mws/microrest'
processGoPackage 'mws/server'
