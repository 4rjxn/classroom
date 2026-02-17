package models

type Response struct {
	Courses []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Sub  string `json:"subject"`
	} `json:"courses"`
}
