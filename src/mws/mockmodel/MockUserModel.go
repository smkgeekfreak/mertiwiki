// MockUserDB.go

package mockmodel

import (
	"fmt"
	"log"
	"mws/dto"
	"reflect"
)

type MockUserModel struct {
	create int `vdl:"testoperation"`
}

func (mum MockUserModel) RetrieveById(id int64) (*dto.UserDetails, error) {
	log.Printf("Calling Retrieve with %d", id)
	user := dto.User{
		Id: id,
	}
	u, err := mum.Retrieve(user)
	return u, err
}

func (mum MockUserModel) Retrieve(u dto.User) (*dto.UserDetails, error) {
	log.Printf("Mocking Call %s with %s", reflect.TypeOf(mum), u.Id)
	if user, found := userStore[u.Id]; found {
		//details := wikidto.UserDetails{user, dto.NewTags(0)}
		var details = dto.InitUserDetails(user)
		//
		// Get userTags
		if tagRefs, has := userTags[u.Id]; has {
			for _, v := range tagRefs {
				if tag, tFound := tagStore[v]; tFound {
					details.Tags = append(details.Tags, tag)
				}
			}
		}
		log.Printf("Returning %#v", details)
		return &details, nil
	} else {
		return nil, fmt.Errorf("User %v not found", u.Name)
	}

}

func (um MockUserModel) Create(vdlOpeartion string, uRef *dto.User) (newUser *dto.User, retCode int, err error) {
	u := *uRef
	log.Printf("Mocking call Create this type %v", reflect.TypeOf(u))

	retCode = -1
	for key, found := range userStore {
		if found.Name == u.Name {
			userStore[key] = u
			retCode = 0 //change after unique name constraint
			return &u, retCode, nil
		}
	}
	//User not in database, so create new
	userStore[u.Id] = u
	log.Printf("New User:%#v", u)
	retCode = 0
	//serv.ResponseBuilder().Created("http://localhost:8686/tag/Username1") //Created, http 201
	return &u, retCode, nil
}

func (um MockUserModel) Update(userId int64, refUser *dto.User) (retUser *dto.User, retCode int, err error) {
	return refUser, 0, nil
}
func (um MockUserModel) Delete(uid int64) (newUser *dto.User, retCode int, err error) {
	return &dto.User{}, 0, nil
}
