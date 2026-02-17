package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/browser"
)

const courseListUrl string = "https://classroom.googleapis.com/v1/courses"

func generateToken() string {
	type AuthResponce struct {
		AccessToken string `json:"access_token"`
	}
	baseUrl := "https://accounts.google.com/o/oauth2/v2/auth"
	tokenUrl := "https://oauth2.googleapis.com/token"
	redirectUri := "http://localhost:4321"
	scope := "https://www.googleapis.com/auth/classroom.courses.readonly https://www.googleapis.com/auth/classroom.courseworkmaterials.readonly"
	code := ""
	authServerUrl, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SC")

	query := authServerUrl.Query()
	query.Add("client_id", clientId)
	query.Add("redirect_uri", redirectUri)
	query.Add("response_type", "code")
	query.Add("scope", scope)

	authServerUrl.RawQuery = query.Encode()
	channel := make(chan bool)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		http.Redirect(w, r, "https://google.com", http.StatusMovedPermanently)
		channel <- true

	})
	server := &http.Server{Addr: ":4321"}
	go func() {
		close := <-channel
		if close {
			if err := server.Shutdown(context.Background()); err != nil {
				fmt.Println("failed to shutdown")
			}
		}
	}()

	browser.OpenURL(authServerUrl.String())
	server.ListenAndServe()

	res, err := http.PostForm(tokenUrl, url.Values{
		"code":          {code},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectUri},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	tokenData := AuthResponce{}
	json.Unmarshal(body, &tokenData)
	return tokenData.AccessToken

}

func doGetRequest(url, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

type CourseModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Sub  string `json:"subject"`
}
type Response struct {
	Courses []CourseModel `json:"courses"`
}

func listCourses(token string) {
	res, err := doGetRequest(courseListUrl, token)
	if err != nil {
		fmt.Println("get request err")
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("body read err")
	}
	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		fmt.Println("parse err")
	}
	for _, course := range response.Courses {
		fmt.Printf("ID: %s, Name: %s, Subject: %s\n", course.Id, course.Name, course.Sub)
	}
}

func listMaterialsInCourse(token string, courseId string) {
	baseUrl := fmt.Sprintf("https://classroom.googleapis.com/v1/courses/%s/courseWorkMaterials", courseId)
	res, err := doGetRequest(baseUrl, token)
	if err != nil {
		fmt.Println("get request err")
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("body read err")
	}
	materialModel := MaterialModel{}
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

func main() {
	token := generateToken()
	var action string
	var param string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("== hello there ==\n")
	for {
		fmt.Print("Available actions:\n")
		fmt.Print(" list courses [list]\n")
		fmt.Print(" list materials [choose <id>]\n")
		fmt.Print(" quit [q]\n")
		fmt.Print(">> ")
		scanner.Scan()
		parts := strings.Fields(scanner.Text())
		if len(parts) > 0 {
			action = parts[0]
		}
		if len(parts) > 1 {
			param = parts[1]
		}
		if action == "q" {
			return
		}
		if action == "list" {
			fmt.Println("======================================")
			listCourses(token)
			fmt.Println("======================================")
		} else if action == "choose" && param != "" {
			fmt.Println("======================================")
			listMaterialsInCourse(token, param)
			fmt.Println("======================================")
		}
	}
}
