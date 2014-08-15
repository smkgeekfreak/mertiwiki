#!/bin/sh
function checkReturnCode {
rc=$?
if [[ $rc != 0 ]] ; then
	echo $1 ' Failed (' $rc ')'	
	exit $rc
fi
echo $1 ' Succeeded (' $rc ')'	
}

psql -U wikiadmin_test -d MeritWiki_test -a -f drop_acct.plsql
checkReturnCode 'drop_acct.plsql'
psql -U postgres -a -f drop_testdbrole.plsql
checkReturnCode 'drop_testdbrole.plsql'
