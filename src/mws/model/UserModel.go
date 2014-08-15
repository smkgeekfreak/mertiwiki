// UserModel.go
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

type UserModel struct {
	// Model annotations for VDL mapping to stored procedure/functions
	//vdl_opearations int  `create:"addpage" update:"updatepage" delete:"deletepage"`
	create             int `vdl:"create_account"`
	update             int `vdl:"update_account"`
	updatestatus       int `vdl:"update_accountstatus"`
	delete             int `vdl:"update_accountstatus"`
	findbyid           int `vdl:"findbyid_account"`
	authorizeOwnership int `vdl:"authorizeOwnership_account"`
	findrating         int `vdl:"find_user_rating"`
}

func (um UserModel) FindRating(vdlOperation string, uid int64) (*dto.UserRating, int, error) {
	retMap, retCode, err := processGenericVDL(vdlOperation, uid)
	if err != nil {
		return nil, -1, err
	}
	if len(retMap) > 0 {
		result := &dto.UserRating{
			Uid:     retMap[0]["ret_user_uid"].(int64),
			Rating:  retMap[0]["ret_user_rating"].(int64),
			Updated: &dto.JsonTime{retMap[0]["ret_updated"].(time.Time), time.RFC3339},
		}
		return result, retCode, nil
	}
	return nil, 1, errors.New("No data returned")

}

func (sm UserModel) FindById(vdlOperation string, uid int64) (*dto.User, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), uid)
	return sm.processVDLCall(vdlOperation, uid)
}

func (sm UserModel) Create(vdlOperation string, refDTO *dto.User) (*dto.User, int, error) {
	funcName := util.GetCallerName()
	newDTO := *refDTO

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), newDTO)
	return sm.processVDLCall(
		vdlOperation,
		newDTO.Name,
		newDTO.PasswordHash,
		newDTO.Email,
		dto.INITIALIZED)
}

func (sm UserModel) Update(vdlOperation string, uid int64, refDTO *dto.User) (*dto.User, int, error) {
	funcName := util.GetCallerName()
	modDTO := *refDTO

	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), modDTO)
	return sm.processVDLCall(vdlOperation, uid, modDTO.Email, modDTO.Status.StatusCode)
}

func (sm UserModel) UpdateStatus(vdlOperation string, uid int64, status dto.StatusType) (*dto.User, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), uid)
	return sm.processVDLCall(vdlOperation, uid, status)
}

func (sm UserModel) Delete(vdlOperation string, uid int64) (*dto.User, int, error) {
	funcName := util.GetCallerName()
	log.Printf("%s Calling %s with %d", funcName, reflect.TypeOf(sm), uid)
	return sm.processVDLCall(vdlOperation, uid, dto.DELETED)
}

func (sm UserModel) processVDLMultiResults(vdlOperation string, variadicArgs ...interface{}) (*[]dto.User, int, error) {
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

	dbResults := make([]dto.User, 0)
	for _, row := range retMap {
		//construct transfer object with return values
		dbResult := dto.User{
			Id:       row["uid"].(int64),
			Name:     row["username"].(string),
			Email:    row["email"].(string),
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

func (rm UserModel) processVDLCall(vdlOperation string, variadicArgs ...interface{}) (*dto.User, int, error) {
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
		dbResults := &dto.User{
			Id:       retMap[0]["ret_uid"].(int64),
			Name:     retMap[0]["ret_username"].(string),
			Email:    retMap[0]["ret_email"].(string),
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
