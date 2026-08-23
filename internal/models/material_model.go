package models

type MaterialModel struct {
	Materials     []CourseWorkMaterial `json:"courseWorkMaterial"`
	NextPageToken string               `json:"nextPageToken,omitempty"`
}

type CourseWorkMaterial struct {
	ID            string        `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description,omitempty"`
	State         string        `json:"state,omitempty"`
	AlternateLink string        `json:"alternateLink,omitempty"`
	CreationTime  string        `json:"creationTime,omitempty"`
	UpdateTime    string        `json:"updateTime,omitempty"`
	Materials     []RawMaterial `json:"materials,omitempty"`
}

func (m CourseWorkMaterial) GetAttachments() []Attachment {
	return ConvertRawMaterials(m.Materials)
}
