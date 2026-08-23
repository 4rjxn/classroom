package models

type AnnouncementsModel struct {
	Announcements []CourseAnnouncement `json:"announcements"`
	NextPageToken string               `json:"nextPageToken,omitempty"`
}

type CourseAnnouncement struct {
	CourseID      string        `json:"courseId"`
	ID            string        `json:"id"`
	Text          string        `json:"text"`
	State         string        `json:"state,omitempty"`
	AlternateLink string        `json:"alternateLink,omitempty"`
	CreationTime  string        `json:"creationTime,omitempty"`
	UpdateTime    string        `json:"updateTime,omitempty"`
	CreatorUserID string        `json:"creatorUserId,omitempty"`
	Materials     []RawMaterial `json:"materials,omitempty"`
}

func (a CourseAnnouncement) GetAttachments() []Attachment {
	return ConvertRawMaterials(a.Materials)
}
