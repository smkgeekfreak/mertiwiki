// TagModel.go
package model

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"mws/db/pgreflect"
	"mws/dto"
	"mws/util"
	"reflect"
	"time"
)

func GetVDLOperation(model interface{}, operationField string) (string, error) {
	modelType := reflect.TypeOf(model)
	vdlField, hasVdl := modelType.FieldByName(operationField)
	if hasVdl { // == true
		//
		//If vdl operation provided then prepend it to the argument list
		vdlOperation := vdlField.Tag.Get("vdl")
		if len(vdlOperation) == 0 {
			log.Printf("%v VDL Tag (%s)) not annotated", model, operationField)
			return "", errors.New(fmt.Sprintf("%v VDL Opeartion Field (%s)) not found", model, operationField))
		}
		return vdlOperation, nil
	} else {
		log.Printf("%v VDL Opeartion Field (%s)) not found", model, operationField)
		return "", errors.New(fmt.Sprintf("%v VDL Opeartion Field (%s)) not found", model, operationField))
	}
}

type TagModel struct {
	// Operation fields should match the associated function names (but in all lower case) for the model
	// Model annotations for VDL mapping to stored procedure/functions using key "vdl"
	// 	with the stored procedure function name as the value.
	create             int `vdl:"create_tag"`
	update             int `vdl:"update_tag"`
	updatestatus       int `vdl:"update_tagstatus"`
	delete             int `vdl:"update_tagstatus"`
	findbyid           int `vdl:"findbyid_tag"`
	authorizeOwnership int `vdl:"authorizeOwnership_tag"`
}

func (tm TagModel) FindById(vdlOperation string, uid int64) (*dto.Tag, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(tm), uid)
	return tm.processVDLCall(vdlOperation, uid)
}

func (tm TagModel) Create(vdlOperation string, newRef *dto.Tag) (foundTag *dto.Tag, retCode int, err error) {
	funcName := util.GetCallerName()
	newTag := *newRef
	log.Printf("%s Calling %s with %s", funcName, reflect.TypeOf(tm), newTag)
	return tm.processVDLCall(vdlOperation, newTag.Name, newTag.Description, newTag.Status.StatusCode)
}

func (tm TagModel) Update(vdlOperation string, uid int64, refTag *dto.Tag) (retTag *dto.Tag, retCode int, err error) {
	funcName := util.GetCallerName()
	modTag := *refTag
	log.Printf("%s -- Calling %s with %s", funcName, reflect.TypeOf(tm), modTag)
	return tm.processVDLCall(vdlOperation, uid, modTag.Name, modTag.Description, modTag.Status.StatusCode)
}

func (tm TagModel) UpdateStatus(vdlOperation string, tagId int64, upStatus dto.StatusType) (tag *dto.Tag, retCode int, err error) {
	funcName := util.GetCallerName()
	log.Printf("%s -- Calling %s with %d", funcName, reflect.TypeOf(tm), tagId)
	return tm.processVDLCall(vdlOperation, tagId, upStatus)
}

func (tm TagModel) Delete(vdlOperation string, tagId int64) (tag *dto.Tag, retCode int, err error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %s", funcName, reflect.TypeOf(tm), tagId)
	return tm.processVDLCall(vdlOperation, tagId, dto.DELETED)
}

func (tm TagModel) processVDLCall(vdlOperation string, variadicArgs ...interface{}) (tag *dto.Tag, retCode int, err error) {
	funcName := util.GetCallerName()
	if vdlOperation == "" {
		return nil, -1, errors.New(fmt.Sprintf("%s: VDL annotation not set", funcName))
	}

	log.Printf("%s -- Calling operation %s with %v", funcName, vdlOperation, variadicArgs)
	// Get the Postgres Function information
	pgfunc, _ := pgreflect.GetPgFuncInfo(vdlOperation)
	//
	// Check the argument list size agrees
	if pgfunc.NumArgs != len(variadicArgs) {
		log.Printf("%s -- Error with number of argument for %s has %d, passing in %d", funcName, vdlOperation, pgfunc.NumArgs, len(variadicArgs))
		return nil, -1, errors.New(fmt.Sprintf("Error with number of argument for %s has %d, passing in %d", vdlOperation, pgfunc.NumArgs, len(variadicArgs)))
	}
	//Call the function and get return values
	retMap, retCode, err := pgfunc.VariadicScan(variadicArgs...)
	if err != nil || retMap == nil {
		log.Printf("%s -- Error Calling Postgres Function - %s ( %#v)", funcName, pgfunc.Name, pgfunc)
		return nil, -1, err
	}
	log.Printf("%s -- Postgres Function Returned- %#v, %d, %#v", funcName, retMap, retCode, err)
	if len(retMap) > 0 {
		//construct transfer object with return values
		foundTag := &dto.Tag{
			Id:          retMap[0]["ret_uid"].(int64),
			Name:        retMap[0]["ret_name"].(string),
			Description: retMap[0]["ret_description"].(string),
			Status:      dto.SetStatus(dto.StatusType(retMap[0]["ret_status"].(int64))),
			Created:     &dto.JsonTime{retMap[0]["ret_created"].(time.Time), time.RFC3339},
			Modified:    &dto.JsonTime{retMap[0]["ret_modified"].(time.Time), time.RFC3339},
		}
		log.Printf("%s -- created %v, modified %v", funcName, retMap[0]["ret_created"].(time.Time), retMap[0]["ret_modified"].(time.Time).Format(time.RFC3339))
		log.Printf("%s -- Excution result %#v", funcName, foundTag)
		return foundTag, retCode, nil
	}
	return nil, 1, errors.New("No data returned")
}
