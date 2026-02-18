package domain

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/classroom-cli/internal/models"
)

func ListCourses(token string) models.CourseResponse {
	baseUrl := "https://classroom.googleapis.com/v1/courses"
	res, err := DoGetRequest(baseUrl, token)
	if err != nil {
		fmt.Println("get request err")
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("body read err")
	}
	var response models.CourseResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		fmt.Println("parse err")
	}
	return response
}

func ListMaterialsInCourse(token string, courseId string) models.MaterialModel {
	baseUrl := fmt.Sprintf("https://classroom.googleapis.com/v1/courses/%s/courseWorkMaterials", courseId)
	res, err := DoGetRequest(baseUrl, token)
	if err != nil {
		fmt.Println("get request err")
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("body read err")
	}
	materialModel := models.MaterialModel{}
	err = json.Unmarshal(bodyBytes, &materialModel)
	if err != nil {
		fmt.Println("parse error")
	}
	return materialModel

}

func ListAnnouncementsInCourse(token string, courseId string) models.AnnouncementsModel {
	baseUrl := fmt.Sprintf("https://classroom.googleapis.com/v1/courses/%s/announcements", courseId)
	res, err := DoGetRequest(baseUrl, token)
	if err != nil {
		fmt.Println("get request err")
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("body read err")
	}
	announcementsModel := models.AnnouncementsModel{}
	json.Unmarshal(bodyBytes, &announcementsModel)
	return announcementsModel

}
