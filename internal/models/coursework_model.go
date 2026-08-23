package models

import "fmt"

type CourseWorkResponse struct {
	CourseWork    []CourseWorkModel `json:"courseWork"`
	NextPageToken string            `json:"nextPageToken,omitempty"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

func (d Date) IsZero() bool {
	return d.Year == 0 && d.Month == 0 && d.Day == 0
}

func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

type TimeOfDay struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

func (t TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hours, t.Minutes)
}

type CourseWorkModel struct {
	ID            string        `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description,omitempty"`
	State         string        `json:"state,omitempty"`
	AlternateLink string        `json:"alternateLink,omitempty"`
	CreationTime  string        `json:"creationTime,omitempty"`
	UpdateTime    string        `json:"updateTime,omitempty"`
	DueDate       Date          `json:"dueDate,omitempty"`
	DueTime       TimeOfDay     `json:"dueTime,omitempty"`
	MaxPoints     float64       `json:"maxPoints,omitempty"`
	WorkType      string        `json:"workType,omitempty"`
	Materials     []RawMaterial `json:"materials,omitempty"`
}

func (c CourseWorkModel) GetAttachments() []Attachment {
	return ConvertRawMaterials(c.Materials)
}

func (c CourseWorkModel) DueString() string {
	if c.DueDate.IsZero() {
		return "No due date"
	}
	if c.DueTime.Hours != 0 || c.DueTime.Minutes != 0 {
		return fmt.Sprintf("%s %s", c.DueDate.String(), c.DueTime.String())
	}
	return c.DueDate.String()
}
