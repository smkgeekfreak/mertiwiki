// Tag.go
package dto

import (
	"fmt"
)

type Tag struct {
	Id          int64
	Name        string
	Description string
	Status      StatusDetail
	Created     *JsonTime `json:",omitempty"`
	Modified    *JsonTime `json:",omitempty"`
}

//type TagList struct {
//	Tags []Tag
//}

func NewTag(name string, description string) Tag {
	var t = Tag{}
	t.Name = name
	t.Description = description
	t.Status = StatusDetail{INITIALIZED, fmt.Sprint(StatusType(INITIALIZED))}
	return t
}

//func NewTagList(numTags int) TagList {
//	var t = TagList{make([]Tag, numTags)}
//	return t
//}

func (tag *Tag) SetStatus(code StatusType) {
	var ts = StatusDetail{code, fmt.Sprint(StatusType(code))}
	(*tag).Status = ts
}
