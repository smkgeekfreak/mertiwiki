// microAuth.go
package microrest

import (
	"github.com/emicklei/go-restful"
	"log"
	"mws/util"
	//"os"
	"fmt"
	"mws/db/pgreflect"
	"mws/model"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func init() {
	// This indicate that client preauthorization header should be used
	// for all requests.  *Could alternately set on for specific micro services
	// or just for individual requests
	//restful.Filter(clientPreAuthUserAuthenticate)
	restful.Filter(NCSACommonLogFormatLogger())
	//restful.Filter(accountRatingAuthenticate)
}

//var logger *log.Logger = log.New(os.Stdout, "", 0)

func NCSACommonLogFormatLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var username = "-"
		if req.Request.URL.User != nil {
			if name := req.Request.URL.User.Username(); name != "" {
				username = name
			}
		}
		chain.ProcessFilter(req, resp)
		log.Printf("%s - %s [%s] \"%s %s %s\" %d %d",
			strings.Split(req.Request.RemoteAddr, ":")[0],
			username,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			req.Request.Method,
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			resp.StatusCode(),
			resp.ContentLength(),
		)
	}
}

func authorizeRating(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	funcName := util.GetCallerName()
	log.Printf("%s Authorizing Request %s", funcName, req.Request.URL.RequestURI())
	// Get client preauth user id
	preAuthId, err := strconv.ParseInt(req.Request.Header.Get("X-meritwiki-client-preauth-id"), 0, 64)
	if err != nil {
		resp.WriteErrorString(403, "403: Rating authoriztion was unavailable")
		log.Printf("%s  - %s 403: Rating authoriztion was unavailable", funcName, req.SelectedRoutePath())
		return
	}
	// Get information to call for postgres function to get user rating
	vdlOperation := "find_user_rating"
	pgfunc, err := pgreflect.GetPgFuncInfo(vdlOperation)
	if err != nil {
		log.Printf("%s -- Error Getting Postgres Function Information- %s ", funcName, vdlOperation)
		return
	}
	// Call postgres function
	var rating int64
	retMap, retCode, err := pgfunc.VariadicScan(preAuthId)
	if err != nil || retMap == nil || retCode != 0 {
		log.Printf("%s -- Error Calling Postgres Function - %s ( %#v)", funcName, pgfunc.Name, err)
		resp.WriteErrorString(403, "403: Rating authoriztion was unavailable")
		return
	} else if len(retMap) <= 0 {
		log.Printf("%s No Rating for user  [%d]", funcName, preAuthId)
		rating = 0
	} else {
		rating = retMap[0]["ret_user_rating"].(int64)
		log.Printf("%s Rating for user (%d) = %d", funcName,
			retMap[0]["ret_user_uid"].(int64),
			retMap[0]["ret_user_rating"].(int64))
	}
	globalThreshold, err := strconv.ParseInt(req.Request.Header.Get("X-meritwiki-global-threshold"), 0, 64)
	if err != nil {
		globalThreshold = int64(0)
	}
	// Check rating against
	if rating < globalThreshold {
		log.Printf("%s -- Insufficent rating for - %s", funcName, req.Request.URL.RequestURI())
		resp.WriteErrorString(403, fmt.Sprintf("403: Insufficient Rating to Perform %s", req.Request.URL.RequestURI()))
		return
	}
	chain.ProcessFilter(req, resp)
}

