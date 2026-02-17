package models

type MaterialModel struct {
	CourseWorkMaterial []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Materials []struct {
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
	} `json:"courseWorkMaterial"`
}
