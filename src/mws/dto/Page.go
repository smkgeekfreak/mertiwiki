// Page.go
package dto

import (
	//"errors"
	"fmt"
	//"net/url"
	//	"strings"
)

//type Title string

//func (p *Page) SetPath() {
//	p.Path = "/" + strings.Replace(string(p.Title), " ", "", -1)
//}

type Page struct {
	Id       int64
	AuthorId int64 `json:",omitempty"` // need to fix tests and remove this
	Title    string
	Status   StatusDetail
	Tags     []Tag     `json:",omitempty"`
	Created  *JsonTime `json:",omitempty"`
	Modified *JsonTime `json:",omitempty"`
}

type Section struct {
	Id, PageId int64
	AuthorId   int64
	Name       string
	OrderNum   int
	Status     StatusDetail
	Tags       []Tag     `json:",omitempty"`
	Created    *JsonTime `json:",omitempty"`
	Modified   *JsonTime `json:",omitempty"`
}

type Revision struct {
	Id, SectionId, PageId int64 //pk, fk's
	AuthorId              int64
	Content               string
	Tags                  []Tag `json:",omitempty"`
	Status                StatusDetail
	Created               *JsonTime `json:",omitempty"`
	Modified              *JsonTime `json:",omitempty"`
}

func NewPage(id int64, title string) Page {
	var p = Page{}
	p.Id = id
	p.Title = title
	p.Status = StatusDetail{INITIALIZED, fmt.Sprint(StatusType(INITIALIZED))}
	//	p.SetPath()
	return p
}
