// SectionModel_test.go
// Testing the SectionModel implmentation for intracting with
// the actual data model / virtual data layer / database
package model

import (
	"log"
	"mws/dto"
	"mws/util"
	"testing"
)

func TestSectionModelCreate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestPageModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	getThisSection := dto.Section{
		Name:     "Rate My First Section",
		PageId:   pageId,
		OrderNum: 1,
		//Status:   dto.SetStatus(dto.INITIALIZED),
	}
	smodel := SectionModel{}
	operation, _ := GetVDLOperation(smodel, "create")
	author_id := int64(userId)
	gotSection, retCode, err := smodel.Create(operation, author_id, &getThisSection)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := retCode, 0; got != want {
		t.Errorf("%s got  %d, wanted  %d", funcName, got, want)
		count.FailCount++
		return
	}
	sectionId = (*gotSection).Id
	log.Printf("%s Successful for %#v", funcName, gotSection)
	count.SuccessCount++
	log.SetPrefix("")
}
func TestSectionModelUpdate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestSectionModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	getThisSection := dto.Section{
		Name:     "Update My First Section",
		PageId:   pageId,
		OrderNum: 1,
		Status:   dto.SetStatus(dto.ACTIVE),
	}
	smodel := SectionModel{}
	operation, _ := GetVDLOperation(smodel, "update")
	section_id := int64(sectionId)
	gotSection, retCode, err := smodel.Update(operation, section_id, &getThisSection)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := retCode, 0; got != want {
		t.Errorf("%s got  %d, wanted  %d", funcName, got, want)
		count.FailCount++
		return
	}
	sectionId = (*gotSection).Id
	log.Printf("%s Successful for %#v", funcName, gotSection)
	count.SuccessCount++
	log.SetPrefix("")
}
func TestSectionModelDelete(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestSectionModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	smodel := SectionModel{}
	operation, _ := GetVDLOperation(smodel, "delete")
	section_id := int64(sectionId)
	gotSection, retCode, err := smodel.Delete(operation, section_id)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := retCode, 0; got != want {
		t.Errorf("%s got  %d, wanted  %d", funcName, got, want)
		count.FailCount++
		return
	}
	sectionId = (*gotSection).Id
	log.Printf("%s Successful for %#v", funcName, gotSection)
	count.SuccessCount++
	log.SetPrefix("")
}
func TestSectionModelFindById(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestSectionModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	smodel := SectionModel{}
	operation, _ := GetVDLOperation(smodel, "findbyid")
	section_id := int64(sectionId)
	gotSection, retCode, err := smodel.FindById(operation, section_id)
	if err != nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got, want := retCode, 0; got != want {
		t.Errorf("%s got  %d, wanted  %d", funcName, got, want)
		count.FailCount++
		return
	}
	sectionId = (*gotSection).Id
	log.Printf("%s Successful for %#v", funcName, gotSection)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestSectionModelFindByIdNotFound(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestSectionModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	smodel := SectionModel{}
	operation, _ := GetVDLOperation(smodel, "findbyid")
	gotSection, retCode, err := smodel.FindById(operation, 98888)
	if err == nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if gotSection != nil || retCode == 0 {
		t.Errorf("%s recCode = 0, wanted non-zero return", funcName)
		count.FailCount++
		return
	}
	log.Printf("%s Successful got return Code = %d", funcName, retCode)
	count.SuccessCount++
	log.SetPrefix("")
}

