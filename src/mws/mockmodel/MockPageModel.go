// MockPageModel.go
package mockmodel

import (
	//"errors"
	"fmt"
	"log"
	"mws/dto"
	"reflect"
)

var (
	pageSequenceNumber int = 0
)

type MockPageModel struct {
	create int `vdl:"testoperation"`
	delete int `vdl:"testopearation"`
}

func (mpm MockPageModel) RetrieveById(id int64) (*dto.Page, error) {
	log.Printf("Calling Retrieve with %d", id)
	u, err := mpm.Retrieve(id)
	return u, err
}

func (mpm MockPageModel) Retrieve(findId int64) (*dto.Page, error) {
	log.Printf("Mocking Call %v with %v", reflect.TypeOf(mpm), findId)
	//serv.ResponseBuilder().SetResponseCode(409)
	for _, p := range pageStore {
		if p.Id == findId {
			log.Printf("Found Page:%#v", p)
			//serv.ResponseBuilder().SetResponseCode(200)
			return &p, nil
		}

	}
	return nil, fmt.Errorf("Page with Id: %d not found", findId)
}

func (mpm MockPageModel) Create(vdlOperation string, authId int64, pRef *dto.Page) (foundPage *dto.Page, retCode int, err error) {
	p := *pRef
	log.Printf("Mocking call Create this type %s", reflect.TypeOf(p))

	retCode = -1
	for _, found := range pageStore {
		if p.Title == found.Title {
			log.Printf("Page with Title = %v already exists", p.Title)
			retCode = 1
			return &p, retCode, nil
		}
	}

	var newPage = dto.NewPage(1, string(p.Title))
	pageStore[1] = newPage
	log.Printf("Created New Page :%s\n %#v", p.Title, newPage)

	retCode = 0
	return &newPage, retCode, nil
}

func (mpm MockPageModel) Update(pageId int64, refPage *dto.Page) (retPage *dto.Page, retCode int, err error) {
	return refPage, 0, nil
}
func (mpm MockPageModel) Delete(uid int64) (int64, int, error) {
	return uid, 0, nil
}
