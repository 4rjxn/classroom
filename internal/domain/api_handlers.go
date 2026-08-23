package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/4rjxn/classroom/internal/models"
)

// ListCourses retrieves all active and provisioned courses for the authenticated user.
func ListCourses(token string) (models.CourseResponse, error) {
	endpoint := "https://classroom.googleapis.com/v1/courses?courseStates=ACTIVE&pageSize=50"
	res, err := DoGetRequest(endpoint, token)
	if err != nil {
		return models.CourseResponse{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return models.CourseResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.CourseResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return models.CourseResponse{}, fmt.Errorf("failed to parse courses JSON: %w", err)
	}

	return response, nil
}

// ListMaterialsInCourse retrieves course work materials for a specific course.
func ListMaterialsInCourse(token string, courseID string) (models.MaterialModel, error) {
	endpoint := fmt.Sprintf("https://classroom.googleapis.com/v1/courses/%s/courseWorkMaterials?pageSize=50", url.PathEscape(courseID))
	res, err := DoGetRequest(endpoint, token)
	if err != nil {
		return models.MaterialModel{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return models.MaterialModel{}, fmt.Errorf("failed to read materials response body: %w", err)
	}

	var materialModel models.MaterialModel
	if err := json.Unmarshal(bodyBytes, &materialModel); err != nil {
		return models.MaterialModel{}, fmt.Errorf("failed to parse materials JSON: %w", err)
	}

	return materialModel, nil
}

// ListAnnouncementsInCourse retrieves announcements for a specific course.
func ListAnnouncementsInCourse(token string, courseID string) (models.AnnouncementsModel, error) {
	endpoint := fmt.Sprintf("https://classroom.googleapis.com/v1/courses/%s/announcements?pageSize=50", url.PathEscape(courseID))
	res, err := DoGetRequest(endpoint, token)
	if err != nil {
		return models.AnnouncementsModel{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return models.AnnouncementsModel{}, fmt.Errorf("failed to read announcements response body: %w", err)
	}

	var announcementsModel models.AnnouncementsModel
	if err := json.Unmarshal(bodyBytes, &announcementsModel); err != nil {
		return models.AnnouncementsModel{}, fmt.Errorf("failed to parse announcements JSON: %w", err)
	}

	return announcementsModel, nil
}

// ListCourseWorkInCourse retrieves coursework (assignments, questions) for a specific course.
func ListCourseWorkInCourse(token string, courseID string) (models.CourseWorkResponse, error) {
	endpoint := fmt.Sprintf("https://classroom.googleapis.com/v1/courses/%s/courseWork?pageSize=50", url.PathEscape(courseID))
	res, err := DoGetRequest(endpoint, token)
	if err != nil {
		return models.CourseWorkResponse{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return models.CourseWorkResponse{}, fmt.Errorf("failed to read coursework response body: %w", err)
	}

	var courseWorkResponse models.CourseWorkResponse
	if err := json.Unmarshal(bodyBytes, &courseWorkResponse); err != nil {
		return models.CourseWorkResponse{}, fmt.Errorf("failed to parse coursework JSON: %w", err)
	}

	return courseWorkResponse, nil
}
