// microService.go
package microrest

import (
	"bytes"
	"fmt"
	"io"
	//"net/http"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/emicklei/go-restful"
	"log"
	"mws/dto"
	"mws/util"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var typeMap map[string]reflect.Type

func init() {
	typeMap = make(map[string]reflect.Type)
	typeMap["int64"] = reflect.TypeOf(int64(1))
	typeMap["string"] = reflect.TypeOf(string(""))
	typeMap["dto.StatusType"] = reflect.TypeOf(dto.StatusType(0))
}

func isAllowedParamType(typeName string) bool {
	for key, _ := range typeMap {
		if key == strings.ToLower(typeName) {
			return true
		}
	}
	return false
}

func getType(typeStr string) (reflect.Type, error) {
	//funcName := util.GetCallerName()
	log.Printf("Get Type from %#v", typeMap)
	if typeMap != nil {
		return typeMap[typeStr], nil
	}
	return nil, fmt.Errorf("Model not found for %s", typeStr)
}

func (ms *microService) AddModel(model interface{}) error {
	funcName := util.GetCallerName()
	key := ms.GetFullPath()
	if ms.models == nil {
		log.Printf("%s mS Model initialized for %s", funcName, key)
		ms.models = make(map[string]interface{})
	}
	for k, v := range ms.models {
		if k == key {
			log.Printf("ms Model exists for %s with model %v", k, v)
			log.Printf("ms Model added for %s", key)
			ms.models[key] = model
			log.Printf("mS Models %#v", ms.models)
			return fmt.Errorf("Overrode model for key %s with %s", key, model)
		}
	}
	log.Printf("mS Model added for %s", key)
	ms.models[key] = model
	log.Printf("mS Models %#v", ms.models)
	return nil
}

func (ms *microService) getModel(key string) (interface{}, error) {
	//log.SetPrefix(util.GetCallerName() + ":")
	if ms.models != nil {
		log.Printf("ms Get Model for %s", key)
		return ms.models[key], nil
	}
	return nil, fmt.Errorf("Model not found for %s", key)
}

func (ms *microService) AddRelationship(key reflect.Type, relation interface{}) error {
	//log.SetPrefix(util.GetCallerName() + ":")
	//key := ms.GetFullPath()
	if ms.relationships == nil {
		log.Printf("mS Relationship initialized for %s", key)
		ms.relationships = make(map[reflect.Type]interface{})
	}
	for k, v := range ms.relationships {
		if k == key {
			log.Printf("ms Relationship exists for %v with model %v", k, v)
			log.Printf("ms Relationship added for %v", key)
			ms.relationships[key] = relation
			log.Printf("mS Relationship %#v", ms.relationships)
			return fmt.Errorf("Overrode model for key %v with %s", key, relation)
		}
	}
	log.Printf("mS Relationship added for %s", key)
	ms.relationships[key] = relation
	log.Printf("mS Relationship %s", ms.relationships)
	return nil
}

func (ms *microService) getRelationship(key reflect.Type) interface{} {
	//log.SetPrefix(util.GetCallerName() + ":")
	if ms.relationships != nil {
		log.Printf("ms Get Relationship  for %s", key)
		return ms.relationships[key]
	}
	log.Printf("Model not found for relationship %s", key)
	return nil
}

type microOperation struct {
	model            interface{}   // Which data model to operate on
	operation        string        //Operation call
	args             []interface{} //Path Parameters
	readParam        interface{}   //Body Parameter
	writeType        reflect.Type  //Output Type
	vdlOperationName string        //Name of the associate VDL operation
	authIdRequired   bool          // identifies whether or not to pass the authorization id in the argument list
	authId           int64         // The authorized account id for this operation ( may removed this)
}

type microService struct {
	rootPath      string
	version       string
	svcPath       string
	models        map[string]interface{}
	relationships map[reflect.Type]interface{}
}

func (ms *microService) GetVersion() string {
	return ms.version
}

func (ms *microService) GetRootPath() string {
	return ms.rootPath
}

func (ms *microService) GetVersionedRootPath() string {
	var verPath string
	if len(ms.rootPath) > 0 {
		verPath = ms.rootPath
		//log.Printf("Build UserService path %s", verPath)
	}
	if len(ms.version) > 0 {
		if !strings.Contains(ms.version, "/") {
			verPath = verPath + "/"
		}
		verPath = verPath + ms.version
		//log.Printf("Build UserService path %s", verPath)
	}
	return verPath
}

func (ms *microService) GetSvcPath() string {
	return ms.svcPath
}

func (ms *microService) GetFullPath() string {
	fullpath := ms.GetVersionedRootPath()
	if len(ms.svcPath) > 0 {
		if !strings.Contains(ms.svcPath, "/") {
			fullpath = fullpath + "/"
		}
		fullpath = fullpath + ms.svcPath

		//log.Printf("Build UserService path %s", fullpath)
	}
	return fullpath
}

func (ms *microService) CREATE_(op microOperation) (interface{}, int, error) {
	//log.SetPrefix(util.GetCallerName() + ":")
	retCode := http.StatusConflict
	// Get parameters from the request
	//readObj := reflect.ValueOf(op.readParam).Interface()
	log.Printf("CRUD Operation - %s - with readObj (%v)", op.operation, op.args)
	// Get the appropriate data model to use to act on the parameters
	model := op.model
	log.Printf("Model type = %#v", model)
	// Call operation on the data model
	rets, err := invokeModel(model, op.operation, op.args)
	if err != nil {
		return nil, retCode, err
	}
	// Process Return data
	obj := rets[0].Interface()
	code := rets[1].Interface()
	retErr := rets[2].Interface()
	log.Printf("Creating Returned = %#v, %#v", obj, retErr)
	if retErr != nil || obj == nil {
		log.Printf("Error creating due to ' %s ' ", retErr)
		retCode = http.StatusExpectationFailed
		return nil, retCode, retErr.(error)
	}
	//TODO: Fix this convention in the DataModel for CREATE and UPDATE
	// Move to a 0, 1, -1 convention
	// 0 = created
	// 1 = exists
	// anything else is erroneous
	switch code {
	case 0:
		retCode = http.StatusCreated
	case 1:
		retCode = http.StatusOK
	default:
		retCode = http.StatusNotModified
	}
	return obj, retCode, nil
}

func (ms *microService) RETRIEVEBYID_(op microOperation) (interface{}, int, error) {
	retCode := http.StatusConflict
	//log.SetPrefix(util.GetCallerName() + ":")
	//TODO: does it make sense to try and handle more than one arg here?
	//for i, _ := range args {
	//log.Printf("Inside RETRIEVE By=%v, ", reflect.ValueOf(op.args[0]))
	id := reflect.ValueOf(op.args[0]).Interface()
	log.Printf("CRUD Operation - %s (%v)", op.operation, id)
	//}
	// Get the appropriate data model to use to act on the parameters
	model := op.model
	log.Printf("Model type = %#v", model)
	// Call operation on the data model
	log.Printf("Invoking Model type =( %v, %s, %s", model, op.operation, id)
	rets, err := invokeModel(model, op.operation, id)
	if err != nil {
		return nil, retCode, err
	}
	// Process Return data
	ret := rets[0].Interface()
	retErr := rets[1].Interface()
	log.Printf("Retrieve Returned = %#v, %#v", ret, retErr)

	if retErr != nil || &ret == nil {
		err1 := retErr.(error)
		log.Printf("Error finding due to ' %s ' ", err1)
		return ret, http.StatusNotFound, err1
	}
	return ret, http.StatusOK, nil
}

func (ms *microService) UPDATE_(op microOperation) (interface{}, int, error) {
	//log.SetPrefix(util.GetCallerName() + ":")
	retCode := http.StatusConflict
	// Get parameters from the request
	//readObj := reflect.ValueOf(op.readParam).Interface()
	//log.Printf("CRUD Operation - %s (%v)", op.operation, readObj)
	// Get the appropriate data model to use to act on the parameters
	model := op.model
	log.Printf("Model type = %#v", model)
	// Call operation on the data model
	rets, err := invokeModel(model, op.operation, op.args)
	if err != nil {
		return nil, retCode, err
	}
	// Process Return data
	obj := rets[0].Interface()
	code := rets[1].Interface()
	retErr := rets[2].Interface()
	log.Printf("Updating Returned = %#v, %#v", obj, retErr)
	if retErr != nil || obj == nil {
		log.Printf("Error updating  due to ' %s ' ", retErr)
		retCode = http.StatusExpectationFailed
		return nil, retCode, retErr.(error)
	}
	//TODO: Fix this convention in the DataModel for CREATE and UPDATE
	// Move to a 0, 1, -1 convention
	// 0 = created
	// 1 = exists
	// anything else is erroneous
	switch code {
	case 0:
		retCode = http.StatusOK
	default:
		retCode = http.StatusNotModified
	}
	return obj, retCode, nil
}

func (ms *microService) DELETE_(op microOperation) (interface{}, int, error) {
	//log.SetPrefix(util.GetCallerName() + ":")
	retCode := http.StatusConflict
	// Get parameters from the request
	id := reflect.ValueOf(op.args[0]).Interface()
	log.Printf("CRUD Operation - %s (%v)", op.operation, id)
	// Get the appropriate data model to use to act on the parameters
	model := op.model
	log.Printf("Model type = %#v", model)
	// Call operation on the data model
	rets, err := invokeModel(model, op.operation, id)
	if err != nil {
		return nil, retCode, err
	}
	// Process Return data
	obj := rets[0].Interface()
	code := rets[1].Interface()
	retErr := rets[2].Interface()
	log.Printf("Deleting Returned = %v, %v %v", obj, code, retErr)
	if retErr != nil || obj == nil {
		log.Printf("Error deleting due to ' %s ' ", retErr)
		retCode = http.StatusExpectationFailed
		return nil, retCode, retErr.(error)
	}
	//TODO: Fix this convention in the DataModel for CREATE and UPDATE
	// Move to a 0, 1, -1 convention
	// 0 = created
	// 1 = exists
	// anything else is erroneous
	switch code {
	case 0:
		retCode = http.StatusOK
	case 1:
		retCode = http.StatusNotModified
		// default case retCode remain StatusConflict
	}
	return obj, retCode, nil
}

func (ms *microService) UPDATESTATUS_(op microOperation) (interface{}, int, error) {
	//log.SetPrefix(util.GetCallerName() + ":")
	retCode := http.StatusConflict
	// Get parameters from the request
	id := reflect.ValueOf(op.args[0]).Interface()
	status := reflect.ValueOf(op.args[1]).Interface()
	log.Printf("CRUD Operation - %s (%d, %d)", op.operation, id, status)
	// Get the appropriate data model to use to act on the parameters
	model := op.model
	log.Printf("Model type = %v", model)
	// Call operation on the data model
	//insert VDL op name
	rets, err := invokeModel(model, op.operation, id, status)
	if err != nil {
		return nil, retCode, err
	}
	// Process Return data
	obj := rets[0].Interface()
	code := rets[1].Interface()
	retErr := rets[2].Interface()
	log.Printf("Updating Returned = %#v, %#v", obj, retErr)
	if retErr != nil || obj == nil {
		log.Printf("Error updating  due to ' %s ' ", retErr)
		retCode = http.StatusExpectationFailed
		return nil, retCode, retErr.(error)
	}
	//TODO: Fix this convention in the DataModel for CREATE and UPDATE
	// Move to a 0, 1, -1 convention
	// 0 = created
	// 1 = exists
	// anything else is erroneous
	switch code {
	case 0:
		retCode = http.StatusOK
	default:
		retCode = http.StatusNotModified
	}
	return obj, retCode, nil
}

func (ms *microService) CREATERELATIONSHIP_(op microOperation) (interface{}, int, error) {
	//log.SetPrefix(util.GetCallerName() + ":")
	retCode := http.StatusConflict
	// Get parameters from the request
	//rel_id_1 := reflect.ValueOf(op.args[0]).Interface()
	//rel_id_2 := reflect.ValueOf(op.args[1]).Interface()

	log.Printf(" Called %s (%#v)", op.operation, op.args)
	// Get the appropriate data model to use to act on the parameters
	model := op.model
	log.Printf(" Model type = %#v", model)
	// Call operation on the data model
	rets, err := invokeModel(model, op.operation, op.args)
	if err != nil {
		return nil, retCode, err
	}
	// Process Return data
	obj := rets[0].Interface()
	code := rets[1].Interface()
	retErr := rets[2].Interface()
	log.Printf("%s returned = %v, %v, %v", op.operation, obj, code, retErr)
	if retErr != nil || obj == nil {
		log.Printf("Error in %s due to ' %s ' ", op.operation, retErr)
		retCode = http.StatusExpectationFailed
		return nil, retCode, retErr.(error)
	}
	//TODO: Fix this convention in the DataModel for CREATE and UPDATE
	// Move to a 0, 1, -1 convention
	// 0 = created
	// 1 = exists
	// anything else is erroneous
	switch code {
	case 0:
		retCode = http.StatusOK
	case 1:
		retCode = http.StatusNotModified
		// default case retCode remain StatusConflict
	}
	return obj, retCode, nil
}

//
// Support generic proxy to model operations assuming back conventions
// assuming the model implements an exported method for the specified operation
// with arguments of the type specified in the associated MicroService API definition
// in the order of path parameters (in order) followed by the Reads() value (if any) and
// returning some value object, integer return code and error
// If these conventions do not meet the needs, then a custom proxy method can be written, i.e.
// and are named with the uppercase of the opeartion + "_", i.e.
// CREATE_, UPDATE_, CREATERELATIONSHIP_
func (ms *microService) genericModelProxy(op microOperation) (interface{}, int, error) {
	funcName := util.GetCallerName()
	//log.SetPrefix(util.GetCallerName() + ":")
	retCode := http.StatusConflict
	// Get parameters from the request
	log.Printf("%s Called %s (%#v)", funcName, op.operation, op.args)
	// Get the appropriate data model to use to act on the parameters
	model := op.model
	log.Printf(" Model type = %#v", model)
	// Call operation on the data model
	rets, err := invokeModel(model, op.operation, op.args)
	if err != nil {
		return nil, retCode, err
	}
	// Process Return data
	obj := rets[0].Interface()
	code := rets[1].Interface()
	retErr := rets[2].Interface()
	log.Printf("%s returned = %v, %v, %v", op.operation, obj, code, retErr)
	if retErr != nil || obj == nil {
		log.Printf("Error in %s due to ' %s ' ", op.operation, retErr)
		retCode = http.StatusExpectationFailed
		return nil, retCode, retErr.(error)
	}
	//TODO: Fix this convention in the DataModel for CREATE and UPDATE
	// Move to a 0, 1, -1 convention
	// 0 = created
	// 1 = exists
	// anything else is erroneous
	switch code {
	case 0:
		retCode = http.StatusOK
	case 1:
		retCode = http.StatusNotModified
		// default case retCode remain StatusConflict
	}
	return obj, retCode, nil
}

func (ms *microService) handleMicroOp(request *restful.Request, response *restful.Response, microOpRef *microOperation) {
	funcName := util.GetCallerName()
	/* Inspect request*/

	if request == nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Request corrupted")
		return
	}
	microOp := *microOpRef
	activeRoute, activeWS := findActiveRoute(request)

	/* No matching route was found, error return */
	if activeRoute == nil {
		log.Printf("%s No matching route found for %s", funcName, request.SelectedRoutePath())
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("No matching route found for %s", request.SelectedRoutePath()))
		return
	}
	//
	// If the model for the operation is not set, atttempt to get the model using the root path for this active restful.WebService
	if microOp.model == nil {
		//Get model assocaited with the root path for this service
		log.Printf("Models = %#v", ms.models)
		mod, err := ms.getModel(activeWS.RootPath())
		if err != nil || mod == nil {
			log.Printf("No Matching model for " + activeWS.RootPath())
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, "No Matching model for "+activeWS.RootPath())
			return
		}
		log.Printf("Model = %#v", mod)
		microOp.model = mod
		crudMethodName := strings.ToLower(strings.Replace(((*activeRoute).Operation), " ", "", -1))
		field, ok := reflect.TypeOf(mod).FieldByName(crudMethodName)
		if ok == true {
			log.Printf("VDL op for %s = %s", crudMethodName, field.Tag.Get("vdl"))
		}
	}

	log.Printf("Active Route %#v", *activeRoute)
	log.Printf("Consumes %#v", (*activeRoute).Consumes)
	microOp.operation = activeRoute.Operation
	//

	//Construct arguments from  parameters
	microOp.args = make([]interface{}, 0)
	arrArgs := make([]interface{}, 0)
	startIndex := 0
	// If "X-meritwiki-requires-preauth-id"  header exists it indicates this operation needs the
	// authorized account id, and it becomes the first parameter in the argument list
	preAuthHeader := request.Request.Header.Get("X-meritwiki-requires-preauth-id")
	if len(preAuthHeader) > 0 {
		needsPreAuth, err := strconv.ParseBool(preAuthHeader)
		if err != nil {
			log.Printf("Problem determining if preauth id is necessary")
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, "Problem determining if preauth id is necessary")
			return
		}
		if needsPreAuth {
			preAuthId, err := strconv.ParseInt(request.Request.Header.Get("X-meritwiki-client-preauth-id"), 0, 64)
			if err != nil {
				response.WriteErrorString(403, "403: Authorization Id was not parsable")
				log.Printf("%s  - %s 403: Authorization Id was not parsable", funcName, request.SelectedRoutePath())
				return
			}
			log.Printf("%s microOperation authRequired setting authId (%d) in parameter list", funcName, preAuthId)
			microOp.authIdRequired = needsPreAuth
			microOp.authId = preAuthId
			// Set the first arg to the pre auth id
			arrArgs = append(arrArgs, preAuthId)
			log.Printf("%s operation has %d args ", funcName, len(arrArgs))
		}
	}
	log.Printf("%s Preauth id determined not necessary", funcName)
	//
	// process path parameters
	for key, val := range request.PathParameters() {
		log.Printf("Key %s, Val %s, Params %#v", key, val, activeRoute.ParameterDocs[startIndex].Data().DataType)
		tp, _ := getType(activeRoute.ParameterDocs[startIndex].Data().DataType)
		if v, err := makeArg(val, tp, request.HeaderParameter("Content-Type")); err == nil {
			log.Printf("Make Args gave %#v", v)
			arrArgs = append(arrArgs, v.Interface())
		} else {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Problem with path parameters %s", err))
			return //nil, state
		}
		startIndex++
	}
	// process body parameters
	// Construct argument from ReadSample, if present
	if activeRoute.ReadSample != nil {
		log.Printf("Checking Read Sample %s", reflect.TypeOf(activeRoute.ReadSample))
		readSam := reflect.New(reflect.TypeOf(activeRoute.ReadSample)).Interface()
		err := request.ReadEntity(&readSam)
		if err != nil {
			log.Printf("Problem with Read Sample %s", err)
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Problem with Read Sample %s", err))
			return
		}
		microOp.readParam = readSam
		arrArgs = append(arrArgs, reflect.ValueOf(readSam).Interface())
		log.Printf("Arg from ReadSample %v", readSam)
		log.Printf("%s operation has %d args ", funcName, len(arrArgs))
	} else {
		log.Printf("Checking Body Parameter")
		for _, val := range activeRoute.ParameterDocs {
			log.Printf("Kind %#v, Name %s, Type %s", val.Data().Kind, val.Data().Name, val.Data().DataType)
			if val.Data().Kind == restful.BodyParameterKind {
				bodyParam, err := request.BodyParameter(val.Data().Name)

				buf := new(bytes.Buffer)
				io.Copy(buf, request.Request.Body)
				body := buf.String()

				log.Printf("Body param = %s", body)
				log.Printf("Body param = %s", bodyParam)
				if err != nil {
					log.Printf("Problem with Body Param %s", err)
					response.AddHeader("Content-Type", "text/plain")
					response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Problem with Body Param %s", err))
					return
				}
				tp, _ := getType(val.Data().DataType)
				if v, err := makeArg(body, tp, request.HeaderParameter("Content-Type")); err == nil {
					log.Printf("Make Args gave %#v", v)
					arrArgs = append(arrArgs, v.Interface())
				} else {
					response.AddHeader("Content-Type", "text/plain")
					response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Problem with path parameters %s", err))
					return //nil, state
				}
			}
		}
	}

	microOp.args = arrArgs

	if activeRoute.WriteSample != nil {
		microOp.writeType = reflect.TypeOf(activeRoute.WriteSample)
	}
	log.Printf("MicroOp  %#v", microOp)

	inputs := make([]reflect.Value, 1)
	inputs[0] = reflect.ValueOf(microOp)
	//
	//Find matching method by 'Operation' tag on the active Route
	crudMethodName := strings.ToUpper(strings.Replace(((*activeRoute).Operation), " ", "", -1))
	crudMethodName = crudMethodName + "_"
	log.Printf("Translasted Method Name %#v", crudMethodName)

	// Invoke reflective CRUD Method
	var writeObj, retCode, retErr interface{}
	retVals, err := invokeReflectiveMethod(ms, crudMethodName, inputs)
	if err != nil {
		//
		// If the reflective method wasn't found, try the generic model proxy
		obj, rc, er := ms.genericModelProxy(microOp)
		log.Printf("%s Generic Model returned %v, %v, %v", funcName, obj, rc, er)
		if rc == http.StatusNotModified {
			log.Printf("No results returned from %s", microOp.operation)
			//response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusNoContent, fmt.Sprintf("No results returned"))
			return
		}
		writeObj = obj
		retCode = rc
		retErr = er
	} else {
		log.Printf("%s Model returned %v, %v, %v", funcName, retVals[0], retVals[1], retVals[2])
		writeObj = retVals[0].Interface()
		retCode = retVals[1].Interface()
		retErr = retVals[2].Interface()
	}
	//
	// Process return arguments
	if retErr != nil || writeObj == nil || retCode.(int) < 0 {
		log.Printf("Operation %s failed due to %s", microOp.operation, retErr)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(retCode.(int), fmt.Sprintf("Operation failed due to %s", retErr))
		return
	}
	if !reflect.TypeOf(writeObj).AssignableTo(microOp.writeType) && !reflect.TypeOf(writeObj).AssignableTo(reflect.PtrTo(microOp.writeType)) {
		log.Printf("Documented API return type does not match: %s not assignable %s", reflect.TypeOf(writeObj), microOp.writeType)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusExpectationFailed, fmt.Sprintf("Documented API return type does not match: %s not assignable %s", reflect.TypeOf(writeObj), microOp.writeType))
		return
	}
	log.Printf("MicroOp Returned  \r\n Write Object: \r\n %v,  \r\n  ReturnCode(%d), \r\n Error(%v)", writeObj, retCode, retErr)

	response.WriteHeader(retCode.(int))
	response.WriteEntity(writeObj)
	return
}

