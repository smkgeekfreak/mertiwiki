// MicroRevisionSvc.go
package microrest

import (
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	"mws/util"
	"net/http"
	"reflect"
	"strconv"
)

type MicroRevisionSvc struct {
	MicroSvc microService
}

func NewMicroRevisionSvc() *MicroRevisionSvc {
	return &MicroRevisionSvc{
		MicroSvc: microService{},
	}
}
func (mrs *MicroRevisionSvc) authorizeRatingWithOwnership(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	funcName := util.GetCallerName()
	log.Printf("%s Authorizing Request %s", funcName, req.Request.URL.RequestURI())
	model, err := mrs.MicroSvc.getModel(mrs.MicroSvc.GetFullPath())
	if err != nil || model == nil {
		resp.WriteErrorString(http.StatusPreconditionFailed, "Ownership authoriztion was unavailable")
		log.Printf("%s  - %s Ownership authoriztion was unavailable", funcName, req.SelectedRoutePath())
		return
	}
	// assumes the id for the entity to validate was provided as a path parameter
	entity_id, err := strconv.ParseInt(req.PathParameter("revision-id"), 0, 64)
	if err != nil {
		resp.WriteErrorString(http.StatusPreconditionFailed, "Ownership authoriztion entity id was not provided")
		log.Printf("%s  - %s Ownership authoriztion entity id was not provided", funcName, req.SelectedRoutePath())
		return
	}
	if !authorizeRatingWithOwnership(model, entity_id, req, resp) {
		resp.WriteErrorString(http.StatusUnauthorized, "Ownership authoriztion rejected")
		log.Printf("%s  - %s Ownership authoriztion rejected", funcName, req.SelectedRoutePath())
		return
	}
	chain.ProcessFilter(req, resp)
}

func (mrs *MicroRevisionSvc) handleCRUD(request *restful.Request, response *restful.Response) {
	funcName := util.GetCallerName()
	log.Printf("%s Handling Request %#v", funcName, request)
	mrs.MicroSvc.handleMicroOp(request, response, &microOperation{})
}

func (mrs *MicroRevisionSvc) handleTagRelationship(request *restful.Request, response *restful.Response) {
	funcName := util.GetCallerName()
	log.Printf("%s Handling Request %#v", funcName, request)
	mOp := microOperation{
		model: mrs.MicroSvc.getRelationship(reflect.TypeOf(dto.Tag{})),
	}
	mrs.MicroSvc.handleMicroOp(request, response, &mOp)
}

func (mrs *MicroRevisionSvc) handleRatingRelationship(request *restful.Request, response *restful.Response) {
	funcName := util.GetCallerName()
	log.Printf("%s Handling Request %#v", funcName, request)
	mOp := microOperation{
		model: mrs.MicroSvc.getRelationship(reflect.TypeOf(dto.Rating{})),
	}
	mrs.MicroSvc.handleMicroOp(request, response, &mOp)
}

