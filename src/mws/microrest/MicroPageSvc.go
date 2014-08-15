// MicroPageService.go
package microrest

import (
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	"mws/util"
	"net/http"
	"strconv"
)

type MicroPageSvc struct {
	MicroSvc microService
}

func NewMicroPageSvc() *MicroPageSvc {
	return &MicroPageSvc{
		MicroSvc: microService{},
	}
}

func (mps *MicroPageSvc) authorizeRatingWithOwnership(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	funcName := util.GetCallerName()
	log.Printf("%s Authorizing Request %s", funcName, req.Request.URL.RequestURI())
	model, err := mps.MicroSvc.getModel(mps.MicroSvc.GetFullPath())
	if err != nil || model == nil {
		resp.WriteErrorString(http.StatusPreconditionFailed, " Ownership authoriztion was unavailable")
		log.Printf("%s  - %s Ownership authoriztion was unavailable", funcName, req.SelectedRoutePath())
		return
	}
	// assumes the id for the entity to validate was provided as a path parameter
	entity_id, err := strconv.ParseInt(req.PathParameter("page-id"), 0, 64)
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

func (mps *MicroPageSvc) handleCRUD(request *restful.Request, response *restful.Response) {
	log.Printf("Handling Request %#v", request)
	mps.MicroSvc.handleMicroOp(request, response, &microOperation{})
}

func (mps *MicroPageSvc) Register(rootPath, version, svcPath string) error {
	mps.MicroSvc.rootPath = rootPath
	mps.MicroSvc.version = version
	mps.MicroSvc.svcPath = svcPath

	ws := new(restful.WebService)
	ws.
		Path(mps.MicroSvc.GetFullPath()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("Endpoint path for Page api services")

	log.Printf("Registered Page API on \r\n      root = %s  \r\n      ver = %s\r\n      svc path = %s \r\n      full path = %s",
		mps.MicroSvc.rootPath, mps.MicroSvc.version, mps.MicroSvc.svcPath, mps.MicroSvc.GetFullPath())

	// All operations for this service require user pre authorization
	ws.Filter(requiresClientPreAuthHeader)

	ws.Route(ws.GET("/{page-id}").
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Retrive a page based on their unique system identifier").
		Operation("FindById").
		Param(ws.PathParameter("page-id", "identifier of the page").DataType("int64")).
		Writes(dto.Page{})) // on the response

	ws.Route(ws.GET("/author/{author-id}").
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Retrive a page based on their unique system identifier").
		Operation("FindByAuthor").
		Param(ws.PathParameter("author-id", "identifier of the author or the page").DataType("int64")).
		Writes(dto.Page{})) // on the response
	//
	ws.Route(ws.POST("/").
		Filter(usesPreAuthorizedId).
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Add a new page to the system").
		Operation("Create").
		Reads(dto.Page{}).
		Writes(dto.Page{})) // on the response
	//
	ws.Route(ws.DELETE("/{page-id}").
		Filter(mps.authorizeRatingWithOwnership).
		To(mps.handleCRUD).
		// docs
		Doc("Removes a page based on their unique system identifier").
		Operation("Delete").
		Param(ws.PathParameter("page-id", "identifier of the page").DataType("int64")).
		Writes(dto.Page{})) // on the response
	//
	ws.Route(ws.PUT("/{page-id}").
		Filter(mps.authorizeRatingWithOwnership).
		To(mps.handleCRUD).
		// docs
		Doc("Update a page based on their unique system identifier").
		Operation("Update").
		Param(ws.PathParameter("page-id", "identifier of the page").DataType("int64")).
		Reads(dto.Page{}).
		Writes(dto.Page{})) // on the response

	//
	for _, rws := range restful.RegisteredWebServices() {
		if rws.RootPath() == mps.MicroSvc.GetFullPath() {
			log.Printf("%s already registered", mps.MicroSvc.GetFullPath())
			return nil
		}
	}
	restful.Add(ws)
	return nil
}