func findActiveRoute(request *restful.Request) (*restful.Route, *restful.WebService) {
	var activeRoute *restful.Route
	var activeWS *restful.WebService
	/* Search registered web services for this route */
	for _, ws := range restful.DefaultContainer.RegisteredWebServices() {
		for _, route := range ws.Routes() {
			//log.Printf("Look path %s, found path %s", request.SelectedRoutePath(), route.Path)
			// Compare Route based on path and method
			if route.Path == request.SelectedRoutePath() && route.Method == request.Request.Method {
				log.Printf("Matched path %s, found path %s", request.SelectedRoutePath(), route.Path)
				activeRoute = &route
				activeWS = ws
				//Break out of loop, since route has been located
				return activeRoute, activeWS
			}
		}
	}
	return nil, nil
}

func invokeReflectiveMethod(any interface{}, methodName string, args []reflect.Value) (result []reflect.Value, failed error) {
	if method := reflect.ValueOf(any).MethodByName(methodName); method.IsValid() {
		retval := method.Call(args)
		//log.Printf("return from operation invoke %#v", retval[0].Interface())
		return retval, nil
	} else {
		log.Printf("Method %s not found on %s", methodName, reflect.TypeOf(any))
		return nil, errors.New(fmt.Sprintf("Method %s not found on %s", methodName, reflect.TypeOf(any)))
	}
}

