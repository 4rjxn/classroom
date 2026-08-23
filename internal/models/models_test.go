package models_test

import (
	"encoding/json"
	"testing"

	"github.com/classroom-cli/internal/models"
)

func TestDateAndDueString(t *testing.T) {
	d := models.Date{Year: 2026, Month: 8, Day: 24}
	if d.String() != "2026-08-24" {
		t.Errorf("expected 2026-08-24, got %s", d.String())
	}

	cw := models.CourseWorkModel{
		Title:   "Assignment 1",
		DueDate: d,
		DueTime: models.TimeOfDay{Hours: 23, Minutes: 59},
	}
	if cw.DueString() != "2026-08-24 23:59" {
		t.Errorf("expected '2026-08-24 23:59', got '%s'", cw.DueString())
	}

	cwEmpty := models.CourseWorkModel{Title: "No Due Date"}
	if cwEmpty.DueString() != "No due date" {
		t.Errorf("expected 'No due date', got '%s'", cwEmpty.DueString())
	}
}

func TestRawMaterialConversion(t *testing.T) {
	jsonData := `[
		{
			"driveFile": {
				"driveFile": {
					"id": "file-123",
					"title": "Lecture Notes.pdf",
					"alternateLink": "https://drive.google.com/file/d/123/view"
				}
			}
		},
		{
			"youtubeVideo": {
				"id": "yt-456",
				"title": "Course Intro",
				"alternateLink": "https://youtube.com/watch?v=456"
			}
		},
		{
			"link": {
				"title": "Documentation",
				"url": "https://golang.org"
			}
		},
		{
			"form": {
				"title": "Feedback Survey",
				"formUrl": "https://forms.google.com/123"
			}
		}
	]`

	var raw []models.RawMaterial
	if err := json.Unmarshal([]byte(jsonData), &raw); err != nil {
		t.Fatalf("failed to unmarshal raw materials: %v", err)
	}

	atts := models.ConvertRawMaterials(raw)
	if len(atts) != 4 {
		t.Fatalf("expected 4 attachments, got %d", len(atts))
	}

	if atts[0].Type != "drive" || atts[0].Title != "Lecture Notes.pdf" {
		t.Errorf("unexpected drive attachment: %+v", atts[0])
	}
	if atts[1].Type != "youtube" || atts[1].Title != "Course Intro" {
		t.Errorf("unexpected youtube attachment: %+v", atts[1])
	}
	if atts[2].Type != "link" || atts[2].URL != "https://golang.org" {
		t.Errorf("unexpected link attachment: %+v", atts[2])
	}
	if atts[3].Type != "form" || atts[3].URL != "https://forms.google.com/123" {
		t.Errorf("unexpected form attachment: %+v", atts[3])
	}
}
