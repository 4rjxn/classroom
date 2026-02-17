package domain

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/classroom-cli/internal/models"
)

func ListCourses(token string) {
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
	var response models.Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		fmt.Println("parse err")
	}
	for _, course := range response.Courses {
		fmt.Printf("ID: %s, Name: %s, Subject: %s\n", course.Id, course.Name, course.Sub)
	}
}

func ListMaterialsInCourse(token string, courseId string) {
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
	json.Unmarshal(bodyBytes, &materialModel)
	for _, cw := range materialModel.CourseWorkMaterial {
		fmt.Println("CourseWork ID:", cw.ID)
		fmt.Println("Title:", cw.Title)

		for _, mat := range cw.Materials {
			file := mat.DriveFile.DriveFile

			fmt.Println("  File ID:", file.ID)
			fmt.Println("  File Title:", file.Title)
			fmt.Println("  Link:", file.AlternateLink)
			fmt.Println("  Share Mode:", mat.DriveFile.ShareMode)
			fmt.Println()
		}
	}

}

func ListAnnouncementsInCourse(token string, courseId string) {
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
	for _, ann := range announcementsModel.Announcements {
		fmt.Println("Text: ", ann.Text)
		for _, material := range ann.Materials {
			file := material.DriveFile.DriveFile
			fmt.Println("  Title: ", file.Title)
			fmt.Println("  Link: ", file.AlternateLink)
		}
	}

}