func invokeModel(model interface{}, operationName string, args ...interface{}) (result []reflect.Value, failed error) { //interface{} {
	funcName := util.GetCallerName()
	log.Printf("%s model  [%v] [%s] [%#v]", funcName, reflect.TypeOf(model), operationName, args)
	if model == nil {
		log.Printf("%sModel not set", funcName)
		return nil, nil
	}
	inputs := make([]reflect.Value, 0)
	//
	// Use operationName to look up the VDL operation of the model
	// This assumes the model has fields for "operationName" and is annoated with
	// the VDL operations on those fields
	modelType := reflect.TypeOf(model)
	operation := strings.ToLower(operationName) // all postgres operations are lowercase
	vdlField, hasVdl := modelType.FieldByName(operation)
	if hasVdl { // == true
		//
		//If vdl operation provided then prepend it to the argument list
		vdlOperation := vdlField.Tag.Get("vdl")
		log.Printf("%s VDL Opeartion Field [%s]  found [%s]", funcName, operation, vdlOperation)
		inputs = append(inputs, reflect.ValueOf(vdlOperation))
	} else {
		log.Printf("VDL Opeartion Field (%s)) not found", operation)
	}
	//
	// process argument list
	log.Printf("processing %d args ", len(args))
	for i, _ := range args {
		switch reflect.TypeOf(args[i]).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(args[i])
			log.Printf("processing arg with %d args ", s.Len())
			for j := 0; j < s.Len(); j++ {
				a := s.Index(j).Interface()
				log.Printf("processed arg  %v  ", a)
				inputs = append(inputs, reflect.ValueOf(a))
			}
		default:
			inputs = append(inputs, reflect.ValueOf(args[i]))
		}

	}
	log.Printf(" %s calling [%s] on [%v] with %d args %v", funcName, operationName, modelType, len(inputs), inputs)

	if method := reflect.ValueOf(model).MethodByName(operationName); method.IsValid() {
		retval := method.Call(inputs)
		//retval := reflect.ValueOf(model).MethodByName(operationName).Call(inputs)
		result = retval
		return result, nil
	} else {
		log.Printf("Opearation %s not found on %s", operationName, modelType)
		return nil, errors.New(fmt.Sprintf("Method %s not found on %s", operationName, reflect.TypeOf(model)))
	}

}

