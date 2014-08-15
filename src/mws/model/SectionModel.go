// SectionModel.go
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

type SectionModel struct {
	// Model annotations for VDL mapping to stored procedure/functions
	//vdl_opearations int  `create:"addpage" update:"updatepage" delete:"deletepage"`
	create             int `vdl:"create_section"`
	update             int `vdl:"update_section"`
	updatestatus       int `vdl:"update_sectionstatus"`
	delete             int `vdl:"update_sectionstatus"`
	findbyid           int `vdl:"findbyid_section"`
	findbyauthor       int `vdl:"findbyauthor_section"`
	authorizeOwnership int `vdl:"authorizeOwnership_section"`
}

func (sm SectionModel) FindByAuthor(vdlOperation string, userId int64) (*[]dto.Section, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), userId)
	return sm.processVDLMultiResults(vdlOperation, userId)
}

func (sm SectionModel) FindById(vdlOperation string, uid int64) (*dto.Section, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), uid)
	return sm.processVDLCall(vdlOperation, uid)
}

func (sm SectionModel) Create(vdlOperation string, authId int64, refDTO *dto.Section) (*dto.Section, int, error) {
	funcName := util.GetCallerName()
	newDTO := *refDTO

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), newDTO)
	return sm.processVDLCall(
		vdlOperation,
		newDTO.PageId,
		authId,
		newDTO.Name,
		newDTO.OrderNum,
		dto.INITIALIZED)
}

func (sm SectionModel) Update(vdlOperation string, uid int64, refDTO *dto.Section) (*dto.Section, int, error) {
	funcName := util.GetCallerName()
	modDTO := *refDTO

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), modDTO)
	return sm.processVDLCall(vdlOperation, uid, modDTO.Name, modDTO.OrderNum, modDTO.Status.StatusCode)
}

func (sm SectionModel) UpdateStatus(vdlOperation string, uid int64, status dto.StatusType) (*dto.Section, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), uid)
	return sm.processVDLCall(vdlOperation, uid, status)
}

func (sm SectionModel) Delete(vdlOperation string, uid int64) (*dto.Section, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), uid)
	return sm.processVDLCall(vdlOperation, uid, dto.DELETED)
}

func (sm SectionModel) processVDLMultiResults(vdlOperation string, variadicArgs ...interface{}) (*[]dto.Section, int, error) {
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

	dbResults := make([]dto.Section, 0)
	for _, row := range retMap {
		//construct transfer object with return values
		dbResult := dto.Section{
			Id:       row["uid"].(int64),
			PageId:   row["page_uid"].(int64),
			AuthorId: row["user_uid"].(int64),
			Name:     row["name"].(string),
			OrderNum: int(row["ordernum"].(int64)),
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

func (rm SectionModel) processVDLCall(vdlOperation string, variadicArgs ...interface{}) (*dto.Section, int, error) {
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
		dbResults := &dto.Section{
			Id:       retMap[0]["ret_uid"].(int64),
			PageId:   retMap[0]["ret_page_uid"].(int64),
			AuthorId: retMap[0]["ret_user_uid"].(int64),
			Name:     retMap[0]["ret_name"].(string),
			OrderNum: int(retMap[0]["ret_ordernum"].(int64)),
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
