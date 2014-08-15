/*
MeritWiki Server
*/
package main

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"log"
	"mws/dto"
	//"mws/mockmodel"
	"mws/microrest"
	"mws/model"
	"mws/util"
	"net/http"
	"reflect"
)

func main() {
	log.Printf("------------------------------------------------------------------------------------------------------------------------")
	log.Printf("Starting server  API Version = %s at root %s", util.GetConfig().ApiVersion, util.GetConfig().RootPath)
	log.Printf("------------------------------------------------------------------------------------------------------------------------")
	/* Process command line arguments*/
	util.ReadFlags()
	sysConfig := util.GetConfig()
	//TODO: add Container and collection for microservices set root path and version?
	// Register User Micro Service
	mus := microrest.NewMicroUserSvc()
	mus.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/user")
	mus.MicroSvc.AddModel(model.UserModel{})
	//page
	mps := microrest.NewMicroPageSvc()
	mps.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/page")
	mps.MicroSvc.AddModel(model.PageModel{})
	//section
	mss := microrest.NewMicroSectionSvc()
	mss.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/section")
	mss.MicroSvc.AddModel(model.SectionModel{})
	//revision
	mrs := microrest.NewMicroRevisionSvc()
	mrs.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/revision")
	mrs.MicroSvc.AddModel(model.RevisionModel{})
	mrs.MicroSvc.AddRelationship(reflect.TypeOf(dto.Tag{}), model.TagRevisionModel{})
	mrs.MicroSvc.AddRelationship(reflect.TypeOf(dto.Rating{}), model.RatingModel{})
	//tag
	mts := microrest.NewMicroTagSvc()
	mts.Register(util.GetConfig().RootPath, util.GetConfig().ApiVersion, "/tag")
	mts.MicroSvc.AddModel(model.TagModel{})

	//TODO: need to put url and swagger locations in config
	config := swagger.Config{
		WebServices:    restful.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: sysConfig.Swagger_WebServicesUrl + ":" + sysConfig.Port,
		//WebServicesUrl: "http://localhost:8686",
		ApiPath: "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: sysConfig.SwaggerFilePath}
	//SwaggerFilePath: "/Users/skircher/CodeProjects/go/src/swagger-ui/dist"}

	swagger.InstallSwaggerService(config)

	http.ListenAndServe(":"+sysConfig.Port, nil)
}
