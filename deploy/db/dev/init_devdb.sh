#!/bin/sh
function checkReturnCode {
rc=$?
if [[ $rc != 0 ]] ; then
	echo $1 ' Failed (' $rc ')'	
	exit $rc
fi
echo $1 ' Succeeded (' $rc ')'	
}

if [ $# -eq 0 ]; then
	echo ' Failed - Not enough arguments, please provide environment to init'
	exit 1
fi
if  [  -z $1  ] ; then
	echo ' Failed - Not enough arguments, please provide environment to init'
	exit 1
fi

psql -U postgres -a -f create_testrole.plsql
checkReturnCode 'create_testrole.plsql'
psql -U postgres -a -f create_new_testdb.plsql
checkReturnCode 'create_new_testdb.plsql'
psql MeritWiki_test wikiadmin_test << EOF
CREATE LANGUAGE plpgsql
EOF
goose -env=$1 -path="../." status
#read -t 10 -p "Hit ENTER or wait ten seconds"
goose -env=$1 -path="../." up 
checkReturnCode 'goose -env=$1 -path="../." up'
#goose -env=$1 -path="../." down 
#checkReturnCode 'goose -env=$1 -path="../." down'
goose -env=$1 -path="../." status
