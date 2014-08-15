// PageModel_test.go
// Testing the PageModel implmentation for intracting with
// the actual data model / virtual data layer / database
package model

import (
	"log"
	"mws/dto"
	"mws/util"
	"testing"
)

func init() {
}

func TestPageModelCreate(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	TestUserModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	getThisPage := dto.Page{
		Title: "Rate My First Page",
	}
	pmodel := PageModel{}
	operation, _ := GetVDLOperation(pmodel, "create")
	author_id := int64(userId)
	gotPage, retCode, err := pmodel.Create(operation, author_id, &getThisPage)
	if err != nil {
		log.Fatalf(" PageModel Create Failed: %#v", err)
		count.FailCount++
		return
	}

	if got, want := retCode, 0; got != want {
		t.Errorf("PageModelCreate Broken: got  %d, wanted  %d", got, want)
		count.FailCount++
		return
	}
	pageId = (*gotPage).Id
	log.Printf("PageModel Create Successful for %#v", gotPage)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestPageModelFindById(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	TestPageModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	searchId := int64(pageId)
	pmodel := PageModel{}
	operation, _ := GetVDLOperation(pmodel, "findbyid")
	page, retCode, err := pmodel.FindById(operation, searchId)
	if err != nil || retCode != 0 {
		t.Errorf("PageModel RetrieveById Failed: %v", err.Error)
		count.FailCount++
		return
	}

	if want, got := searchId, page.Id; got != want {
		t.Errorf("PageModel RetrieveById Broken: got id %s, wanted id %s", got, want)
		count.FailCount++
		return
	}

	count.SuccessCount++
	log.SetPrefix("")
}
func TestPageModelRetrieveByIdNotFound(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	searchId := int64(99991999)
	pmodel := PageModel{}
	operation, _ := GetVDLOperation(pmodel, "findbyid")
	page, retCode, err := pmodel.FindById(operation, searchId)
	if err != nil || retCode != 0 {
		count.SuccessCount++
		return

	}
	t.Errorf("PageModel RetrieveById  found %v: %s", page, err)
	count.FailCount++
	return

	count.SuccessCount++
	log.SetPrefix("")
}

func TestPageModelDelete(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	TestPageModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	searchId := int64(pageId)
	pmodel := PageModel{}
	operation, _ := GetVDLOperation(pmodel, "delete")
	page, retCode, err := pmodel.Delete(operation, searchId)
	if err != nil || retCode != 0 {
		t.Errorf("PageModel Delete Failed: %v", err)
		count.FailCount++
		return
	}

	if want, got := 0, retCode; got != want || page.Status.StatusCode != dto.DELETED {
		t.Errorf("PageModel Delete Broken: got %s, wanted  %s", got, want)
		count.FailCount++
		return
	}

	log.Printf("PageModel Delete Successful for %#v", page)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestPageModelzzzz(t *testing.T) {
	TestZZZZ(t)
}
