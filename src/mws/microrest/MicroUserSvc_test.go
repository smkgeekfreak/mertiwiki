// MicroUserSvc_test.go
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
	//"os"
	"reflect"
	"testing"
)

var mus *MicroUserSvc

func init() {
	log.Printf("........................ Init Micro User Service Test ...............................")
	//config := util.GetConfig()
	//err := os.Setenv("WIKIENV", "test")
	//if err != nil {
	//	log.Fatalf("Could not set WIKIENV %#v, %v", err, config)
	//}
	mus = NewMicroUserSvc()
}

func TestMicroUserSvcAddUser(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mus.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/user")
	mus.MicroSvc.AddModel(mockmodel.MockUserModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "AddUserHappyCase",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/",
			Sub: map[string]string{
				"%1": mus.MicroSvc.GetFullPath(),
			},
			Method:     "POST",
			JSONParams: `{"Name":"UsernameNew9","Email":"UsernameNew@mail.com","PasswordHash":"23lafjalsdflasfkdjlf","Status":{"StatusCode":0,"Desc":"Initialized"}}`,
			Status:     http.StatusCreated,
			MatchVal: dto.User{
				Name:   "UsernameNew9",
				Email:  "UsernameNew@mail.com",
				Status: dto.StatusDetail{dto.INITIALIZED, fmt.Sprint(dto.StatusType(dto.INITIALIZED))},
			},
			MatchFields:     []string{"Name", "Email", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
		//{
		//	Desc:       "AddUserDuplicate",
		//	Handler:    restful.DefaultContainer.ServeHTTP,
		//	Path:       mus.MicroSvc.GetFullPath(),
		//	Method:     "PUT",
		//	JSONParams: `{"Name":"UsernameNew9","Email":"UsernameNew@mail.com","PasswordHash":"23lafjalsdflasfkdjlf","Status":{"StatusCode":1,"Desc":"Pending"}}`,
		//	Status:     http.StatusCreated, // change after unique name constraint
		//	MatchVal: dto.User{
		//		Name:   "UsernameNew9",
		//		Email:  "UsernameNew@mail.com",
		//		Status: dto.StatusDetail{dto.PENDING, fmt.Sprint(dto.StatusType(dto.PENDING))},
		//	},
		//	MatchFields: []string{"Name", "Email", "Status", "StatusCode", "Desc"},
		//},
		//{
		//	Desc:    "GetUserHappyCase",
		//	Handler: restful.DefaultContainer.ServeHTTP,
		//	Path:    "%1/%2",
		//	Sub: map[string]string{
		//		"%1": mus.MicroSvc.GetFullPath(),
		//		"%2": "1",
		//	},
		//	Method: "GET",
		//	Status: http.StatusOK,
		//	MatchVal: dto.User{
		//		Name:  "Testing Create Should Delete",
		//		Email: "thisisatest@go.com",
		//		//Status: dto.StatusDetail{dto.DELETED, fmt.Sprint(dto.StatusType(dto.DELETED))}, //TODO: update this
		//	},
		//	MatchFields: []string{"Name", "Email"},
		//},
	}
	mod, _ := mus.MicroSvc.getModel(mus.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	// Test w/Actual database
	mus.MicroSvc.AddModel(model.UserModel{})
	mod, _ = mus.MicroSvc.getModel(mus.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}

func TestMicroUserSvcUpdateUser(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mus.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/user")
	mus.MicroSvc.AddModel(mockmodel.MockUserModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "UpdateUser",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mus.MicroSvc.GetFullPath(),
				"%2": "1",
			},
			Method:     "PUT",
			JSONParams: `{"Id":1,"Email":"thisisatest@go.com","Status":{"StatusCode":1,"Desc":"Pending"}}`,
			Status:     http.StatusOK,
			MatchVal: dto.User{
				Id:     1,
				Name:   "Testing Create Should Delete",
				Email:  "thisisatest@go.com",
				Status: dto.StatusDetail{dto.PENDING, fmt.Sprint(dto.StatusType(dto.PENDING))},
			},
			MatchFields:     []string{"Id", "Email", "Status", "StatusCode", "Desc"},
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}
	//mod, _ := mus.MicroSvc.getModel(mus.MicroSvc.GetFullPath())
	//log.Println("----------------------------------------------------------------------------")
	//log.Println("----------------------------------------------------------------------------")
	//log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	//log.Println("----------------------------------------------------------------------------")
	//log.Println("----------------------------------------------------------------------------")
	//resttest.RunTestSet(t, tests)
	// Test w/Actual database
	mus.MicroSvc.AddModel(model.UserModel{})
	mod, _ := mus.MicroSvc.getModel(mus.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)
	log.SetPrefix("")
}

func TestMicroUserSvcDeleteUser(t *testing.T) {
	log.SetPrefix(util.GetCallerName() + " : ")
	mus.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/user")
	mus.MicroSvc.AddModel(mockmodel.MockUserModel{})
	tests := []resttest.RESTTestContainer{
		{
			Desc:    "DeleteUser",
			Handler: restful.DefaultContainer.ServeHTTP,
			Path:    "%1/%2",
			Sub: map[string]string{
				"%1": mus.MicroSvc.GetFullPath(),
				"%2": "3",
			},
			Method:          "DELETE",
			Status:          http.StatusOK,
			PreAuthId:       3,
			OwnThreshold:    0,
			GlobalThreshold: 0,
		},
	}
	// Test w/Actual database
	mus.MicroSvc.AddModel(model.UserModel{})
	mod, _ := mus.MicroSvc.getModel(mus.MicroSvc.GetFullPath())
	log.Println("----------------------------------------------------------------------------")
	log.Printf("----- Testing with Model ( %s )---------------------", reflect.TypeOf(mod))
	log.Println("----------------------------------------------------------------------------")
	resttest.RunTestSet(t, tests)

}
