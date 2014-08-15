#!/bin/sh
function checkReturnCode {
 rc=$?
 if [[ $rc != 0 ]] ; then
 	echo $1 ' Failed (' $rc ')'	
 	exit $rc
 fi
 echo $1 ' Succeeded (' $rc ')'	
}

# 
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
 echo '          * PROCESSING PACKAGE '$1'    *'
 PACKAGE=$1
 cd $PROJPATH/src/$PACKAGE
 echo "******************************************"
 echo "Building  $PROJPATH/src/$PACKAGE  "
 echo "******************************************"
 go get
 checkReturnCode 'go get '$PACKAGE
 go build
 checkReturnCode 'go build '$PACKAGE
 go test
 checkReturnCode 'go test '$PACKAGE
 go install
 checkReturnCode 'go install '$PACKAGE
}


source ./setpath.sh
echo "Set gopath = $PROJPATH"
processGoPackage 'mws/util' 'get'
processGoPackage 'mws/resttest'
processGoPackage 'mws/dto'
processGoPackage 'mws/db'
processGoPackage 'mws/db/pgreflect'
processGoPackage 'mws/model'
processGoPackage 'mws/mockmodel'
processGoPackage 'mws/microrest'
processGoPackage 'mws/server'
