package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/classroom-cli/internal/models"
	"github.com/pkg/browser"
)

func GenerateToken(config models.Config) string {
	type AuthResponce struct {
		AccessToken string `json:"access_token"`
	}
	baseUrl := "https://accounts.google.com/o/oauth2/v2/auth"
	tokenUrl := "https://oauth2.googleapis.com/token"
	redirectUri := "http://localhost:4321"
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
