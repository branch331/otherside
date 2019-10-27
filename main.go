package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/branch331/otherside/api/seatgeekLayer"
	"github.com/branch331/otherside/api/spotifyLayer"

	"github.com/gorilla/mux"
)

var applicationPort = "8081"

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/authenticate", Authenticate)
	router.HandleFunc("/callback", Callback)
	router.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":"+applicationPort, router))
}

func index(w http.ResponseWriter, r *http.Request) {

}

func test(w http.ResponseWriter, r *http.Request) {
	localSeatGeekEvents := seatgeekLayer.findLocalEvents("78745", "20")

	playlistID := spotifyLayer.GeneratePlayList(spotifyLayer.spotifyClient, "Best playlist2!", "Desc")

	for _, event := range localSeatGeekEvents {
		if event.eventType == "concert" || event.eventType == "music_festival" {
			for _, performer := range event.performers {
				artistID := spotifyLayer.SearchAndFindSpotifyArtistID(spotifyLayer.spotifyClient, performer)
				topTracks := spotifyLayer.GetTopFourSpotifyArtistTracks(spotifyLayer.spotifyClient, artistID)
				spotifyLayer.AddTracksToPlaylist(spotifyLayer.spotifyClient, playlistID, topTracks)
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
	fmt.Printf("%s", spotifyLayer.redirectURL)
	http.Redirect(w, r, authenticationUrl, http.StatusMovedPermanently)
}
