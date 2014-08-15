// Rating.go
package dto

import (
//"fmt"
)

type Rating struct {
	AccountId  int64
	RevisionId int64
	Rating     int64
	Created    *JsonTime `json:",omitempty"`
	Modified   *JsonTime `json:",omitempty"`
}
