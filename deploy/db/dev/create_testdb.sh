#!/bin/sh
function checkReturnCode {
rc=$?
if [[ $rc != 0 ]] ; then
	echo $1 ' Failed (' $rc ')'	
	exit $rc
fi
echo $1 ' Succeeded (' $rc ')'	
}

source ./init_devdb.sh
psql -U wikiadmin_test -d MeritWiki_test -a -f create_new_user_id_seq.plsql
checkReturnCode 'create_new_user_id_seq.plsql'
psql -U wikiadmin_test -d MeritWiki_test -a -f create_usermodel.plsql
checkReturnCode 'create_usermodel.plsql'
psql -U wikiadmin_test -d MeritWiki_test -a -f create_userfunc.plsql
checkReturnCode 'create_userfunc.plsql'
