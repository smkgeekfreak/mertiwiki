// zzzz_test.go
package model

import (
	"log"
	"mws/util"
	"os"
	"testing"
)

var (
	count                                   util.TestCounter
	createId                                int64
	tagId, revId, pageId, sectionId, userId int64
)

func init() {
	log.Printf("........................ Init Model Tests ...............................")
	err := os.Setenv("WIKIENV", "test")
	if err != nil {
		log.Fatalf("Could not set WIKENVl %#v", err)
	}
	err = os.Setenv("WIKICONFIGROOT", "../../../deploy/db")
	if err != nil {
		log.Fatalf("Could not set WIKICONFIGROOT %#v", err)
	}
}

func TestZZZZ(t *testing.T) {
	total := count.FailCount + count.SuccessCount
	log.Printf("-----Completed (%d) Tests ---------------------", total)
	log.Printf("Success Count = %d Fail Count %d", count.SuccessCount, count.FailCount)
	log.Printf("-------------------------------------------------------------------------------")
}
