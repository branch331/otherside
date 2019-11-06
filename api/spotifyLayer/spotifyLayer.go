package spotifyLayer

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
)

var spotifyAuth = spotify.NewAuthenticator(redirectURL, spotify.ScopePlaylistModifyPublic)
var spotifyClient spotify.Client

var applicationPort = "8081"
var baseURL = "http://localhost:" + applicationPort + "/"
var redirectURL = baseURL + "callback"

var SPOTIFY_ID = os.Getenv("SPOTIFY_ID")
var SPOTIFY_SECRET = os.Getenv("SPOTIFY_SECRET")

func ObtainAuthenticationURL(state string) string {
	spotifyAuth.SetAuthInfo(SPOTIFY_ID, SPOTIFY_SECRET)

	return spotifyAuth.AuthURL(state)
}

func SetNewSpotifyClient(w http.ResponseWriter, r *http.Request, state string) {
	token, err := spotifyAuth.Token(state, r)

	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	spotifyClient = spotifyAuth.NewClient(token)
}

func GetTopSpotifyArtistTrack(artistID spotify.ID) spotify.FullTrack {
	topTracks, err := spotifyClient.GetArtistsTopTracks(artistID, "US")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return topTracks[0]
	} else {
		return topTracks[0]
	}
}

func GeneratePlayList(playlistName string, description string) spotify.ID {
	flag.Parse()

	currentUser, err := spotifyClient.CurrentUser()
	if err != nil {
		fmt.Printf("Error getting current Spotify user.")
	}

	displayName := currentUser.DisplayName

	playList, err := spotifyClient.CreatePlaylistForUser(displayName, playlistName, description, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return ""
	} else {
		fmt.Printf("Created playlist %s for user %s\n", playlistName, displayName)
	}

	return playList.ID
}

func AddTracksToPlaylist(playlistID spotify.ID, tracksToAdd []spotify.FullTrack) {
	for _, track := range tracksToAdd {
		_, err := spotifyClient.AddTracksToPlaylist(playlistID, track.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
	}
}

func SearchAndFindSpotifyArtistID(artistName string) spotify.ID {
	searchResults, err := spotifyClient.Search(artistName, spotify.SearchTypeArtist)
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
