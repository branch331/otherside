package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
)

var userID = flag.String("user", "themobil9", "the Spotify user ID to look up")
var applicationPort = "8081"
var baseURL = "http://localhost:" + applicationPort + "/"
var redirectURL = baseURL + "callback"

var SPOTIFY_ID = os.Getenv("SPOTIFY_ID")
var SPOTIFY_SECRET = os.Getenv("SPOTIFY_SECRET")

var SEATGEEK_ID = os.Getenv("SEATGEEK_ID")
var spotifyAuth = spotify.NewAuthenticator(redirectURL, spotify.ScopePlaylistModifyPublic)
var spotifyClient spotify.Client

var baseSeatGeekURL = "https://api.seatgeek.com/2/events?client_id=" + SEATGEEK_ID

type SeatGeekEvent struct {
	title         string
	eventType     string
	url           string
	performers    []string
	localShowtime string
	venueName     string
	venueLocation string
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/authenticate", Authenticate)
	router.HandleFunc("/callback", Callback)
	log.Fatal(http.ListenAndServe(":"+applicationPort, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	findLocalConcerts()
	/*
		GetAndPrintProfileData(spotifyClient)

		GeneratePlayList(spotifyClient, "testing the waters", "a playlist to test the waters. duh!")

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	*/
}

func Callback(w http.ResponseWriter, r *http.Request) {
	token, err := spotifyAuth.Token("dummyState", r)

	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	spotifyClient = spotifyAuth.NewClient(token)

	http.Redirect(w, r, baseURL, http.StatusMovedPermanently)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	dummyState := "dummyState"
	authenticationUrl := ObtainAuthenticationURL(dummyState)
	fmt.Printf("%s", redirectURL)
	http.Redirect(w, r, authenticationUrl, http.StatusMovedPermanently)
}

func ObtainAuthenticationURL(state string) string {
	spotifyAuth.SetAuthInfo(SPOTIFY_ID, SPOTIFY_SECRET)

	return spotifyAuth.AuthURL(state)
}

func GetAndPrintProfileData(client spotify.Client) {
	flag.Parse()

	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
		flag.Usage()
		return
	}

	user, err := client.GetUsersPublicProfile(spotify.ID(*userID))

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("Display name:", user.DisplayName)
	fmt.Println("Spotify URI:", string(user.URI))
	fmt.Println("Endpoint:", user.Endpoint)
	fmt.Println("Followers:", user.Followers.Count)
}

func GeneratePlayList(client spotify.Client, playlistName string, description string) {
	flag.Parse()

	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
		flag.Usage()
		return
	}

	playList, err := client.CreatePlaylistForUser(*userID, playlistName, description, true)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	} else {
		fmt.Printf("Created playlist %s for user %s", playlistName, *userID)
	}

	//var tracksToAdd []spotify.SimpleTrack
	var trackIDsToAdd []spotify.ID

	//Dummy data
	trackIDsToAdd = append(trackIDsToAdd, "6rPO02ozF3bM7NnOV4h6s2") //Despacito
	trackIDsToAdd = append(trackIDsToAdd, "6rPO02ozF3bM7NnOV4h6s2")

	for _, trackID := range trackIDsToAdd {

		_, err := client.AddTracksToPlaylist(playList.ID, trackID)

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		} else {
			fmt.Printf("Added track ID %s to playlist ID %s", trackID, playList.ID)
		}
	}
}

func findLocalConcerts() []SeatGeekEvent {
	testURL := "https://api.seatgeek.com/2/events?client_id=" + SEATGEEK_ID + "&geoip=78745&range=4mi"

	resp, err := http.Get(testURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	} else {
		fmt.Printf("Obtained local events data from SeatGeek.\n")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData) //Fields are empty?
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	eventsFromResponse := responseData["events"].([]interface{})

	seatGeekEvents := make([]SeatGeekEvent, len(eventsFromResponse))

	for i, event := range eventsFromResponse {
		eventData := event.(map[string]interface{})

		seatGeekEvents[i].title = eventData["title"].(string)
		seatGeekEvents[i].eventType = eventData["type"].(string)
		seatGeekEvents[i].localShowtime = eventData["datetime_local"].(string)

		venueData := eventData["venue"].(map[string]interface{})

		seatGeekEvents[i].venueName = venueData["name"].(string)
		seatGeekEvents[i].venueLocation = venueData["display_location"].(string)
		seatGeekEvents[i].url = venueData["url"].(string)

		performersArray := eventData["performers"].([]interface{})

		for _, performer := range performersArray {
			performerData := performer.(map[string]interface{})
			seatGeekEvents[i].performers = append(seatGeekEvents[i].performers, performerData["short_name"].(string))
		}
	}

	return seatGeekEvents
}

func getJson(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
