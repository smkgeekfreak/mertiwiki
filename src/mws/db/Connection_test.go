// Connection_test.go
package db

import (
	"log"
	"mws/util"
	"os"
	"testing"
)

//got := []string{dbconf.MigrationsDir, dbconf.Env, dbconf.Driver.Name, dbconf.Driver.OpenStr}
//	want := []string{"../../db-sample/migrations", "test", "postgres", "user=liam dbname=tester sslmode=disable"}

//	for i, s := range got {
//		if s != want[i] {
//			t.Errorf("Unexpected DBConf value. got %v, want %v", s, want[i])
//		}
//	}
func TestGetDBConfigTest(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	config, err := getDBConfig("test", "../../../deploy/db")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Getting Connection to Database: %s , %s", config.Driver.Name, config.Driver.OpenStr)
	}
}

func TestGetDBConfigDev(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	config, err := getDBConfig("dev", "../../../deploy/db")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Getting Connection to Database: %s , %s", config.Driver.Name, config.Driver.OpenStr)
	}
}

func TestGetDBConfigProd(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	config, err := getDBConfig("prod", "../../../deploy/db")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Getting Connection to Database: %s , %s", config.Driver.Name, config.Driver.OpenStr)
	}
}

func TestGetDBConfigNotFound(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	_, err := getDBConfig("shouldnteverfindthisconfig", "../../../deploy/db")
	if err != nil {
		log.Printf("Successfully did not find: shouldnteverfindthisconfig , %#v", err)
	} else {
		log.Fatalf("Should have failed to find (shouldnteverfindthisconfig) ")
	}
}

func TestLoadConnection(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	err := os.Setenv("WIKIENV", "test")
	if err != nil {
		log.Fatalf("Could not set WIKENV %#v", err)
	}
	err = os.Setenv("WIKICONFIGROOT", "../../../deploy/db")
	if err != nil {
		log.Fatalf("Could not set WIKICONFIGROOT %#v", err)
	}
	dbConn, err := loadConnection()
	if err != nil || dbConn == nil {
		t.Errorf("Get Connection unsuccessful %#v", err)
	}
	connErr := dbConn.Ping()
	if connErr != nil {
		t.Errorf("Connection query unsuccessful %#v", connErr)
		t.FailNow()
	}
}
func TestLoadConnectionBadEnv(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	err := os.Setenv("WIKIENV", "shouldnteverfindthisconfigd")
	if err != nil {
		log.Fatalf("Could not set WIKENVl %#v", err)
	}
	dbConn, err := loadConnection()
	if err != nil || dbConn == nil {
		log.Printf("Get Connection unsuccessful %#v", "shouldnteverfindthisconfig")
		return
	}
	t.Errorf("Should have failed to load (shouldnteverfindthisconfig)")

}

func TestGetConnection(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	err := os.Setenv("WIKIENV", "dev")
	if err != nil {
		t.Errorf("Could not set WIKENVl %#v", err)
	}
	dbConn, connErr := GetConnection()
	stmt, connErr := dbConn.Prepare("SELECT * From account")
	_, connErr = stmt.Exec()
	if connErr != nil {
		t.Errorf("Connection test unsuccessful %#v", connErr)
	}
}
