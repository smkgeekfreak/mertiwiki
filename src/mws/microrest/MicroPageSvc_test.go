// MicroPageSvc_test.go
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

var mps *MicroPageSvc

func init() {
	log.Printf("........................ Init Micro Page Service Test ...............................")
	mps = NewMicroPageSvc()
}

func TestMicroPageSvcAddPage(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + ":")
	mps.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/page")
	mps.MicroSvc.AddModel(mockmodel.MockPageModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "AddPage",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/",
			Sub: map[string]string{
				"%1": mps.MicroSvc.GetFullPath(),
			},
			Method:     "POST",
			JSONParams: `{"Title":"Rate My First Page"}`,
			Status:     http.StatusCreated,
			MatchVal: dto.Page{
				Title:  "Rate My First Page",
				Status: dto.StatusDetail{dto.INITIALIZED, fmt.Sprint(dto.StatusType(dto.INITIALIZED))},
			},
			MatchFields:     []string{"Title", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    -1,
			GlobalThreshold: -1,
		},
	}

	mod, _ := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	//	 Test w/Actual database
	mps.MicroSvc.AddModel(model.PageModel{})
	mod, _ = mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}
func TestMicroPageSvcUpdatePage(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + ":")
	mps.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/page")
	mps.MicroSvc.AddModel(mockmodel.MockPageModel{})
	timeStr := time.Now().Format(time.RFC850)
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "UpdatePage",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mps.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method:     "PUT",
			JSONParams: `{"Title":"Updating a Test Page ` + timeStr + `", "Status":{"StatusCode":3,"Desc":"Deactivated"}} `,
			Status:     http.StatusOK,
			MatchVal: dto.Page{
				Id:     4,
				Title:  "Updating a Test Page " + timeStr,
				Status: dto.StatusDetail{dto.DEACTIVATED, fmt.Sprint(dto.StatusType(dto.DEACTIVATED))},
			},
			MatchFields:     []string{"Id", "Title", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}

	//mod, _ := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	//log.Println("----------------------------------------------------------------------------")
	//log.Println("----------------------------------------------------------------------------")
	//log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	//log.Println("----------------------------------------------------------------------------")
	//log.Println("----------------------------------------------------------------------------")
	//resttest.RunTestSet(t, tests)
	// Test w/Actual database
	mps.MicroSvc.AddModel(model.PageModel{})
	mod, _ := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}

func TestMicroPageSvcFindByIdPage(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + ":")
	mps.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/page")
	mps.MicroSvc.AddModel(mockmodel.MockPageModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "FindByIdPage",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mps.MicroSvc.GetFullPath(),
				"%2": "4",
			},
			Method: "GET",
			Status: http.StatusOK,
			MatchVal: dto.Page{
				Id: 4,
			},
			MatchFields:     []string{"Id"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}

	//mod, _ := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	//log.Println("----------------------------------------------------------------------------")
	//log.Println("----------------------------------------------------------------------------")
	//log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	//log.Println("----------------------------------------------------------------------------")
	//log.Println("----------------------------------------------------------------------------")
	//resttest.RunTestSet(t, tests)
	//	 Test w/Actual database
	mps.MicroSvc.AddModel(model.PageModel{})
	mod, _ := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}

func TestMicroPageSvcDeletePage(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + ":")
	mps.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/page")
	//mps.MicroSvc.AddModel(mockmodel.MockPageModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "DeletePage",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mps.MicroSvc.GetFullPath(),
				"%2": "3",
			},
			Method:          "DELETE",
			Status:          http.StatusOK,
			PreAuthId:       7,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}
	//mod, _ := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	//log.Println("----------------------------------------------------------------------------")
	//log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	//log.Println("----------------------------------------------------------------------------")
	//resttest.RunTestSet(t, tests)
	// Test w/Actual database
	mps.MicroSvc.AddModel(model.PageModel{})
	mod, _ := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}
