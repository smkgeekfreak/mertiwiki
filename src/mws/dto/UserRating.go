// UserRating.go
package dto

type UserRating struct {
	Uid     int64
	Rating  int64
	Updated *JsonTime `json:",omitempty"`
}
