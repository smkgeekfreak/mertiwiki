// User.go
package dto

import (
	//"encoding/json"
	"fmt"
	"time"
)

type User struct {
	Id           int64
	Name         string
	Email        string
	PasswordHash string `json:"-"`
	Status       StatusDetail
	Created      *JsonTime `json:",omitempty"`
	Modified     *JsonTime `json:",omitempty"`
}

type UserDetails struct {
	User       //embedded field
	Tags []Tag `json:",omitempty"` //embedded field
}

func InitUserDetails(user User) UserDetails {
	var u = UserDetails{}
	u.User = user
	//	u.TagList = TagList{make([]Tag, 0)}
	return u
}

func NewUser(name string) UserDetails {
	var u = UserDetails{}
	u.Name = name
	u.Status = StatusDetail{INITIALIZED, fmt.Sprint(StatusType(INITIALIZED))}
	u.Tags = make([]Tag, 0)
	return u
}

func (user *UserDetails) SetStatus(code StatusType) {
	var us = StatusDetail{code, fmt.Sprint(StatusType(code))}
	(*user).Status = us
}

type JsonTime struct {
	time.Time
	FormatStr string
}

func (j JsonTime) format() string {
	return j.Time.Format(j.FormatStr)
}

func (j JsonTime) MarshalText() ([]byte, error) {
	return []byte(j.format()), nil
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}
