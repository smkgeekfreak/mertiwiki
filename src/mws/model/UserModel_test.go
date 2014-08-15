// Testing the UserModel implmentation for intracting with
// the actual data model / virtual data layer / database
package model

import (
	"log"
	"mws/dto"
	"mws/util"
	"testing"
)

func TestUserModelCreate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	log.SetPrefix(util.GetCallerName() + " : ")
	getThisUser := dto.User{
		Name:  "Rate My First User",
		Email: "Email@mail.com",
		//Status:   dto.SetStatus(dto.INITIALIZED),
	}
	umodel := UserModel{}
	operation, _ := GetVDLOperation(umodel, "create")
	gotUser, retCode, err := umodel.Create(operation, &getThisUser)
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
	userId = (*gotUser).Id
	log.Printf("%s Successful for %#v", funcName, gotUser)
	count.SuccessCount++
	log.SetPrefix("")
}
func TestUserModelUpdate(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestUserModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	getThisUser := dto.User{
		Email:  "NewEmail@mail.com",
		Status: dto.SetStatus(dto.ACTIVE),
	}
	umodel := UserModel{}
	operation, _ := GetVDLOperation(umodel, "update")
	user_id := int64(userId)
	gotUser, retCode, err := umodel.Update(operation, user_id, &getThisUser)
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
	userId = (*gotUser).Id
	log.Printf("%s Successful for %#v", funcName, gotUser)
	count.SuccessCount++
	log.SetPrefix("")
}
func TestUserModelDelete(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestUserModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	umodel := UserModel{}
	operation, _ := GetVDLOperation(umodel, "delete")
	user_id := int64(userId)
	gotUser, retCode, err := umodel.Delete(operation, user_id)
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
	userId = (*gotUser).Id
	log.Printf("%s Successful for %#v", funcName, gotUser)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestUserModelFindById(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestUserModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	umodel := UserModel{}
	operation, _ := GetVDLOperation(umodel, "findbyid")
	user_id := int64(userId)
	gotUser, retCode, err := umodel.FindById(operation, user_id)
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
	userId = (*gotUser).Id
	log.Printf("%s Successful for %#v", funcName, gotUser)
	count.SuccessCount++
	log.SetPrefix("")
}

func TestUserModelFindByIdNotFound(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	TestUserModelCreate(t)
	log.SetPrefix(util.GetCallerName() + " : ")
	umodel := UserModel{}
	operation, _ := GetVDLOperation(umodel, "findbyid")
	gotUser, retCode, err := umodel.FindById(operation, 98888)
	if err == nil {
		log.Fatalf("%s Failed: %#v", funcName, err)
		count.FailCount++
		return
	}

	if gotUser != nil || retCode == 0 {
		t.Errorf("%s recCode = 0, wanted non-zero return", funcName)
		count.FailCount++
		return
	}
	log.Printf("%s Successful got return Code = %d", funcName, retCode)
	count.SuccessCount++
	log.SetPrefix("")
}
