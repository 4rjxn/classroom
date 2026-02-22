package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/classroom-cli/internal/models"
	"github.com/classroom-cli/internal/utils"
)

var tokenUrl string = "https://oauth2.googleapis.com/token"
var redirectUri string = "http://localhost:4321"

type AuthResponce struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func OffileGeneration(config models.Config) string {
	refreshToken := utils.ReadRefreshToken()
	if refreshToken != "" {
		res, err := http.PostForm(tokenUrl, url.Values{
			"refresh_token": {refreshToken},
			"client_id":     {config.ClientId},
			"client_secret": {config.ClientSecret},
			"redirect_uri":  {redirectUri},
			"grant_type":    {"refresh_token"},
		})
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		aut := AuthResponce{}
		json.Unmarshal(body, &aut)
		return aut.AccessToken

	}
	return GenerateToken(config)
}
func GenerateToken(config models.Config) string {
	baseUrl := "https://accounts.google.com/o/oauth2/v2/auth"
	scope := "https://www.googleapis.com/auth/classroom.courses.readonly https://www.googleapis.com/auth/classroom.courseworkmaterials.readonly https://www.googleapis.com/auth/classroom.announcements.readonly"
	code := ""
	authServerUrl, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}
	clientId := config.ClientId
	clientSecret := config.ClientSecret

	query := authServerUrl.Query()
	query.Add("client_id", clientId)
	query.Add("redirect_uri", redirectUri)
	query.Add("response_type", "code")
	query.Add("access_type", "offline")
	query.Add("scope", scope)
	query.Add("prompt", "consent")

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

	exec.Command("xdg-open", authServerUrl.String()).Start()

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
	utils.StoreRefreshToken(body)
	tokenData := AuthResponce{}
	json.Unmarshal(body, &tokenData)
	return tokenData.AccessToken

}
