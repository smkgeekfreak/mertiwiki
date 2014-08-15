// PageModel.go
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

type PageModel struct {
	// Model annotations for VDL mapping to stored procedure/functions
	//vdl_opearations int  `create:"addpage" update:"updatepage" delete:"deletepage"`
	create             int `vdl:"create_page"`
	update             int `vdl:"update_page"`
	updatestatus       int `vdl:"update_pagestatus"`
	delete             int `vdl:"update_pagestatus"`
	findbyid           int `vdl:"findbyid_page"`
	findbyauthor       int `vdl:"findbyauthor_page"`
	authorizeOwnership int `vdl:"authorizeOwnership_page"`
}

func (pm PageModel) FindByAuthor(vdlOperation string, userId int64) (*[]dto.Page, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(pm), userId)
	return pm.processVDLMultiResults(vdlOperation, userId)
}

func (pm PageModel) FindById(vdlOperation string, uid int64) (*dto.Page, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(pm), uid)
	return pm.processVDLCall(vdlOperation, uid)
}

func (pm PageModel) Create(vdlOperation string, authId int64, refDTO *dto.Page) (*dto.Page, int, error) {
	funcName := util.GetCallerName()
	newDTO := *refDTO

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(pm), newDTO)
	return pm.processVDLCall(
		vdlOperation,
		authId,
		newDTO.Title,
		dto.INITIALIZED)
}

func (pm PageModel) Update(vdlOperation string, uid int64, refDTO *dto.Page) (*dto.Page, int, error) {
	funcName := util.GetCallerName()
	modDTO := *refDTO

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(pm), modDTO)
	return pm.processVDLCall(vdlOperation, uid, modDTO.Title, modDTO.Status.StatusCode)
}

func (pm PageModel) UpdateStatus(vdlOperation string, uid int64, status dto.StatusType) (*dto.Page, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(pm), uid)
	return pm.processVDLCall(vdlOperation, uid, status)
}

func (pm PageModel) Delete(vdlOperation string, uid int64) (*dto.Page, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(pm), uid)
	return pm.processVDLCall(vdlOperation, uid, dto.DELETED)
}

func (pm PageModel) processVDLMultiResults(vdlOperation string, variadicArgs ...interface{}) (*[]dto.Page, int, error) {
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
	if err != nil || retMap == nil || len(retMap) == 0 {
		log.Printf("%s -- Error Calling Postgres Function - %s ( %#v)", funcName, pgfunc.Name, pgfunc)
		return nil, -1, err
	}
	log.Printf("%s -- Postgres Function Returned- %#v, %d, %#v", funcName, retMap, retCode, err)

	dbResults := make([]dto.Page, 0)
	for _, row := range retMap {
		//construct transfer object with return values
		dbResult := dto.Page{
			Id:       row["uid"].(int64),
			AuthorId: row["user_uid"].(int64),
			Title:    row["title"].(string),
			Status:   dto.SetStatus(dto.StatusType(row["status"].(int64))),
			Created:  &dto.JsonTime{row["created"].(time.Time), time.RFC3339},
			Modified: &dto.JsonTime{row["modified"].(time.Time), time.RFC3339},
		}
		log.Printf("%s -- created %v, modified %v", funcName, row["created"].(time.Time), row["modified"].(time.Time).Format(time.RFC3339))
		log.Printf("%s -- Excution result %#v", funcName, dbResult)
		dbResults = append(dbResults, dbResult)
	}
	return &dbResults, retCode, nil
}

func (rm PageModel) processVDLCall(vdlOperation string, variadicArgs ...interface{}) (*dto.Page, int, error) {
	funcName := util.GetCallerName()
	if vdlOperation == "" {
		return nil, -1, errors.New(fmt.Sprintf("%s: VDL annotation not set", funcName))
	}

	log.Printf("%s -- Calling operation %s with %v", funcName, vdlOperation, variadicArgs)
	// Get the Postgres Function information
	pgfunc, err := pgreflect.GetPgFuncInfo(vdlOperation)
	if err != nil {
		log.Printf("%s -- Error Getting Postgres Function Information- %s ", funcName, vdlOperation)
		return nil, -1, err
	}
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
	//construct transfer object with return values
	if len(retMap) > 0 {
		dbResults := &dto.Page{
			Id:       retMap[0]["ret_uid"].(int64),
			AuthorId: retMap[0]["ret_user_uid"].(int64),
			Title:    retMap[0]["ret_title"].(string),
			Status:   dto.SetStatus(dto.StatusType(retMap[0]["ret_status"].(int64))),
			Created:  &dto.JsonTime{retMap[0]["ret_created"].(time.Time), time.RFC3339},
			Modified: &dto.JsonTime{retMap[0]["ret_modified"].(time.Time), time.RFC3339},
		}
		log.Printf("%s -- created %v, modified %v", funcName, retMap[0]["ret_created"].(time.Time), retMap[0]["ret_modified"].(time.Time).Format(time.RFC3339))
		log.Printf("%s -- Excution result %#v", funcName, dbResults)
		return dbResults, retCode, nil
	}

	return nil, 1, errors.New("No data returned")
}
