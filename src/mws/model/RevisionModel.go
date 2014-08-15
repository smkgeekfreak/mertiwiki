// RevisionModel.go
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

type RevisionModel struct {
	create             int `vdl:"create_revision"`
	update             int `vdl:"update_revision"`
	updatestatus       int `vdl:"update_revisionstatus"`
	delete             int `vdl:"update_revisionstatus"`
	findbyid           int `vdl:"findbyid_revision"`
	findbyauthor       int `vdl:"findbyauthor_revision"`
	authorizeOwnership int `vdl:"authorizeOwnership_revision"`
}

func (revm RevisionModel) FindByAuthor(vdlOperation string, userId int64) (*[]dto.Revision, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(revm), userId)
	return revm.processVDLMultiResults(vdlOperation, userId)
}

func (revm RevisionModel) FindById(vdlOperation string, revisionId int64) (*dto.Revision, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(revm), revisionId)
	return revm.processVDLCall(vdlOperation, revisionId)
}

func (revm RevisionModel) Create(vdlOperation string, authId int64, newRef *dto.Revision) (*dto.Revision, int, error) {
	funcName := util.GetCallerName()
	newRevision := *newRef

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(revm), newRevision)
	return revm.processVDLCall(
		vdlOperation,
		newRevision.SectionId,
		authId,
		newRevision.Content,
		dto.INITIALIZED)
}

func (revm RevisionModel) Update(vdlOperation string, revisionId int64, refRevision *dto.Revision) (*dto.Revision, int, error) {
	funcName := util.GetCallerName()
	modRevision := *refRevision

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(revm), modRevision)
	return revm.processVDLCall(vdlOperation, revisionId, modRevision.Content, modRevision.Status.StatusCode)
}

func (revm RevisionModel) UpdateStatus(vdlOperation string, revisionId int64, status dto.StatusType) (*dto.Revision, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(revm), revisionId)
	return revm.processVDLCall(vdlOperation, revisionId, status)
}

func (revm RevisionModel) Delete(vdlOperation string, revisionId int64) (*dto.Revision, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(revm), revisionId)
	return revm.processVDLCall(vdlOperation, revisionId, dto.DELETED)
}

func (revm RevisionModel) processVDLMultiResults(vdlOperation string, variadicArgs ...interface{}) (*[]dto.Revision, int, error) {
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

	dbRevisions := make([]dto.Revision, 0)
	for _, row := range retMap {
		//construct transfer object with return values
		dbRevision := dto.Revision{
			Id:        row["uid"].(int64),
			SectionId: row["sec_uid"].(int64),
			PageId:    row["page_uid"].(int64),
			AuthorId:  row["user_uid"].(int64),
			Content:   row["body"].(string),
			Status:    dto.SetStatus(dto.StatusType(row["status"].(int64))),
			Created:   &dto.JsonTime{row["created"].(time.Time), time.RFC3339},
			Modified:  &dto.JsonTime{row["modified"].(time.Time), time.RFC3339},
		}
		log.Printf("%s -- created %v, modified %v", funcName, row["created"].(time.Time), row["modified"].(time.Time).Format(time.RFC3339))
		log.Printf("%s -- Excution result %#v", funcName, dbRevision)
		dbRevisions = append(dbRevisions, dbRevision)
	}
	return &dbRevisions, retCode, nil
}

func (revm RevisionModel) processVDLCall(vdlOperation string, variadicArgs ...interface{}) (*dto.Revision, int, error) {
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
		dbRevision := &dto.Revision{
			Id:        retMap[0]["ret_uid"].(int64),
			SectionId: retMap[0]["ret_sec_uid"].(int64),
			PageId:    retMap[0]["ret_page_uid"].(int64),
			AuthorId:  retMap[0]["ret_user_uid"].(int64),
			Content:   retMap[0]["ret_body"].(string),
			Status:    dto.SetStatus(dto.StatusType(retMap[0]["ret_status"].(int64))),
			Created:   &dto.JsonTime{retMap[0]["ret_created"].(time.Time), time.RFC3339},
			Modified:  &dto.JsonTime{retMap[0]["ret_modified"].(time.Time), time.RFC3339},
		}
		log.Printf("%s -- created %v, modified %v", funcName, retMap[0]["ret_created"].(time.Time), retMap[0]["ret_modified"].(time.Time).Format(time.RFC3339))
		log.Printf("%s -- Excution result %#v", funcName, dbRevision)
		return dbRevision, retCode, nil
	}

	return nil, 1, errors.New("No data returned")
}
