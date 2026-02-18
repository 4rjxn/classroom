package models

type MaterialModel struct {
	Materials []CourseWorkMaterial `json:"courseWorkMaterial"`
}

type CourseWorkMaterial struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Materials   []struct {
		DriveFile struct {
			DriveFile struct {
				ID            string `json:"id"`
				Title         string `json:"title"`
				AlternateLink string `json:"alternateLink"`
			} `json:"driveFile"`
			ShareMode string `json:"shareMode"`
		} `json:"driveFile"`
	} `json:"materials"`
	AlternateLink string `json:"alternateLink"`
}
