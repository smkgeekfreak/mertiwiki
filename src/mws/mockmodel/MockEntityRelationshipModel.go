// MockEntityRelationshipModel.go
package mockmodel

import (
	//"errors"
	//"log"
	"mws/dto"
	//"reflect"
)

var (
//tagSequenceNumber int = 0
)

type MockEntityRelationshipModel struct {
	createrelationship   int `vdl:"testoperation"`
	retrieverelationship int `vdl:"testopearation"`
}

func (merm MockEntityRelationshipModel) RetrieveRelationship(op string, revisionId int64, tagId int64) (relationship *dto.EntityRelationship, retCode int, err error) {
	return nil, 0, nil
}

func (merm MockEntityRelationshipModel) CreateRelationship(op string, a int64, b int64) (newER *dto.EntityRelationship, retCode int, err error) {
	return nil, 0, nil
}
