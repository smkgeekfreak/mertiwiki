// Status.go
package dto

import (
	"fmt"
)

type StatusType int

//TODO: Refactor to use StatusType instead of StatusDetail.
// Can then just use the String() function to get the text version
type StatusDetail struct {
	StatusCode StatusType
	Desc       string
}

const (
	INITIALIZED StatusType = iota
	PENDING
	ACTIVE
	DEACTIVATED
	BANNED
	DELETED
)

var StatusDescMap = map[StatusType]string{
	INITIALIZED: "Initalized",
	PENDING:     "Pending",
	ACTIVE:      "Active",
	DEACTIVATED: "Deactivated",
	BANNED:      "Banned",
	DELETED:     "Deleted",
}

func (st StatusType) String() string {
	if v, ok := StatusDescMap[StatusType(st)]; ok {
		return v
	}
	return fmt.Sprintf("Unknown StatusType[%d]", st)
}

//TODO: Can this be refactored into an generic function?
func SetStatus(code StatusType) StatusDetail {
	var sd = StatusDetail{code, fmt.Sprint(StatusType(code))}
	//(*page).Status = sd
	return sd
}
