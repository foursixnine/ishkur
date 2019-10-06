package main

import (

    "io/ioutil"
	"fmt"
	"github.com/foursixnine/strava/oauth2"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"regexp"
)

var (
	strava_Oauth_config *oauth2.Config
)

const port = 8080

func init() {
	strava_Oauth_config = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/login/callback",
		ClientID:     os.Getenv("CID"),
		ClientSecret: os.Getenv("STRAVA_TOKEN"),
		Scopes:       []string{"activity:read_all", "activity:write"},
		Endpoint:     strava.Endpoint,
	}
}

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

func handle_strava_login(w http.ResponseWriter, r *http.Request) {
	url := strava_Oauth_config.AuthCodeURL(oauthStateString)
	fmt.Printf("%s\n", url)
	re := regexp.MustCompile(`([\+])`)
	url = re.ReplaceAllString(url, ",")
	fmt.Printf("%s\n", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handle_strava_callback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Content: %s\n", r.Body)
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := strava_Oauth_config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.strava.com/api/v3/athlete?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

func main() {
	defer fmt.Println("CID, and possible other variables are not set, source production.env. Bye bye")
	if os.Getenv("CID") == "" {
		return
	}
	//http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handle_strava_login)
	http.HandleFunc("/login/callback", handle_strava_callback)
	fmt.Printf("Visit http://localhost:%d/ to view the demo\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
