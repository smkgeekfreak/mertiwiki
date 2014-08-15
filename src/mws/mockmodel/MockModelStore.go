// MockModelStore.go
package mockmodel

import (
	//"crypto/rand"
	//"encoding/hex"
	"fmt"
	"mws/dto"
	"strconv"
)

var (
	userStore    map[int64]dto.User
	tagStore     map[string]dto.Tag
	userTags     map[int64][]string
	pageStore    map[int64]dto.Page
	sectionStore map[int64]dto.Section
)

func init() {
	userStore = make(map[int64]dto.User, 0)
	tagStore = make(map[string]dto.Tag, 0)
	pageStore = make(map[int64]dto.Page, 0)
	sectionStore = make(map[int64]dto.Section, 0)
	userTags = make(map[int64][]string)
	initUsers()
	initTags()
	initPages()
}

func initPages() {
	var uid = int64(99999)
	var title = "Sample Page Demo"
	pageStore[uid] = dto.Page{
		Id:     uid,
		Title:  title,
		Status: dto.StatusDetail{dto.ACTIVE, fmt.Sprint(dto.StatusType(dto.PENDING))},
		//TagList: dto.NewTagList(0),
	}
}

func initUsers() {
	for i := 1; i <= 10; i++ {
		var name = "Testing Create Should Delete"
		if i > 1 {
			name = "Testing Create Should Delete" + strconv.Itoa(i)
		}
		var uid = int64(i)
		userStore[uid] = dto.User{
			Id:           uid,
			Name:         name,
			Email:        "thisisatest@go.com",
			PasswordHash: "13lafjalsdflasfkdjlf",
			Status:       dto.StatusDetail{dto.ACTIVE, fmt.Sprint(dto.StatusType(dto.ACTIVE))},
		}
		if i%2 == 0 {
			userTags[uid] = append(userTags[uid], "Tech")
			userTags[uid] = append(userTags[uid], "Math")
		} else {
			userTags[uid] = append(userTags[uid], "Health")
			userTags[uid] = append(userTags[uid], "Science")
		}
	}
}
func initTags() {
	tagStore["Tech"] = dto.Tag{
		Id:          1,
		Name:        "Tech",
		Description: "Technology",
		Status: dto.StatusDetail{
			dto.ACTIVE,
			fmt.Sprint(dto.StatusType(dto.ACTIVE)),
		},
	}
	tagStore["Science"] = dto.Tag{
		Id:          2,
		Name:        "Science",
		Description: "Scientific Discoveries",
		Status: dto.StatusDetail{
			dto.ACTIVE,
			fmt.Sprint(dto.StatusType(dto.ACTIVE)),
		},
	}
	tagStore["Math"] = dto.Tag{
		Id:          3,
		Name:        "Math",
		Description: "Mathematics",
		Status: dto.StatusDetail{
			dto.ACTIVE,
			fmt.Sprint(dto.StatusType(dto.ACTIVE)),
		},
	}
	tagStore["Health"] = dto.Tag{
		Id:          4,
		Name:        "Health",
		Description: "Healthy Living",
		Status: dto.StatusDetail{
			dto.ACTIVE,
			fmt.Sprint(dto.StatusType(dto.ACTIVE)),
		},
	}
}

//func genUUID() (string, error) {
//	uuid := make([]byte, 16)
//	n, err := rand.Read(uuid)
//	if n != len(uuid) || err != nil {
//		return "", err
//	}
//	// TODO: verify the two lines implement RFC 4122 correctly
//	uuid[8] = 0x80 // variant bits see page 5
//	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

//	return hex.EncodeToString(uuid), nil
//}
