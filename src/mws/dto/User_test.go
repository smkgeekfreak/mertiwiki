// User_test.go
package dto

import (
	"encoding/json"
	"fmt"
	"log"
	//	"mws/resttest"
	"testing"
	//"time"
	"mws/util"
)

func TestCreateUser(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	//now := time.Now();

	wantUserJSON := `{"Id":0,"Name":"sCamelCase","Email":"test@test.com","Status":{"StatusCode":0,"Desc":"Initalized"}}`
	u := &User{
		Name:   "sCamelCase",
		Email:  "test@test.com",
		Status: StatusDetail{INITIALIZED, fmt.Sprint(StatusType(INITIALIZED))},
	}
	//u.Created = &jsonTime{time.Now(), time.RFC3339}
	//u.Modified = &jsonTime{time.Now(), time.RFC3339}
	gotJSON, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
		return
	}
	//
	if wantUserJSON != string(gotJSON) {
		t.Errorf("Received = %s\n Expected %s", gotJSON, wantUserJSON)
		return
	}
	log.Printf("Successfully created user %s", string(gotJSON))

}