//
// Helper function for route filters to call to centalize logic for call
// "model" authorizationWithOwnership function
func authorizeRatingWithOwnership(dbModel interface{}, entityId int64, req *restful.Request, resp *restful.Response) bool {
	funcName := util.GetCallerName()
	log.Printf("%s Authorizing Ownership for Request %s", funcName, req.Request.URL.RequestURI())
	// Get client preauth user id
	preAuthId, err := strconv.ParseInt(req.Request.Header.Get("X-meritwiki-client-preauth-id"), 0, 64)
	if err != nil {
		resp.WriteErrorString(http.StatusPreconditionFailed, "Rating authoriztion was unavailable")
		log.Printf("%s  - %s Rating authoriztion was unavailable", funcName, req.SelectedRoutePath())
		return false
	}
	//		get the globalThreshold
	globalThreshold, err := strconv.ParseInt(req.Request.Header.Get("X-meritwiki-global-threshold"), 0, 64)
	if err != nil {
		log.Printf("%s - no global threshold set to Perform %s %s.", funcName, req.Request.Method, req.Request.URL.RequestURI())
		resp.WriteErrorString(http.StatusPreconditionFailed, fmt.Sprintf(" no global threshold set to Perform %s %s", req.Request.Method, req.Request.URL.RequestURI()))
		return false
	}
	// 		get the ownerThreshold
	ownerThreshold, err := strconv.ParseInt(req.Request.Header.Get("X-meritwiki-owner-threshold"), 0, 64)
	if err != nil {
		log.Printf("%s - no threshold set to Perform %s %s.", funcName, req.Request.Method, req.Request.URL.RequestURI())
		resp.WriteErrorString(http.StatusPreconditionFailed, fmt.Sprintf(" no threshold set to perform %s %s", req.Request.Method, req.Request.URL.RequestURI()))
		return false
	}
	// get the vdl operation from the model
	vdlOperation, err := model.GetVDLOperation(dbModel, "authorizeOwnership")
	if err != nil {
		resp.WriteErrorString(http.StatusPreconditionFailed, "Ownership authoriztion was unavailable")
		log.Printf("%s  - %s Ownership  authoriztion was unavailable", funcName, req.SelectedRoutePath())
		return false
	}
	log.Printf("%s Ownership operation for [%v] - [%s]", funcName, reflect.TypeOf(dbModel), vdlOperation)
	// call the function for the operation
	pgfunc, err := pgreflect.GetPgFuncInfo(vdlOperation)
	if err != nil {
		log.Printf("%s -- Error Getting Postgres Function Information- %s ", funcName, vdlOperation)
		return false
	}
	retMap, retCode, err := pgfunc.VariadicScan(preAuthId, entityId)
	//process the results
	var user_rating int64
	var is_owner bool
	if err != nil || retMap == nil || retCode != 0 {
		log.Printf("%s -- Error Calling Postgres Function - %s ( %#v)", funcName, pgfunc.Name, err)
		resp.WriteErrorString(http.StatusPreconditionFailed, "Ownership authoriztion was unavailable")
		return false
	} else if len(retMap) <= 0 {
		log.Printf("%s No Ownership for user  [%d]", funcName, preAuthId)
		user_rating = 0 //??
	} else {
		// 		get user rating
		user_rating = retMap[0]["ret_user_rating"].(int64)
		// 		get ownership (true/false)
		is_owner = retMap[0]["ret_owns"].(bool)
		log.Printf("%s Ownership for user (%d) = %t and has %d", funcName,
			retMap[0]["ret_auth_uid"].(int64),
			is_owner,
			user_rating)
	}
	//if owned enity
	if is_owner {
		// 		compare against the user rating
		if user_rating < ownerThreshold {
			log.Printf("%s - [%d] has insufficient rating [%d] as an owner of [%d] to Perform %s %s, need [%d]",
				funcName, preAuthId, user_rating, entityId, req.Request.Method, req.Request.URL.RequestURI(), ownerThreshold)
			resp.WriteErrorString(http.StatusUnauthorized, fmt.Sprintf("[%d] has insufficient rating [%d] as an owner of [%d] to Perform %s %s, need [%d]",
				preAuthId, user_rating, entityId, req.Request.Method, req.Request.URL.RequestURI(), ownerThreshold))
			return false
		}
		log.Printf("%s -- [%d] as owner of [%d] has rating [%d] meets threshold of [%d]", funcName, preAuthId, entityId, user_rating, ownerThreshold)
		return true
	} else {

		//		compare against the user rating
		if user_rating < globalThreshold {
			log.Printf("%s - [%d] has insufficient rating [%d] as a non-owner of [%d] to Perform %s %s, need [%d]",
				funcName, preAuthId, user_rating, entityId, req.Request.Method, req.Request.URL.RequestURI(), globalThreshold)
			resp.WriteErrorString(http.StatusUnauthorized, fmt.Sprintf("[%d] has insufficient rating [%d] as a non-owner of [%d] to Perform %s %s, need [%d]",
				preAuthId, user_rating, entityId, req.Request.Method, req.Request.URL.RequestURI(), globalThreshold))
			return false
		}

		log.Printf("%s -- [%d] as non-owner of [%d] has rating [%d] meets global threshold of [%d]", funcName, preAuthId, entityId, user_rating, globalThreshold)
		return true
	}
	log.Printf("%s -- Insufficent rating for - %s %s", funcName, req.Request.Method, req.Request.URL.RequestURI())
	return false
}

//
// Use a pre authorized account id set by the client in the request header for
// authentication of the request
func requiresClientPreAuthHeader(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	funcName := util.GetCallerName()
	log.Printf("%s Authenticating Request %#v", funcName, req.Request)
	auth_id := req.Request.Header.Get("X-meritwiki-client-preauth-id")
	log.Printf("%s Authorizing Request %#v", funcName, auth_id)
	if _, err := strconv.ParseInt(auth_id, 0, 64); len(auth_id) == 0 || err != nil {
		// Call valid user code
		resp.AddHeader("WWW-Authenticate", "Client PreAuthenticated Account Realm=Protected Area")
		resp.WriteErrorString(401, "Client PreAuthenticated Account Not Provided")
		log.Printf("%s  - %s  client pre authentication not provide", funcName, req.SelectedRoutePath())
		return
	}
	chain.ProcessFilter(req, resp)
}

//
// Used to indicate that the pre authorized user id needs to be put into the argument list for processing
// down the chain
func usesPreAuthorizedId(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	funcName := util.GetCallerName()
	log.Printf("%s Authorizing Request %s", funcName, req.Request.URL.RequestURI())
	req.Request.Header.Set("X-meritwiki-requires-preauth-id", "true")
	route := req.Request.URL.Path
	log.Printf("%s Authorizing Route Path %#v", funcName, route)
	chain.ProcessFilter(req, resp)
}

// NOT USED
func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	funcName := util.GetCallerName()
	log.Printf("%s Authenticating Request %#v", funcName, req)
	encoded := req.Request.Header.Get("Authorization")
	log.Printf("%s Authorizing Request %#v", funcName, encoded)
	// usr/pwd = admin/admin
	// call acutal username / password auth
	if len(encoded) == 0 || "Basic YWRtaW46YWRtaW4=" != encoded {

		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		log.Printf("%s  - %s failed basic authentication", funcName, req.SelectedRoutePath())
		return
	}
	chain.ProcessFilter(req, resp)
}
