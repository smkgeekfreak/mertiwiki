// MicroUserService.go
package microrest

import (
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	"mws/util"
	"net/http"
	"strconv"
)

type MicroUserSvc struct {
	MicroSvc microService
}

func NewMicroUserSvc() *MicroUserSvc {
	return &MicroUserSvc{
		MicroSvc: microService{},
	}
}

func (mus *MicroUserSvc) authorizeRatingWithOwnership(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	funcName := util.GetCallerName()
	log.Printf("%s Authorizing Request %s", funcName, req.Request.URL.RequestURI())
	model, err := mus.MicroSvc.getModel(mus.MicroSvc.GetFullPath())
	if err != nil || model == nil {
		resp.WriteErrorString(http.StatusPreconditionFailed, " Ownership authoriztion was unavailable")
		log.Printf("%s  - %s Ownership authoriztion was unavailable", funcName, req.SelectedRoutePath())
		return
	}
	// assumes the id for the entity to validate was provided as a path parameter
	entity_id, err := strconv.ParseInt(req.PathParameter("user-id"), 0, 64)
	if err != nil {
		resp.WriteErrorString(http.StatusPreconditionFailed, " Ownership authoriztion entity id was not provided")
		log.Printf("%s  - %s Ownership authoriztion entity id was not provided", funcName, req.SelectedRoutePath())
		return
	}
	if !authorizeRatingWithOwnership(model, entity_id, req, resp) {
		resp.WriteErrorString(http.StatusUnauthorized, " Ownership authoriztion rejected")
		log.Printf("%s  - %s Ownership authoriztion rejected", funcName, req.SelectedRoutePath())
		return
	}
	chain.ProcessFilter(req, resp)
}

func (mus *MicroUserSvc) handleCRUD(request *restful.Request, response *restful.Response) {
	log.Printf("Handling Request %#v", request)
	mus.MicroSvc.handleMicroOp(request, response, &microOperation{})
}

func (mus *MicroUserSvc) Register(rootPath, version, svcPath string) error {
	mus.MicroSvc.rootPath = rootPath
	mus.MicroSvc.version = version
	mus.MicroSvc.svcPath = svcPath

	ws := new(restful.WebService)
	ws.
		Path(mus.MicroSvc.GetFullPath()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("Endpoint path for User api services")

	log.Printf("Registered User API on \r\n      root = %s  \r\n      ver = %s\r\n      svc path = %s \r\n      full path = %s",
		mus.MicroSvc.rootPath, mus.MicroSvc.version, mus.MicroSvc.svcPath, mus.MicroSvc.GetFullPath())

	// All operations for this service require user pre authorization
	//	ws.Filter(requiresClientPreAuthHeader)

	ws.Route(ws.GET("/{user-id}").
		Filter(requiresClientPreAuthHeader).
		Filter(authorizeRating).
		To(mus.handleCRUD).
		// docs
		Doc("Retrive a user based on their unique system identifier").
		Operation("FindById").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("int64")).
		Writes(dto.User{})) // on the response
	//
	ws.Route(ws.POST("/").
		//Filter(authorizeRating).
		To(mus.handleCRUD).
		// docs
		Doc("Add a new user to the system").
		Operation("Create").
		Reads(dto.User{}).
		Writes(dto.User{})) // on the response
	//
	ws.Route(ws.DELETE("/{user-id}").
		Filter(requiresClientPreAuthHeader).
		Filter(mus.authorizeRatingWithOwnership).
		To(mus.handleCRUD).
		// docs
		Doc("Removes a user based on their unique system identifier").
		Operation("Delete").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("int64")).
		Writes(dto.User{})) // on the response
	//
	ws.Route(ws.PUT("/{user-id}").
		Filter(requiresClientPreAuthHeader).
		Filter(mus.authorizeRatingWithOwnership).
		To(mus.handleCRUD).
		// docs
		Doc("Update a user based on their unique system identifier").
		Operation("Update").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("int64")).
		Reads(dto.User{}).
		Writes(dto.User{})) // on the response
	//
	ws.Route(ws.PUT("/{user-id}/status/{status}").
		Filter(requiresClientPreAuthHeader).
		Filter(mus.authorizeRatingWithOwnership).
		To(mus.handleCRUD).
		// docs
		Doc("Update a user based on their unique system identifier").
		Operation("UpdateStatus").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("int64")).
		Param(ws.PathParameter("status", "status to update tag").DataType("dto.StatusType")).
		Writes(dto.User{})) // on the response

	ws.Route(ws.GET("/{user-id}/rating").
		Filter(requiresClientPreAuthHeader).
		Filter(mus.authorizeRatingWithOwnership).
		To(mus.handleCRUD).
		// docs
		Doc("Get a user rating  based on their unique system identifier").
		Operation("FindRating").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("int64")).
		Writes(dto.UserRating{})) // on the response

	//
	for _, rws := range restful.RegisteredWebServices() {
		if rws.RootPath() == mus.MicroSvc.GetFullPath() {
			log.Printf("%s already registered", mus.MicroSvc.GetFullPath())
			return nil
		}
	}
	restful.Add(ws)
	return nil
}
