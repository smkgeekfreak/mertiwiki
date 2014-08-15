// MicroSectionSvc_test.go
// MicroSectionSvc_test.go
package microrest

import (
	//"fmt"
	"fmt"
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	"mws/model"
	"mws/resttest"
	"mws/util"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var mss *MicroSectionSvc

func init() {

	log.Printf("........................ Init Micro TagSection Service Test ...............................")
	mss = NewMicroSectionSvc()
}

func TestMicroSectionSvcCreateSection(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mss.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/section")
	timeStr := time.Now().Format(time.RFC850)
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "CreateSection",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/",
			Sub: map[string]string{
				"%1": mss.MicroSvc.GetFullPath(),
			},
			Method:     "POST",
			JSONParams: `{"PageId":2, "Name":"Creating a Test Section` + timeStr + `"} `,
			Status:     http.StatusCreated,
			MatchVal: dto.Section{
				PageId: 2,
				Name:   "Creating a Test Section  " + timeStr,
				Status: dto.StatusDetail{dto.INITIALIZED, fmt.Sprint(dto.StatusType(dto.INITIALIZED))},
			},
			MatchFields:     []string{"PageId", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}
	mss.MicroSvc.AddModel(model.SectionModel{})
	mod, _ := mss.MicroSvc.getModel(mss.MicroSvc.GetFullPath())
	rel := mss.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{}))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	//// Test w/Actual database
	log.SetPrefix("")
}

func TestMicroSectionSvcUpdateSection(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mss.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/section")
	timeStr := time.Now().Format(time.RFC850)
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "UpdateSection",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mss.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:     "PUT",
			JSONParams: `{"Name":"Updating a Test Section ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:     http.StatusOK,
			MatchVal: dto.Section{
				Id:     4,
				Status: dto.StatusDetail{dto.DEACTIVATED, fmt.Sprint(dto.StatusType(dto.DEACTIVATED))},
			},
			MatchFields:     []string{"Id", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
		{
			Desc:    "UpdateSectionFailsGlobalThreshold",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mss.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:          "PUT",
			JSONParams:      `{"Name":"Updating a Test Section ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:          http.StatusUnauthorized,
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 2500000,
		},
		{
			Desc:    "UpdateSectionFailsThresholdNotProvided",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mss.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:          "PUT",
			JSONParams:      `{"Name":"Updating a Test Section ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:          http.StatusPreconditionFailed,
			PreAuthId:       3,
			OwnThreshold:    100,
			GlobalThreshold: -1,
		},
	}
	mss.MicroSvc.AddModel(model.SectionModel{})
	mod, _ := mss.MicroSvc.getModel(mss.MicroSvc.GetFullPath())
	rel := mss.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{}))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	//// Test w/Actual database
	log.SetPrefix("")
}
