package models

type CourseResponse struct {
	Courses []CourseModel `json:"courses"`
}

type CourseModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Sub  string `json:"subject"`
}
