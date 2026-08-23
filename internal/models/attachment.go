package models

import "fmt"

// Attachment represents a generic resource attached to coursework, materials, or announcements.
type Attachment struct {
	Type  string // "drive", "youtube", "link", "form"
	Title string
	URL   string
}

// Icon returns an emoji icon based on the attachment type.
func (a Attachment) Icon() string {
	switch a.Type {
	case "drive":
		return "📄 [Drive]"
	case "youtube":
		return "🎥 [YouTube]"
	case "link":
		return "🔗 [Link]"
	case "form":
		return "📝 [Form]"
	default:
		return "📎 [File]"
	}
}

func (a Attachment) String() string {
	if a.Title != "" {
		return fmt.Sprintf("%s %s", a.Icon(), a.Title)
	}
	if a.URL != "" {
		return fmt.Sprintf("%s %s", a.Icon(), a.URL)
	}
	return a.Icon()
}

// RawMaterial represents the raw material JSON object returned by Google Classroom API.
type RawMaterial struct {
	DriveFile *struct {
		DriveFile struct {
			ID            string `json:"id"`
			Title         string `json:"title"`
			AlternateLink string `json:"alternateLink"`
			ThumbnailURL  string `json:"thumbnailUrl"`
		} `json:"driveFile"`
		ShareMode string `json:"shareMode"`
	} `json:"driveFile,omitempty"`

	YoutubeVideo *struct {
		ID            string `json:"id"`
		Title         string `json:"title"`
		AlternateLink string `json:"alternateLink"`
		ThumbnailURL  string `json:"thumbnailUrl"`
	} `json:"youtubeVideo,omitempty"`

	Link *struct {
		URL          string `json:"url"`
		Title        string `json:"title"`
		ThumbnailURL string `json:"thumbnailUrl"`
	} `json:"link,omitempty"`

	Form *struct {
		FormURL      string `json:"formUrl"`
		ResponseURL  string `json:"responseUrl"`
		Title        string `json:"title"`
		ThumbnailURL string `json:"thumbnailUrl"`
	} `json:"form,omitempty"`
}

// ToAttachment converts a RawMaterial to a simplified Attachment.
func (r RawMaterial) ToAttachment() (Attachment, bool) {
	if r.DriveFile != nil {
		title := r.DriveFile.DriveFile.Title
		if title == "" {
			title = "Drive Attachment"
		}
		return Attachment{
			Type:  "drive",
			Title: title,
			URL:   r.DriveFile.DriveFile.AlternateLink,
		}, true
	}

	if r.YoutubeVideo != nil {
		title := r.YoutubeVideo.Title
		if title == "" {
			title = "YouTube Video"
		}
		return Attachment{
			Type:  "youtube",
			Title: title,
			URL:   r.YoutubeVideo.AlternateLink,
		}, true
	}

	if r.Link != nil {
		title := r.Link.Title
		if title == "" {
			title = r.Link.URL
		}
		return Attachment{
			Type:  "link",
			Title: title,
			URL:   r.Link.URL,
		}, true
	}

	if r.Form != nil {
		title := r.Form.Title
		if title == "" {
			title = "Google Form"
		}
		url := r.Form.FormURL
		if url == "" {
			url = r.Form.ResponseURL
		}
		return Attachment{
			Type:  "form",
			Title: title,
			URL:   url,
		}, true
	}

	return Attachment{}, false
}

// ConvertRawMaterials converts a slice of RawMaterial to a slice of Attachment.
func ConvertRawMaterials(raws []RawMaterial) []Attachment {
	var atts []Attachment
	for _, r := range raws {
		if att, ok := r.ToAttachment(); ok {
			atts = append(atts, att)
		}
	}
	return atts
}
