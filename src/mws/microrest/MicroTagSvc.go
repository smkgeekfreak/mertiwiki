// MicroTagService.go
package microrest

import (
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	//"mws/util"
	//"net/http"
	//"strconv"
)

type MicroTagSvc struct {
	MicroSvc microService
}

func NewMicroTagSvc() *MicroTagSvc {
	return &MicroTagSvc{
		MicroSvc: microService{},
	}
}

func (mps *MicroTagSvc) handleCRUD(request *restful.Request, response *restful.Response) {
	log.Printf("Handling Request %#v", request)
	mps.MicroSvc.handleMicroOp(request, response, &microOperation{})
}

func (mps *MicroTagSvc) Register(rootPath, version, svcPath string) error {
	mps.MicroSvc.rootPath = rootPath
	mps.MicroSvc.version = version
	mps.MicroSvc.svcPath = svcPath

	ws := new(restful.WebService)
	ws.
		Path(mps.MicroSvc.GetFullPath()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("Endpoint path for Tag api services")

	log.Printf("Registered Tag API on \r\n      root = %s  \r\n      ver = %s\r\n      svc path = %s \r\n      full path = %s",
		mps.MicroSvc.rootPath, mps.MicroSvc.version, mps.MicroSvc.svcPath, mps.MicroSvc.GetFullPath())

	// All operations for this service require user pre authorization
	ws.Filter(requiresClientPreAuthHeader)

	ws.Route(ws.GET("/{tag-id}").
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Retrive a tag based on their unique system identifier").
		Operation("FindById").
		Param(ws.PathParameter("tag-id", "identifier of the tag").DataType("int64")).
		Writes(dto.Tag{})) // on the response

	ws.Route(ws.POST("/").
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Add a new tag to the system").
		Operation("Create").
		Reads(dto.Tag{}).
		Writes(dto.Tag{})) // on the response
	//
	ws.Route(ws.DELETE("/{tag-id}").
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Removes a tag based on their unique system identifier").
		Operation("Delete").
		Param(ws.PathParameter("tag-id", "identifier of the tag").DataType("int64")).
		Writes(dto.Tag{})) // on the response
	//
	ws.Route(ws.PUT("/{tag-id}").
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Update a tag based on their unique system identifier").
		Operation("Update").
		Param(ws.PathParameter("tag-id", "identifier of the tag").DataType("int64")).
		Reads(dto.Tag{}).
		Writes(dto.Tag{})) // on the response
	//
	ws.Route(ws.PUT("/{tag-id}/status/{status}").
		Filter(authorizeRating).
		To(mps.handleCRUD).
		// docs
		Doc("Update a tag based on their unique system identifier").
		Operation("Update").
		Param(ws.PathParameter("tag-id", "identifier of the tag").DataType("int64")).
		Param(ws.PathParameter("status", "status to update tag").DataType("dto.StatusType")).
		Writes(dto.Tag{})) // on the response

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
