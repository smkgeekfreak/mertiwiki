// MockTagModel.go
package mockmodel

import (
	//"errors"
	"log"
	"mws/dto"
	"reflect"
)

var (
	tagSequenceNumber int = 2000
)

type MockTagModel struct {
}

func (mpm MockTagModel) Create(tag *dto.Tag) (newId dto.Tag, retCode int, err error) {
	log.Printf("Mocking call Create this type %s", reflect.TypeOf(tag))
	retCode = -1
	tag.Id = int64(tagSequenceNumber + 1)
	newId = *tag
	err = nil
	for _, found := range tagStore {
		if tag.Id == found.Id {
			log.Printf("Tag with Id = %v already exists", tag.Id)
			retCode = 1
			return
		}
	}
	tagStore[tag.Name] = *tag
	tagSequenceNumber = tagSequenceNumber + 1
	log.Printf("Created New Tag : %d, %s, %s, %s", tag.Id, tag.Name, tag.Description, tag.Status)
	retCode = 0
	return
}
