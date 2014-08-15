// RatingModel.go
package model

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	//"mws/db"
	"mws/db/pgreflect"
	"mws/dto"
	"mws/util"
	"reflect"
	"time"
)

type RatingModel struct {
	// Operation fields should match the associated function names (but in all lower case) for the model
	// Model annotations for VDL mapping to stored procedure/functions using key "vdl"
	// 	with the stored procedure function name as the value.
	createrelationship int `vdl:"addrating"`
	//update       int `vdl:"updatetag"`
	//updatestatus int `vdl:"updatetagstatus`
	//delete       int `vdl:"deletetag"`
}

func (tm RatingModel) CreateRelationship(vdlOperation string, account_uid int64, revision_uid int64, rating int64) (foundRating *dto.Rating, retCode int, err error) {
	funcName := util.GetCallerName()
	//newRating := newRef.(dto.Rating)
	//newRating := *newRef\\
	log.Printf("%s Calling %s with %s, %d, %d %d", funcName, reflect.TypeOf(tm), vdlOperation, account_uid, revision_uid, rating)
	//
	//Use field annotations to find method
	//tt := reflect.TypeOf(tm)
	//vdlField, ok := tt.FieldByName(operation)
	//if ok != true {
	//	return nil, -1, errors.New(fmt.Sprintf("%s: VDL Opeartion Field (%s)) not found", funcName, operation))
	//}
	//vdlOperation := vdlField.Tag.Get("vdl")
	if vdlOperation == "" {
		return nil, -1, errors.New(fmt.Sprintf("%s: VDL annotation not set", funcName))
	} else {
		// Get the Postgres Function information
		pgfunc, _ := pgreflect.GetPgFuncInfo(vdlOperation)
		retMap, retCode, err := pgfunc.VariadicScan(
			account_uid,
			revision_uid,
			rating)
		if err != nil || retMap == nil {
			log.Printf("%s Error Calling Postgres Function - %#v", funcName, err)
			return nil, -1, err
		}
		log.Printf("%s Postgres Function Returned- %#v, %d, %#v", funcName, retMap, retCode, err)

		foundRating := &dto.Rating{
			AccountId:  retMap[0]["ret_account_uid"].(int64),
			RevisionId: retMap[0]["ret_revision_uid"].(int64),
			Rating:     retMap[0]["ret_rating"].(int64),
			Created:    &dto.JsonTime{retMap[0]["ret_created"].(time.Time), time.RFC3339},
			Modified:   &dto.JsonTime{retMap[0]["ret_modified"].(time.Time), time.RFC3339},
		}
		log.Printf("%s Excution result %#v", funcName, foundRating)
		return foundRating, retCode, nil
	}

}
