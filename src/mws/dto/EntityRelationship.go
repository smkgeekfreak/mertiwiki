// EntityRelationship.go
package dto

import (
	"reflect"
)

type EntityRelationship struct {
	RelId1, RelId2     int64                  //pk, fk's
	RelName1, RelName2 string                 `json:",omitempty"`
	RelType1, RelType2 reflect.Type           `json:"-"`
	Values             map[string]interface{} `json:",omitempty"`
	Status             StatusDetail
	Created            *JsonTime `json:",omitempty"`
	Modified           *JsonTime `json:",omitempty"`
}
