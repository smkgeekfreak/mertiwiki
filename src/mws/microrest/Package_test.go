// Package_test.go
package microrest

import (
	"log"
	"mws/util"
	"os"
)

var (
	//count util.TestCounter
	tagRevId, revRevId int64
)

func init() {
	wd, _ := os.Getwd()
	log.Printf("WORKING DIR %s", wd)

	config := util.GetConfig()
	err := os.Setenv("WIKIENV", "test")
	if err != nil {
		log.Fatalf("Could not set WIKIENV %#v, %v", err, config)
	}
	err = os.Setenv("WIKICONFIGROOT", "../../../deploy/db")
	if err != nil {
		log.Fatalf("Could not set WIKICONFIGROOT %#v", err)
	}
}
