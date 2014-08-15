//// TagRevisionModel.go
package model

import (
	//"fmt"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"mws/db"
	"mws/db/pgreflect"
	"mws/dto"
	"mws/util"
	"reflect"
	"time"
)

type TagRevisionModel struct {
	createrelationship int `vdl:"addtagrevision"`
	checkauthorization int `vdl:"revisionauthorization"`
}

func (tm TagRevisionModel) RetrieveRelationship(revisionId int64, tagId int64) (relationship *dto.EntityRelationship, retCode int, err error) {

	funcName := util.GetCallerName()
	log.Printf("Calling %s (%d, %d)", funcName, tagId, revisionId)

	db, err := db.GetConnection()
	defer db.Close()
	if err != nil {
		log.Printf("%s Connection err %#v", funcName, err)
		return nil, -1, err
	}
	log.Printf("%s Connected to Postgres", funcName)
	defer db.Close()
	stmt, err := db.Prepare(`SELECT * FROM tag_revisiondata WHERE revision_uid =$1 AND tag_uid = $2`)

	if err != nil {
		log.Printf("Prepare err %#v", err)
		//stmt.Close()
		//db.Close()
		return nil, -1, err
	}
	log.Printf("%s Executing ", funcName)
	result, err := stmt.Query(
		revisionId,
		tagId)

	if err != nil {
		log.Printf("%s Execution problem: %#v", funcName, result)
		//result.Close()
		//stmt.Close()
		//db.Close()
		return nil, -1, err
	}
	for result.Next() {
		var status int
		var created, modified time.Time
		if err := result.Scan(&tagId, &revisionId, &status, &created, &modified); err != nil {
			log.Printf("%s Results Error %#v", funcName, err)
			//result.Close()
			//stmt.Close()
			//db.Close()
			return nil, -1, err
		}
		log.Printf("%s Processing return values for entity relationship ids %d, %d", funcName, tagId, revisionId)
		newTRRelationship := &dto.EntityRelationship{
			RelId2:   tagId,
			RelId1:   revisionId,
			RelName1: "dto.Revision",
			RelName2: "dto.Tag",
			RelType2: reflect.TypeOf(dto.Tag{}),
			RelType1: reflect.TypeOf(dto.Revision{}),
			Status:   dto.SetStatus(dto.StatusType(status)),
			Created:  &dto.JsonTime{created, time.RFC3339},
			Modified: &dto.JsonTime{created, time.RFC3339},
		}
		//result.Close()
		//stmt.Close()
		//db.Close()
		log.Printf("%s excution result %v", funcName, newTRRelationship)
		return newTRRelationship, 0, nil
	}

	//result.Close()
	//stmt.Close()
	//db.Close()
	return nil, 1, nil
}

func (trm TagRevisionModel) CreateRelationship(vdlOperation string, revisionId int64, tagId int64) (relationship *dto.EntityRelationship, retCode int, err error) {
	funcName := util.GetCallerName()
	log.Printf("Calling %s (%d, %d)", funcName, tagId, revisionId)
	return trm.processVDLCall(vdlOperation, revisionId, tagId, dto.ACTIVE)
}

func (trm TagRevisionModel) processVDLCall(vdlOperation string, variadicArgs ...interface{}) (relationship *dto.EntityRelationship, retCode int, err error) {
	funcName := util.GetCallerName()
	if vdlOperation == "" {
		return nil, -1, errors.New(fmt.Sprintf("%s: VDL annotation not set", funcName))
	}
	log.Printf("%s -- Calling operation %s with %v", funcName, vdlOperation, variadicArgs)
	// Get the Postgres Function information
	pgfunc, _ := pgreflect.GetPgFuncInfo(vdlOperation)
	//Call the function and get return values
	retMap, retCode, err := pgfunc.VariadicScan(variadicArgs...)
	if err != nil || retMap == nil {
		log.Printf("%s -- Error Calling Postgres Function - %s ( %#v)", funcName, pgfunc.Name, pgfunc)
		return nil, -1, err
	}
	log.Printf("%s -- Postgres Function Returned- %#v, %d, %#v", funcName, retMap, retCode, err)
	//construct transfer object with return values
	newTRRelationship := &dto.EntityRelationship{
		RelId2:   retMap[0]["ret_tag_uid"].(int64),
		RelId1:   retMap[0]["ret_revision_uid"].(int64),
		RelName1: "dto.Revision",
		RelName2: "dto.Tag",
		RelType2: reflect.TypeOf(dto.Tag{}),
		RelType1: reflect.TypeOf(dto.Revision{}),
		Status:   dto.SetStatus(dto.StatusType(retMap[0]["ret_status"].(int64))),
		Created:  &dto.JsonTime{retMap[0]["ret_created"].(time.Time), time.RFC3339},
		Modified: &dto.JsonTime{retMap[0]["ret_modified"].(time.Time), time.RFC3339},
	}
	log.Printf("%s -- created %v, modified %v", funcName, retMap[0]["ret_created"].(time.Time), retMap[0]["ret_modified"].(time.Time).Format(time.RFC3339))
	log.Printf("%s -- Excution result %#v", funcName, newTRRelationship)
	return newTRRelationship, retCode, nil
}