func makeArg(data string, template reflect.Type, mime string) (reflect.Value, error) {
	log.Printf("Make Arg got %s, %v, %s", data, template, mime)
	i := reflect.New(template).Interface()
	log.Printf("Arg type %v", i)

	if data == "" {
		return reflect.ValueOf(i).Elem(), nil
	} else {
		log.Println("Data sent: ", data)
	}

	buf := bytes.NewBufferString(data)
	err := bytesToInterface(buf, i, mime)

	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(i).Elem(), nil
}

func bytesToInterface(buf *bytes.Buffer, i interface{}, mime string) error {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:

		n, err := strconv.ParseBool(buf.String())
		if err != nil {
			return errors.New("Invalid value. " + err.Error())
		}
		reflect.ValueOf(i).Elem().SetBool(n)
		break
	case reflect.String:
		reflect.ValueOf(i).Elem().SetString(buf.String())
		break
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		m := getMarshallerByMime(mime)
		return m.Unmarshal(buf.Bytes(), i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		n, err := strconv.ParseInt(buf.String(), 10, 64)
		if err != nil || v.OverflowInt(n) {
			return errors.New("Invalid value. " + err.Error())
		}
		reflect.ValueOf(i).Elem().SetInt(n)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(buf.String(), 10, 64)
		if err != nil || v.OverflowUint(n) {
			return errors.New("Invalid value. " + err.Error())
		}

		reflect.ValueOf(i).Elem().SetUint(n)
		break
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(buf.String(), v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return errors.New("Invalid value. " + err.Error())
		}
		reflect.ValueOf(i).Elem().SetFloat(n)
	default:
		return errors.New("Type " + v.Type().Name() + " is not handled by GoRest.")
	}
	return nil

}

