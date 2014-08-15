// MockSectionModel.go
package mockmodel

import (
	//"errors"
	"log"
	"mws/dto"
	"reflect"
)

var (
	secSequenceNumber int = 0
)

type MockSectionModel struct {
	create int `vdl:"testoperation"`
}

func (mpm MockSectionModel) Create(vdlOperation string, authId int64, sec dto.Section) (newId dto.Section, retCode int, err error) {
	log.Printf("Mocking call Create this type %s", reflect.TypeOf(sec))
	retCode = 409
	sec.Id = int64(secSequenceNumber + 1)
	newId = sec
	err = nil
	for _, found := range sectionStore {
		if sec.Id == found.Id {
			log.Printf("Section with Id = %v already exists", sec.Id)
			retCode = 302
			return
		}
	}
	sectionStore[sec.Id] = sec
	secSequenceNumber = secSequenceNumber + 1
	log.Printf("Created New Section : %#v, %#v, %#v", sec.Id, sec.PageId, sec.Name)
	retCode = 201
	return
}

func (msm MockSectionModel) Update(refSec *dto.Section) (retSec *dto.Section, retCode int, err error) {
	log.Printf("Mocking call Update. This section updated: %v", refSec.Id)
	return refSec, 0, nil
}
func (msm MockSectionModel) Delete(uid int64) (int64, int, error) {
	return uid, 0, nil
}
