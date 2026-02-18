package models

type AnnouncementsModel struct {
	Announcements []CourseAnnouncement `json:"announcements"`
}

type CourseAnnouncement struct {
	CourseID  string `json:"courseId"`
	ID        string `json:"id"`
	Text      string `json:"text"`
	Materials []struct {
		DriveFile struct {
			DriveFile struct {
				ID            string `json:"id"`
				Title         string `json:"title"`
				AlternateLink string `json:"alternateLink"`
			} `json:"driveFile"`
		} `json:"driveFile"`
	} `json:"materials"`
	AlternateLink string `json:"alternateLink"`
}
