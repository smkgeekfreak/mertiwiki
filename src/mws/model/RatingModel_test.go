// TagRevisionModel_test.go
// Testing the TagRevisionModel implmentation for intracting with
// the actual data model / virtual data layer / database
package model

import (
	"log"
	//	"mws/dto"
	"mws/util"
	"testing"
)

func TestRatingModelCreate(t *testing.T) {
	testName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", testName)
	TestUserModelCreate(t)
	//userId := createId
	TestRevisionModelCreate(t)
	//revId := createId
	trmodel := RatingModel{}
	if gotER, retCode, err := trmodel.CreateRelationship("addrating", userId, revId, 100); gotER == nil || err != nil {
		//if gotER == nil || err != nil {
		log.Fatalf("%s Failed: %#v", testName, err)
		count.FailCount++
		return
	} else {

		if got, want := retCode, 0; got != want {
			t.Errorf("%s Broken: got  %d, wanted  %d", testName, got, want)
			count.FailCount++
			return
		}
		//		createId = (*gotER).Id
		log.Printf("%s Create Successful for relationship (%d)", testName, gotER.AccountId)
		count.SuccessCount++
		return
	}
}
