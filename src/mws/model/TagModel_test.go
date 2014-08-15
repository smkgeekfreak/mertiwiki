// TagModel_test.go
// Testing the TagModel implmentation for intracting with
// the actual data model / virtual data layer / database
package model

import (
	"log"
	"mws/dto"
	"mws/util"
	//"reflect"
	"testing"
)

func TestTagModelCreate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	log.SetPrefix(util.GetCallerName() + " : ")
	getThisTag := dto.Tag{
		Name:        "test tag",
		Description: "Content for my first tag <b>has tags</b>",
		Status:      dto.SetStatus(dto.INITIALIZED),
	}
	smodel := TagModel{}
	operation, _ := GetVDLOperation(smodel, "create")
	gotTag, retCode, err := smodel.Create(operation, &getThisTag)
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
	tagId = (*gotTag).Id
	log.Printf("%s Successful for %#v", funcName, gotTag)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestTagModelFindById(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestTagModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	smodel := TagModel{}
	operation, _ := GetVDLOperation(smodel, "findbyid")
	tag_id := int64(tagId)
	gotSection, retCode, err := smodel.FindById(operation, tag_id)
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

func TestTagModelFindByIdNotFound(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	log.SetPrefix(util.GetCallerName() + " : ")
	smodel := TagModel{}
	operation, _ := GetVDLOperation(smodel, "findbyid")
	got, retCode, err := smodel.FindById(operation, 98888)
	if err == nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if got != nil || retCode == 0 {
		t.Errorf("%s recCode = 0, wanted non-zero return", funcName)
		count.FailCount++
		return
	}
	log.Printf("%s Successful got return Code = %d", funcName, retCode)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestTagModelUpdate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestTagModelCreate(t)
	updateId := int64(tagId)
	tmodel := TagModel{}

	upThisTag := dto.Tag{
		Id:          updateId,
		Name:        "updated tag",
		Description: "Content for my updated tag <b>has tags</b> and updated with more<b> tags </b>",
		Status:      dto.SetStatus(dto.BANNED),
	}
	operation, _ := GetVDLOperation(tmodel, "update")
	tag_id := int64(tagId)
	got, retCode, err := tmodel.Update(operation, tag_id, &upThisTag)
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
	tagId = (*got).Id
	log.Printf("%s Successful for %#v", funcName, got)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestTagModelUpdateStatus(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestTagModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	smodel := TagModel{}
	operation, _ := GetVDLOperation(smodel, "updatestatus")
	tag_id := int64(tagId)
	gotTag, retCode, err := smodel.UpdateStatus(operation, tag_id, dto.BANNED)
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
	tagId = (*gotTag).Id
	log.Printf("%s Successful for %#v", funcName, gotTag)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestTagModelDelete(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestTagModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	smodel := TagModel{}
	operation, _ := GetVDLOperation(smodel, "delete")
	tag_id := int64(tagId)
	gotTag, retCode, err := smodel.Delete(operation, tag_id)
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
	tagId = (*gotTag).Id
	log.Printf("%s Successful for %#v", funcName, gotTag)
	count.SuccessCount++
	log.SetPrefix("")
}
