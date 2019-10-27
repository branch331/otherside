package main

import (
	"log"
	"net/http"
	"otherside/api/seatgeekLayer"
	"otherside/api/spotifyLayer"

	"github.com/gorilla/mux"
)

var applicationPort = "8081"

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/authenticate", Authenticate)
	router.HandleFunc("/callback", Callback)
	router.HandleFunc("/test", Test)
	log.Fatal(http.ListenAndServe(":"+applicationPort, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}

func Test(w http.ResponseWriter, r *http.Request) {
	localSeatGeekEvents := seatgeekLayer.FindLocalEvents("78745", "20")

	playlistID := spotifyLayer.GeneratePlayList("Best playlist2!", "Desc")

	for _, event := range localSeatGeekEvents {
		if event.EventType == "concert" || event.EventType == "music_festival" {
			for _, performer := range event.Performers {
				artistID := spotifyLayer.SearchAndFindSpotifyArtistID(performer)
				topTracks := spotifyLayer.GetTopFourSpotifyArtistTracks(artistID)
				spotifyLayer.AddTracksToPlaylist(playlistID, topTracks)
			}
		}
	}
}

func Callback(w http.ResponseWriter, r *http.Request) {
	spotifyLayer.SetNewSpotifyClient(w, r)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	dummyState := "dummyState"
	authenticationUrl := spotifyLayer.ObtainAuthenticationURL(dummyState)

	http.Redirect(w, r, authenticationUrl, http.StatusMovedPermanently)
}
