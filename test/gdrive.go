package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var tokFile = "token.json"

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from an enviroment variable.
func tokenFromEnv(param string) (*oauth2.Token, error) {
	t := os.Getenv(param)
	tok := &oauth2.Token{}
	err := json.NewDecoder(strings.NewReader(t)).Decode(tok)
	return tok, err
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	b := os.Getenv("GDRIVE_CREDENTIALS")

	// Check if a token file exists if not load from enviroment
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		// Load token from enviroment variables
		tok, err = tokenFromEnv("GDRIVE_TOKEN")
		if err == nil {
			saveToken(tokFile, tok)
		}
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(b), drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	buildNumber := os.Getenv("TRAVIS_BUILD_NUMBER")
	osName := os.Getenv("TRAVIS_OS_NAME")

	remoteFileName := buildNumber + "/" + osName + ".out"
	localFileName := osName + ".out"
	baseMimeType := "text/plain"

	file, err := os.Open(localFileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer file.Close()
	f := &drive.File{
		Name:     remoteFileName,
		MimeType: baseMimeType,
	}
	res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("%s, %s, %s\n", res.Name, res.Id, res.MimeType)
}