//A Marshaller represents the two functions used to marshal/unmarshal interfaces back and forth.
type marshaller struct {
	Marshal   func(v interface{}) ([]byte, error)
	Unmarshal func(data []byte, v interface{}) error
}

var marshallers map[string]*marshaller

//Register a Marshaller. These registered Marshallers are shared by the client or servers side usage of gorest.
func registerMarshaller(mime string, m *marshaller) {
	if marshallers == nil {
		marshallers = make(map[string]*marshaller, 0)
	}
	if _, found := marshallers[mime]; !found {
		marshallers[mime] = m
	}
}

//Get an already registered Marshaller
func getMarshallerByMime(mime string) (m *marshaller) {
	if marshallers == nil {
		marshallers = make(map[string]*marshaller, 0)
	}
	m, _ = marshallers[mime]
	return
}

//Predefined Marshallers

//JSON: This makes the JSON Marshaller. The Marshaller uses pkg: json
func newJSONMarshaller() *marshaller {
	m := marshaller{jsonMarshal, jsonUnMarshal}
	return &m
}
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func jsonUnMarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

//XML
func newXMLMarshaller() *marshaller {
	m := marshaller{xmlMarshal, xmlUnMarshal}
	return &m
}
func xmlMarshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}
func xmlUnMarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
