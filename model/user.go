package model

type User struct {
	Email            string `json:"email" bson:"email"`
	APIDevKey        string `json:"api_dev_key,omitempty" bson:"api_dev_key"`
	AvailableRequest int    `json:"available_request,omitempty" bson:"available_request"`
}