func (mrs *MicroRevisionSvc) Register(rootPath, version, svcPath string) error {
	mrs.MicroSvc.rootPath = rootPath
	mrs.MicroSvc.version = version
	mrs.MicroSvc.svcPath = svcPath

	ws := new(restful.WebService)
	ws.
		Path(mrs.MicroSvc.GetFullPath()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("Endpoint path for Revision api services")

	log.Printf("Registered Revision API on \r\n      root = %s  \r\n      ver = %s\r\n      svc path = %s \r\n      full path = %s",
		mrs.MicroSvc.rootPath, mrs.MicroSvc.version, mrs.MicroSvc.svcPath, mrs.MicroSvc.GetFullPath())

	// All operations for this service require user pre authorization
	ws.Filter(requiresClientPreAuthHeader)

	ws.Route(ws.GET("/{revision-id}").
		Filter(authorizeRating).
		To(mrs.handleCRUD).
		// docs
		Doc("Find a revision based on their unique system identifier").
		Operation("FindById").
		Param(ws.PathParameter("revision-id", "identifier of the revision").DataType("int64")).
		Writes(dto.Revision{})) // on the response

	ws.Route(ws.GET("/author/{author-id}").
		Filter(authorizeRating).
		To(mrs.handleCRUD).
		// docs
		Doc("Find a revision based on their unique system identifier").
		Operation("FindByAuthor").
		Param(ws.PathParameter("author-id", "identifier of the revision author to search").DataType("int64")).
		Writes([]dto.Revision{})) // on the response
	//
	ws.Route(ws.POST("/").
		Filter(usesPreAuthorizedId).
		Filter(authorizeRating).
		To(mrs.handleCRUD).
		// docs
		Doc("Add a new revision to the system").
		Operation("Create").
		Reads(dto.Revision{}).
		Writes(dto.Revision{})) // on the response
	//
	ws.Route(ws.DELETE("/{revision-id}").
		Filter(mrs.authorizeRatingWithOwnership).
		To(mrs.handleCRUD).
		// docs
		Doc("Removes a revision based on unique system identifier").
		Operation("Delete").
		Param(ws.PathParameter("revision-id", "identifier of the revision").DataType("int64")).
		Writes(dto.Revision{})) // on the response
	//
	ws.Route(ws.PUT("/{revision-id}").
		Filter(mrs.authorizeRatingWithOwnership).
		To(mrs.handleCRUD).
		// docs
		Doc("Update a revision based on unique system identifier").
		Operation("Update").
		Param(ws.PathParameter("revision-id", "identifier of the revision").DataType("int64")).
		Reads(dto.Revision{}).
		Writes(dto.Revision{})) // on the response
	//
	ws.Route(ws.GET("/{revision-id}/tag/{tag-id}").
		Filter(authorizeRating).
		To(mrs.handleTagRelationship).
		// docs
		Doc("Retrive a revision tag relationship based on their unique system identifiers").
		Operation("RetrieveRelationship").
		Param(ws.PathParameter("revision-id", "identifier of the revision").DataType("int64")).
		Param(ws.PathParameter("tag-id", "identifier of the tag").DataType("int64")).
		Writes(dto.EntityRelationship{})) // on the response
	//
	ws.Route(ws.PUT("/{revision-id}/tag/{tag-id}").
		Filter(mrs.authorizeRatingWithOwnership).
		To(mrs.handleTagRelationship).
		// docs
		Doc("Tag a revsion based on the unique system identifiers for each").
		Operation("CreateRelationship").
		Param(ws.PathParameter("revision-id", "identifier of the revision").DataType("int64")).
		Param(ws.PathParameter("tag-id", "identifier of the tag").DataType("int64")).
		Writes(dto.EntityRelationship{})) // on the response
	//
	ws.Route(ws.DELETE("/{revision-id}/tag/{tag-id}").
		Filter(mrs.authorizeRatingWithOwnership).
		To(mrs.handleTagRelationship).
		// docs
		Doc("Removes a tag from a revision based on the unique system identifiers for each").
		Operation("DeleteRelationship").
		Param(ws.PathParameter("revision-id", "identifier of the revision").DataType("int64")).
		Param(ws.PathParameter("tag-id", "identifier of the tag").DataType("int64")).
		Writes(dto.EntityRelationship{}))
	//
	ws.Route(ws.PUT("/{revision-id}/rate/{rating}").
		Filter(usesPreAuthorizedId).
		Filter(mrs.authorizeRatingWithOwnership).
		To(mrs.handleRatingRelationship).
		// docs
		Doc("Rate a revision based on the unique system identifier ").
		Operation("CreateRelationship").
		Param(ws.PathParameter("revision-id", "identifier of the revision").DataType("int64")).
		Param(ws.PathParameter("rating", "identifier of the revision").DataType("int64")).
		Writes(dto.Rating{}))
	//
	for _, rws := range restful.RegisteredWebServices() {
		if rws.RootPath() == mrs.MicroSvc.GetFullPath() {
			log.Printf("%s already registered, not updating", mrs.MicroSvc.GetFullPath())
			return nil
		}
	}
	restful.Add(ws)
	return nil
}
