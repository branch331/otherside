package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"otherside/api/seatgeekLayer"
	"otherside/api/spotifyLayer"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
)

var applicationPort = "8081"

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/authenticate", Authenticate)
	router.HandleFunc("/callback", Callback)
	router.HandleFunc("/localevents", LocalEvents)
	router.HandleFunc("/toptracks", TopTracks).Methods("POST")
	router.HandleFunc("/buildplaylist", BuildPlaylist).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+applicationPort, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	//Don't serve a view; this will just be an API so keep frontend/backend separate
}

//GET
func LocalEvents(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	postCode, ok := r.URL.Query()["postcode"]
	if !ok || len(postCode[0]) < 1 {
		fmt.Printf("Postcode parameter missing from localevents request.")
	}

	rangeMiles, ok := r.URL.Query()["miles"]
	if !ok || len(rangeMiles[0]) < 1 {
		fmt.Printf("Mile range parameter missing from localevents request.")
	}

	localSeatGeekEvents := seatgeekLayer.FindLocalEvents(postCode[0], rangeMiles[0])

	localSeatGeekEventsJSON, err := json.Marshal(localSeatGeekEvents)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(localSeatGeekEventsJSON)
	}
}

//POST
func TopTracks(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var localSeatGeekEvents []seatgeekLayer.SeatGeekEvent
	var topTracks []spotify.FullTrack

	err := json.NewDecoder(r.Body).Decode(&localSeatGeekEvents)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	for _, event := range localSeatGeekEvents {
		if event.EventType == "concert" || event.EventType == "music_festival" {
			for _, performer := range event.Performers {
				artistID := spotifyLayer.SearchAndFindSpotifyArtistID(performer)
				topTracks = spotifyLayer.GetTopFourSpotifyArtistTracks(artistID)
			}
		}
	}

	topTracksJSON, err := json.Marshal(topTracks)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(topTracksJSON)
	}
}

//POST
func BuildPlaylist(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	playlistName, ok := r.URL.Query()["name"]
	if !ok || len(playlistName[0]) < 1 {
		fmt.Printf("playlistName parameter missing from buildplaylist request.")
	}

	playlistDesc, ok := r.URL.Query()["desc"]
	if !ok || len(playlistDesc[0]) < 1 {
		fmt.Printf("playlistDesc parameter missing from buildplaylist request.")
	}

	var topTracks []spotify.FullTrack

	err := json.NewDecoder(r.Body).Decode(&topTracks)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	playlistID := spotifyLayer.GeneratePlayList(playlistName[0], playlistDesc[0])

	spotifyLayer.AddTracksToPlaylist(playlistID, topTracks)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	spotifyLayer.SetNewSpotifyClient(w, r)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	dummyState := "dummyState"
	authenticationUrl := spotifyLayer.ObtainAuthenticationURL(dummyState)

	http.Redirect(w, r, authenticationUrl, http.StatusMovedPermanently)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
