// MicroRevisionSvc_test.go
package microrest

import (
	//"fmt"
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	"mws/mockmodel"
	"mws/model"
	"mws/resttest"
	"mws/util"
	"net/http"
	//"os"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var mrs *MicroRevisionSvc

func init() {

	log.Printf("........................ Init Micro TagRevision Service Test ...............................")
	mrs = NewMicroRevisionSvc()
}

func TestMicroRevisionSvcCreateRevision(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mrs.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/revision")
	mrs.MicroSvc.AddModel(model.RevisionModel{})
	timeStr := time.Now().Format(time.RFC850)
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "CreateRevision",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/",
			Sub: map[string]string{
				"%1": mrs.MicroSvc.GetFullPath(),
			},
			Method:     "POST",
			JSONParams: `{"SectionId":2, "Content":"Creating a Test Rev ` + timeStr + `"} `,
			Status:     http.StatusCreated,
			MatchVal: dto.Revision{
				SectionId: 2,
				Content:   "Creating a Test Rev  " + timeStr,
				Status:    dto.StatusDetail{dto.INITIALIZED, fmt.Sprint(dto.StatusType(dto.INITIALIZED))},
			},
			MatchFields:     []string{"SectionId", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}
	mod, _ := mrs.MicroSvc.getModel(mrs.MicroSvc.GetFullPath())
	rel := mrs.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{}))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	//// Test w/Actual database
	log.SetPrefix("")
}

func TestMicroRevisionSvcUpdateRevision(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mrs.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/revision")
	mrs.MicroSvc.AddModel(model.RevisionModel{})
	timeStr := time.Now().Format(time.RFC850)
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "UpdateRevision",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mrs.MicroSvc.GetFullPath(),
				"%2": "7",
			},
			Method:     "PUT",
			JSONParams: `{"Content":"Updating a Test Rev ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:     http.StatusOK,
			MatchVal: dto.Revision{
				Id:     7,
				Status: dto.StatusDetail{dto.DEACTIVATED, fmt.Sprint(dto.StatusType(dto.DEACTIVATED))},
			},
			MatchFields:     []string{"Id", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
		{
			Desc:    "UpdateRevisionFailsGlobalThreshold",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mrs.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:          "PUT",
			JSONParams:      `{"Content":"Updating a Test Rev ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:          http.StatusUnauthorized,
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 2500000,
		},
		{
			Desc:    "UpdateRevisionFailsThresholdNotProvided",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mrs.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:          "PUT",
			JSONParams:      `{"Content":"Updating a Test Rev ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:          http.StatusPreconditionFailed,
			PreAuthId:       3,
			OwnThreshold:    -1,
			GlobalThreshold: -1,
		},
	}
	mod, _ := mrs.MicroSvc.getModel(mrs.MicroSvc.GetFullPath())
	rel := mrs.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{}))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	//// Test w/Actual database
	log.SetPrefix("")
}

func TestMicroRevisionSvcCreateTagRelationship(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mrs.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/revision")
	mrs.MicroSvc.AddModel(model.RevisionModel{})
	mrs.MicroSvc.AddRelationship(reflect.TypeOf(dto.Tag{}), mockmodel.MockEntityRelationshipModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "CreateTagRelationship",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2/tag/%3",
			Sub: map[string]string{
				"%1": mrs.MicroSvc.GetFullPath(),
				"%2": "1",
				"%3": "1",
			},
			Method:          "PUT",
			Status:          http.StatusOK,
			PreAuthId:       7,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}
	mod, _ := mrs.MicroSvc.getModel(mrs.MicroSvc.GetFullPath())
	rel := mrs.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{}))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	//// Test w/Actual database
	log.SetPrefix("")
}
func TestMicroRevisionSvcRetrieveTagRelationship(t *testing.T) {
	mrs.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/revision")
	mrs.MicroSvc.AddModel(model.RevisionModel{})
	mrs.MicroSvc.AddRelationship(reflect.TypeOf(dto.Tag{}), mockmodel.MockEntityRelationshipModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "RetrieveTagRelationship",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2/tag/%3",
			Sub: map[string]string{
				"%1": mrs.MicroSvc.GetFullPath(),
				"%2": "2",
				"%3": "2",
			},
			Method:          "GET",
			Status:          http.StatusOK,
			PreAuthId:       7,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}
	mod, _ := mrs.MicroSvc.getModel(mrs.MicroSvc.GetFullPath())
	rel := mrs.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{}))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}

//Had to comment out this test b/c it failes to do trying to insert
// a duplicate row when the test is run multiple multiple times
//func TestMicroRevisionSvcCreateRateRelationship(t *testing.T) {
//	log.SetPrefix(util.GetCallerName() + " : ")
//	mrs.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/revision")
//	mrs.MicroSvc.AddModel(model.RevisionModel{})
//	mrs.MicroSvc.AddRelationship(reflect.TypeOf(dto.Rating{}), model.RatingModel{})

//	tests := []resttest.RESTTestContainer{
//		{
//			Desc:    "CreateTagRelationship",
//			Handler: restful.DefaultContainer.ServeHTTP,
//			Path:    "%1/%2/rate/%3",

//			Sub: map[string]string{
//				"%1": mrs.MicroSvc.GetFullPath(),
//				"%2": "1",
//				"%3": "17",
//			},
//			Method:          "PUT",
//			Status:          http.StatusOK,
//			PreAuthId:       9,
//			OwnThreshold:    0,
//			GlobalThreshold: 0,
//		},
//	}
//	mod, _ := mrs.MicroSvc.getModel(mrs.MicroSvc.GetFullPath())
//	rel := mrs.MicroSvc.getRelationship(reflect.TypeOf(dto.Rating{}))
//	log.Println("----------------------------------------------------------------------------")
//	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
//	log.Printf("----- Testing with Relationship ( %s )---------------------", reflect.TypeOf(rel))
//	log.Println("----------------------------------------------------------------------------")
//	resttest.RunTestSet(t, tests)
//	//// Test w/Actual database
//	log.SetPrefix("")
//}
