package models

type CourseResponse struct {
	Courses       []CourseModel `json:"courses"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
}

type CourseModel struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Section            string `json:"section,omitempty"`
	Sub                string `json:"subject,omitempty"`
	DescriptionHeading string `json:"descriptionHeading,omitempty"`
	Description        string `json:"description,omitempty"`
	Room               string `json:"room,omitempty"`
	CourseState        string `json:"courseState,omitempty"`
	AlternateLink      string `json:"alternateLink,omitempty"`
	EnrollmentCode     string `json:"enrollmentCode,omitempty"`
}
