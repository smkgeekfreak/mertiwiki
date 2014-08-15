// Page_test.go
package dto

import (
	"encoding/json"
	"fmt"
	"log"
	"mws/util"
	"reflect"
	"strings"
	"testing"
)

func TestEntityRelationShipExtraValue(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	er := &EntityRelationship{
		RelId1:   1,
		RelId2:   1,
		RelName1: "dto.Revision",
		RelName2: "dto.Tag",
		RelType1: reflect.TypeOf(Revision{}),
		RelType2: reflect.TypeOf(Tag{}),
		Status:   SetStatus(StatusType(0)),
		Values:   map[string]interface{}{"rating": 100},
	}
	//
	if er == nil || er.Values == nil || len(er.Values) == 0 {
		t.Errorf("%s Values not set = %v", funcName, er)
		count.FailCount++
		return
	}
	for k, v := range er.Values {
		log.Printf("%s Relationship has values ( %s, %s, %v)", funcName, k, reflect.TypeOf(v), reflect.ValueOf(v).Interface())
		if reflect.TypeOf(v) != reflect.TypeOf(100) {
			t.Errorf("%s Values type not correct = %v", funcName, v)
			count.FailCount++
			return
		}
		var getVal int
		getVal = reflect.ValueOf(v).Interface().(int)
		log.Printf("%s Successfully got value from entity relationship %v", funcName, getVal)
	}

	gotJSON, err := json.Marshal(er)
	if err != nil {
		t.Error("%s Could not construct JSON %v", funcName, err)
		count.FailCount++
		return
	}
	log.Printf("%s Successfully created entity relationship %v", funcName, string(gotJSON))

	count.SuccessCount++
}

func TestEntityRelationShipExtraValues(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	er := &EntityRelationship{
		RelId1:   1,
		RelId2:   1,
		RelName1: "dto.Revision",
		RelName2: "dto.Tag",
		RelType1: reflect.TypeOf(Revision{}),
		RelType2: reflect.TypeOf(Tag{}),
		Status:   SetStatus(StatusType(0)),
		Values:   map[string]interface{}{"rating": 100, fmt.Sprintf("%s", reflect.TypeOf(Tag{})): Tag{Id: 1}},
	}
	//
	if er == nil || er.Values == nil || len(er.Values) == 0 {
		t.Errorf("%s Values not set = %v", funcName, er)
		count.FailCount++
		return
	}
	for k, v := range er.Values {
		log.Printf("%s Relationship has values ( %s, %s, %v)", funcName, k, reflect.TypeOf(v), reflect.ValueOf(v).Interface())

	}

	gotJSON, err := json.Marshal(er)
	if err != nil {
		t.Error("%s Could not construct JSON %v", funcName, err)
		count.FailCount++
		return
	}
	log.Printf("%s Successfully created entity relationship %v", funcName, string(gotJSON))

	count.SuccessCount++
}

func TestEntityRelationShipNoExtraValues(t *testing.T) {
	funcName := util.GetCallerName()
	log.Printf("%s ---------------- Starting  Test ---------------------", funcName)
	er := &EntityRelationship{
		RelId1:   1,
		RelId2:   1,
		RelName1: "dto.Revision",
		RelName2: "dto.Tag",
		RelType1: reflect.TypeOf(Revision{}),
		RelType2: reflect.TypeOf(Tag{}),
		Status:   SetStatus(StatusType(0)),
	}
	//
	if er == nil || er.Values != nil || len(er.Values) != 0 {
		t.Errorf("%s Values  set = %v", funcName, er)
		count.FailCount++
		return
	}

	gotJSON, err := json.Marshal(er)
	if err != nil {
		t.Error("%s Could not construct JSON %v", funcName, err)
		count.FailCount++
		return
	}

	if hasValues := strings.Contains(string(gotJSON), "Values"); hasValues {
		t.Error("%s Values still present in %v", funcName, string(gotJSON))
		count.FailCount++
		return
	}
	log.Printf("%s Successfully created entity relationship %v", funcName, string(gotJSON))

	count.SuccessCount++
}
