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

func TestTagRevisionModelCreate(t *testing.T) {
	testName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", testName)
	TestTagModelCreate(t)
	TestRevisionModelCreate(t)
	trmodel := TagRevisionModel{}
	vdlOperation := "addtagrevision"
	if gotER, retCode, err := trmodel.CreateRelationship(vdlOperation, tagId, revId); gotER == nil || err != nil {
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
		log.Printf("%s Create Successful for relaionship (%d, %d, %v, %v, %v)", testName, gotER.RelId1, gotER.RelId2, gotER.Status, gotER.RelType1, gotER.RelType2)
		count.SuccessCount++
		return
	}
}
func TestTagRevisionModelCreateFailedFK(t *testing.T) {
	testName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", testName)
	vdlOperation := "addtagrevision"
	trmodel := TagRevisionModel{}
	if gotER, retCode, err := trmodel.CreateRelationship(vdlOperation, 888888, 999999); gotER != nil && err == nil {
		//if gotER == nil || err != nil {
		t.Errorf(" %s Failed: %#v", testName, err)
		count.FailCount++
		return
	} else {

		if got, want := retCode, -1; got != want {
			t.Errorf("%s Broken: got  %d, wanted  %d", testName, got, want)
			count.FailCount++
			return
		}

		log.Printf("%s Successful", testName)
		count.SuccessCount++
	}
}

func TestTagRevisionModelRetrieve(t *testing.T) {
	testName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", testName)
	TestTagRevisionModelCreate(t)
	tagId := int64(tagId)
	revId := int64(revId)
	trmodel := TagRevisionModel{}
	if gotER, retCode, err := trmodel.RetrieveRelationship(revId, tagId); gotER == nil || err != nil {
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
		log.Printf("%s Create Successful for relaionship (%d, %d, %v, %v, %v)", testName, gotER.RelId1, gotER.RelId2, gotER.Status, gotER.RelType1, gotER.RelType2)
		count.SuccessCount++
		return
	}
}
