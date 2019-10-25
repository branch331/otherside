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
	genres        []string //possibly nil
	localShowtime string
	venueName     string
	venueLocation string
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/authenticate", Authenticate)
	router.HandleFunc("/callback", Callback)
	router.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":"+applicationPort, router))
}

func Index(w http.ResponseWriter, r *http.Request) {

}

func test(w http.ResponseWriter, r *http.Request) {
	localSeatGeekEvents := findLocalEvents()

	playlistID := GeneratePlayList(spotifyClient, "Best playlist!", "Desc")

	for _, event := range localSeatGeekEvents {
		if event.eventType == "concert" || event.eventType == "music_festival" {
			for _, performer := range event.performers {
				artistID := SearchAndFindSpotifyArtistID(spotifyClient, performer)
				topTracks := GetTopFourSpotifyArtistTracks(spotifyClient, artistID)
				AddTracksToPlaylist(spotifyClient, playlistID, topTracks)
			}
		}
	}
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

func GetTopFourSpotifyArtistTracks(client spotify.Client, artistID spotify.ID) []spotify.FullTrack {
	topTracks, err := spotifyClient.GetArtistsTopTracks(artistID, "US")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	} else {
		var topFourTracks []spotify.FullTrack

		for i, track := range topTracks {
			topFourTracks = append(topFourTracks, track)
			if i == 3 {
				break
			}
		}
		return topFourTracks
	}
}

func GeneratePlayList(client spotify.Client, playlistName string, description string) spotify.ID {
	flag.Parse()

	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
		flag.Usage()
		return ""
	}

	playList, err := client.CreatePlaylistForUser(*userID, playlistName, description, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return ""
	} else {
		fmt.Printf("Created playlist %s for user %s", playlistName, *userID)
	}

	return playList.ID
}

func AddTracksToPlaylist(client spotify.Client, playlistID spotify.ID, tracksToAdd []spotify.FullTrack) {
	for _, track := range tracksToAdd {
		_, err := client.AddTracksToPlaylist(playlistID, track.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
	}
}

func SearchAndFindSpotifyArtistID(client spotify.Client, artistName string) spotify.ID {
	searchResults, err := client.Search(artistName, spotify.SearchTypeArtist)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return ""
	} else {
		if len(searchResults.Artists.Artists) != 0 {
			return searchResults.Artists.Artists[0].ID
		}
	}

	return ""
}

func findLocalEvents() []SeatGeekEvent {
	testURL := "https://api.seatgeek.com/2/events?client_id=" + SEATGEEK_ID + "&geoip=78745&range=20mi"

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

	err = json.Unmarshal(body, &responseData)
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
			seatGeekEvents[i].genres = GetSeatGeekArtistGenres(fmt.Sprintf("%d", int(performerData["id"].(float64))))
		}
	}

	return seatGeekEvents
}

func GetSeatGeekArtistGenres(performerId string) []string {
	SeatGeekPerformerURL := "https://api.seatgeek.com/2/performers/" + performerId + "?client_id=" + SEATGEEK_ID
	fmt.Printf(SeatGeekPerformerURL)

	resp, err := http.Get(SeatGeekPerformerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	var genreArray []string

	if genresFromResponse, keyExists := responseData["genres"].([]interface{}); keyExists {
		for _, genre := range genresFromResponse {
			genreData := genre.(map[string]interface{})
			genreArray = append(genreArray, genreData["slug"].(string))
		}
	}

	return genreArray
}
