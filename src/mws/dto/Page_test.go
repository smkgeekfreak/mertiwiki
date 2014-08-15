// Page_test.go
package dto

import (
	"encoding/json"
	"fmt"
	"log"
	//	"mws/resttest"
	"mws/util"
	"testing"
)

func TestCreatePage(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	wantUserJSON := `{"Id":12345,"Title":"Test New Page","Status":{"StatusCode":0,"Desc":"Initalized"}}`
	u := &Page{
		Id:     12345,
		Title:  "Test New Page",
		Status: StatusDetail{INITIALIZED, fmt.Sprint(StatusType(INITIALIZED))},
	}
	//	u.SetPath()
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
	log.Printf("Successfully created page %s", string(gotJSON))
}

func TestCreateSection(t *testing.T) {
	log.Printf("---------------- Starting  Test ( %s )---------------------", util.GetCallerName())
	wantUserJSON := `{"Id":23456,"PageId":12345,"AuthorId":123,"Name":"New Section Name","OrderNum":1,"Status":{"StatusCode":0,"Desc":"Initalized"}}`
	sec := &Section{
		Id:       23456,
		PageId:   12345,
		AuthorId: 123,
		Name:     "New Section Name",
		Status:   StatusDetail{INITIALIZED, fmt.Sprint(StatusType(INITIALIZED))},
		OrderNum: 1,
	}
	gotJSON, err := json.Marshal(sec)
	if err != nil {
		t.Error(err)
		return
	}
	//
	if wantUserJSON != string(gotJSON) {
		t.Errorf("Received = %s\n Expected %s", gotJSON, wantUserJSON)
		return
	}
	log.Printf("Successfully created section %s", string(gotJSON))

}
