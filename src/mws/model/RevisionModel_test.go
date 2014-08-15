// RevisionModel_test.go
// Testing the RevisionModel implmentation for intracting with
// the actual data model / virtual data layer / database
package model

import (
	"log"
	"mws/dto"
	"mws/util"
	"testing"
)

//func TestRevisionModelFindByAuthor(t *testing.T) {
//	funcName := util.GetCallerName()
//	log.Printf("---------------- Staroting  Test ( %s )---------------------", funcName)
//	TestRevisionModelCreate(t)
//	log.SetPrefix(funcName + " : ")
//	searchId := revId
//	smodel := RevisionModel{}
//	operation, _ := GetVDLOperation(smodel, "findbyauthor")
//	revision, retCode, err := smodel.FindByAuthor(operation, searchId)
//	if err != nil {
//		t.Errorf(" %s Failed: %v", funcName, err.Error)
//		count.FailCount++
//		return
//	}
//	if revision == nil || retCode != 0 {
//		t.Errorf("%s Broken: got id , wanted id ", funcName)
//		count.FailCount++
//		return
//	}
//	if len(*revision) <= 0 {
//		t.Errorf("%s Broken: no rows returned ", funcName)
//		count.FailCount++
//		return
//	}
//	log.Printf("%s Successful ", funcName)
//	count.SuccessCount++
//	log.SetPrefix("")
//}

func TestRevisionModelFindById(t *testing.T) {
	funcName := util.GetCallerName()
	TestRevisionModelCreate(t)
	log.SetPrefix(funcName + " : ")
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	searchId := int64(revId)
	smodel := RevisionModel{}
	operation, _ := GetVDLOperation(smodel, "findbyid")
	revision, retCode, err := smodel.FindById(operation, searchId)
	if err != nil {
		t.Errorf("RevisionModel FindById Failed: %v", err.Error)
		count.FailCount++
		return
	}

	if want, got := searchId, revision.Id; got != want || retCode != 0 {
		t.Errorf("RevisionModel FindById Broken: got id %s, wanted id %s", got, want)
		count.FailCount++
		return
	}
	log.Printf("%s Successful for %#v", funcName)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestRevisionModelCreate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	TestUserModelCreate(t)
	TestSectionModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	getThisRevision := dto.Revision{
		SectionId: sectionId,
		Content:   "Content for my first revision <b>has tags</b>",
	}
	smodel := RevisionModel{}
	if gotRevision, retCode, err := smodel.Create("create_revision", userId, &getThisRevision); gotRevision == nil || err != nil {
		//if gotRevision == nil || err != nil {
		log.Fatalf(" RevisionModel Create Failed: %#v", err)
		count.FailCount++
		return
	} else {

		if got, want := retCode, 0; got != want {
			t.Errorf("RevisionModelCreate Broken: got  %d, wanted  %d", got, want)
			count.FailCount++
			return
		}
		revId = (*gotRevision).Id
		log.Printf("RevisionModel Create Successful for %#v", gotRevision)
		count.SuccessCount++
	}
	log.SetPrefix("")
}

func TestRevisionModelUpdate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	TestRevisionModelCreate(t)
	log.SetPrefix(funcName + " : ")
	getThisRevision := dto.Revision{
		Id:      revId,
		Content: "Content for my first revision <b>has tags</b>",
		Status:  dto.SetStatus(dto.INITIALIZED),
	}
	smodel := RevisionModel{}
	if gotRevision, retCode, err := smodel.Update("update_revision", revId, &getThisRevision); gotRevision == nil || err != nil {
		//if gotRevision == nil || err != nil {
		log.Fatalf(" %s Failed: %#v", funcName, err)
		count.FailCount++
		return
	} else {

		if got, want := retCode, 0; got != want {
			t.Errorf("%s Broken: got  %d, wanted  %d", funcName, got, want)
			count.FailCount++
			return
		}
		revId = (*gotRevision).Id
		log.Printf("%s Successful for %#v", funcName, gotRevision)
		count.SuccessCount++
	}
	log.SetPrefix("")
}

func TestRevisionModelFindByIdNotFound(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	searchId := int64(99991999)
	smodel := RevisionModel{}
	operation, _ := GetVDLOperation(smodel, "findbyid")
	revision, _, err := smodel.FindById(operation, searchId)
	if err != nil {
		count.SuccessCount++
		return

	}
	t.Errorf("RevisionModel RetrieveById  found %v: %s", revision, err)
	count.FailCount++
	return

	count.SuccessCount++
	log.SetPrefix("")
}

func TestRevisionModelDelete(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	TestRevisionModelCreate(t)
	log.SetPrefix(funcName + " : ")
	smodel := RevisionModel{}
	if gotRevision, retCode, err := smodel.Delete("update_revisionstatus", revId); gotRevision == nil || err != nil {
		log.Fatalf(" %s Failed: %#v", funcName, err)
		count.FailCount++
		return
	} else {

		if got, want := retCode, 0; got != want || gotRevision.Status.StatusCode != dto.DELETED {
			t.Errorf("%s Broken: got  %d, wanted  %d", funcName, got, want)
			count.FailCount++
			return
		}
		revId = (*gotRevision).Id
		log.Printf("%s Successful for %#v", funcName, gotRevision)
		count.SuccessCount++
	}
	log.SetPrefix("")
}
func TestRevisionModelUpdateStatus(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("---------------- Starting  Test ( %s )---------------------", funcName)
	TestRevisionModelCreate(t)
	log.SetPrefix(funcName + " : ")
	smodel := RevisionModel{}
	if gotRevision, retCode, err := smodel.UpdateStatus("update_revisionstatus", revId, dto.BANNED); gotRevision == nil || err != nil {
		log.Fatalf(" %s Failed: %#v", funcName, err)
		count.FailCount++
		return
	} else {

		if got, want := retCode, 0; got != want || gotRevision.Status.StatusCode != dto.BANNED {
			t.Errorf("%s Broken: got  %d, wanted  %d", funcName, got, want)
			count.FailCount++
			return
		}
		revId = (*gotRevision).Id
		log.Printf("%s Successful for %#v", funcName, gotRevision)
		count.SuccessCount++
	}
	log.SetPrefix("")
}
