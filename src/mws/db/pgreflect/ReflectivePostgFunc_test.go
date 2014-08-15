// ReflectivePostgFunc_test.go
package pgreflect

import (
	"log"
	//	"mws/db"
	"errors"
	"mws/util"
	"os"
	"strings"
	"testing"
)

var (
	count util.TestCounter
)

func init() {
	log.Printf("........................ Init Model Tests ...............................")
	err := os.Setenv("WIKIENV", "test")
	if err != nil {
		log.Fatalf("Could not set WIKENVl %#v", err)
	}
	err = os.Setenv("WIKICONFIGROOT", "../../../../deploy/db")
	if err != nil {
		log.Fatalf("Could not set WIKICONFIGROOT %#v", err)
	}
}

func TestGetPgFuncInfoWithLowerCase(t *testing.T) {
	funcName := util.GetCallerName()
	opName := "update_revision"
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	pgfunc, err := GetPgFuncInfo(opName)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := pgfunc.Name, opName; got != want {
		t.Errorf("%s Broken: got  %s, wanted  %s", funcName, got, want)
		count.FailCount++
		return
	}
	log.Printf("%s Successful for %#v", funcName, pgfunc)
	count.SuccessCount++
}

func TestGetPgFuncInfoWithUpperCase(t *testing.T) {
	funcName := util.GetCallerName()
	opName := "UPDATE_REVISION"
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	pgfunc, err := GetPgFuncInfo(opName)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := pgfunc.Name, strings.ToLower(opName); got != want {
		t.Errorf("%s Broken: got  %s, wanted  %s", funcName, got, want)
		count.FailCount++
		return
	}
	log.Printf("%s Successful for %#v", funcName, pgfunc)
	count.SuccessCount++
}

func TestGetPgFuncInfoWithCamelCase(t *testing.T) {
	funcName := util.GetCallerName()
	opName := "Update_Revision"
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	pgfunc, err := GetPgFuncInfo(opName)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := pgfunc.Name, strings.ToLower(opName); got != want {
		t.Errorf("%s Broken: got  %s, wanted  %s", funcName, got, want)
		count.FailCount++
		return
	}
	log.Printf("%s Successful for %#v", funcName, pgfunc)
	count.SuccessCount++
}

func TestGetPgFuncInfoWithSpace(t *testing.T) {
	funcName := util.GetCallerName()
	opName := "Update_ Revision"
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	pgfunc, err := GetPgFuncInfo(opName)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := pgfunc.Name, strings.ToLower(strings.Replace(opName, " ", "", -1)); got != want {
		t.Errorf("%s Broken: got  %s, wanted  %s", funcName, got, want)
		count.FailCount++
		return
	}
	log.Printf("%s Successful for %#v", funcName, pgfunc)
	count.SuccessCount++
}

func TestParseArgTypes(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	argmap := parsePgFuncArgTypes(3, "title character varying, path character varying, status integer,")

	if argmap == nil {
		log.Fatalf("%s Failed: %#v", funcName, errors.New("Arguments not parsed"))
		count.FailCount++
		return
	}

	if got, want := len(argmap), 3; got != want {
		t.Errorf("%s Broken: got  %s args, wanted  %s args", funcName, got, want)
		count.FailCount++
		return
	}
	log.Printf("PASS:-- %s for %#v", funcName, argmap)
	count.SuccessCount++
}

func TestParseArgTypesExtraArgs(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	argmap := parsePgFuncArgTypes(3, "title character varying, path character varying, status integer,should not parse, these extras")

	if argmap == nil {
		log.Fatalf("%s Failed: %#v", funcName, errors.New("Arguments not parsed"))
		count.FailCount++
		return
	}

	if got, want := len(argmap), 3; got != want {
		t.Errorf("%s Broken: got  %s args, wanted  %s args", funcName, got, want)
		count.FailCount++
		return
	}
	log.Printf("PASS:-- %s for %#v", funcName, argmap)
	count.SuccessCount++
}

func TestParseArgTypesNameWithNoType(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	argmap := parsePgFuncArgTypes(2, "title character varying, path ,")

	if argmap != nil {
		log.Fatalf("%s Failed: %#v", funcName, errors.New("Arguments not parsed"))
		count.FailCount++
		return
	}

	if got, want := len(argmap), 0; got != want {
		t.Errorf("%s Broken: got  %d args, and only  %d args was valid %#v", funcName, got, 1, argmap)
		count.FailCount++
		return
	}
	log.Printf("PASS:-- %s for %#v", funcName, argmap)
	count.SuccessCount++
}

func Test1(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	//parsePgFuncArgTypes(3, "title character varying, path character varying, status integer,")
	//parsePgFuncArgTypes(3, "uid integer, deletedstatus integer, ")
	pgfunc, _ := GetPgFuncInfo("update_revision")
	retMap, _, _ := pgfunc.VariadicScan(int64(1), "variadic body", 2)
	log.Printf("%s Returned %v", funcName, retMap)
	pgfunc, _ = GetPgFuncInfo("findbyauthor_revision")
	retMap, _, _ = pgfunc.VariadicScan(int64(2))
	log.Printf("%s Returned %v", funcName, retMap)
	//pgfunc, _ = GetPgFuncInfo("updatetagrevisionstatus")
	//retMap, _, _ = pgfunc.VariadicScan(int64(1), int64(2), 1)
	//log.Printf("%s Returned %v", funcName, retMap)
	//pgfunc, _ = GetPgFuncInfo("addtag")
	//retMap, _, _ = pgfunc.VariadicScan("variadic tag name", "variadic tag desc", 2)
	//log.Printf("%s Returned %v", funcName, retMap)

}

func TestZZZZ(t *testing.T) {
	total := count.FailCount + count.SuccessCount
	log.Printf("-----Completed (%d) Tests ---------------------", total)
	log.Printf("Success Count = %d Fail Count %d", count.SuccessCount, count.FailCount)
	log.Printf("-------------------------------------------------------------------------------")
}
