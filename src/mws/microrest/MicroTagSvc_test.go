// MicroTagSvc_test.go
package microrest

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	"mws/mockmodel"
	"mws/model"
	"mws/resttest"
	"mws/util"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var mts *MicroTagSvc

func init() {
	funcName := util.GetCallerName()
	log.Printf("%s ........................ Init Micro Tag Service Tests ...............................", funcName)
	mts = NewMicroTagSvc()
}

func TestMicroTagSvcCreateTag(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	funcName := util.GetCallerName()
	mts.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/tag")
	mts.MicroSvc.AddModel(mockmodel.MockTagModel{})
	timeStr := time.Now().Format(time.RFC850)
	tests := []resttest.RESTTestContainer{
		{
			Desc:       "CreateTag",
			Handler:    restful.DefaultContainer.ServeHTTP,
			Path:       mts.MicroSvc.GetFullPath(),
			Method:     "POST",
			JSONParams: `{"Name":"Creating a RESTTest Tag ` + timeStr + `",  "Description":" Test Tag Desc "} `,
			Status:     http.StatusCreated,
			MatchVal: dto.Tag{
				Name:   "Creating a RESTTest Tag " + timeStr,
				Status: dto.StatusDetail{dto.INITIALIZED, fmt.Sprint(dto.StatusType(dto.INITIALIZED))},
			},
			MatchFields:     []string{"Name", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}

	mod, _ := mts.MicroSvc.getModel(mts.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("%s ----- Testing with Model ( %s )---------------------", funcName, reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	// Test w/Actual database
	mts.MicroSvc.AddModel(model.TagModel{})
	mod, _ = mts.MicroSvc.getModel(mts.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("%s ----- Testing with Model ( %s )---------------------", funcName, reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}
func TestMicroTagSvcUpdateTag(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mts.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/tag")
	timeStr := time.Now().Format(time.RFC850)
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "UpdateTag",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mts.MicroSvc.GetFullPath(),
				"%2": "2",
			},
			Method:     "PUT",
			JSONParams: `{"Name":"Updating a Test Tag ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:     http.StatusOK,
			MatchVal: dto.Tag{
				Id:     2,
				Status: dto.StatusDetail{dto.DEACTIVATED, fmt.Sprint(dto.StatusType(dto.DEACTIVATED))},
			},
			MatchFields:     []string{"Id", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
		{
			Desc:    "UpdateTagFailsGlobalThreshold",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mts.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:          "PUT",
			JSONParams:      `{"Name":"Updating a Test Tag ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:          http.StatusForbidden,
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 2500000,
		},
		{
			Desc:    "UpdateTagFailsThresholdNotProvided",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mts.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:          "PUT",
			JSONParams:      `{"Name":"Updating a Test Tag ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:          http.StatusForbidden,
			PreAuthId:       3,
			OwnThreshold:    200,
			GlobalThreshold: 300,
		},
	}
	mts.MicroSvc.AddModel(model.TagModel{})
	mod, _ := mts.MicroSvc.getModel(mts.MicroSvc.GetFullPath())
	rel := mts.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{}))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	//// Test w/Actual database
	log.SetPrefix("")
}
