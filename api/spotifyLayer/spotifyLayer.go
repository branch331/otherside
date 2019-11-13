package spotifyLayer

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"otherside/api/redisLayer"

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
	artistIDAlreadyCached, err := redisLayer.Exists(string(artistID))
	if err != nil {
		fmt.Printf("Couldn't check if artistID %s exists in Redis: "+err.Error(), string(artistID))
	}

	if artistIDAlreadyCached {
		redisData, err := redisLayer.GetArtistTopTrack(string(artistID))
		if err != nil {
			fmt.Printf("Couldn't get value for artistID key %s in Redis: "+err.Error(), string(artistID))
		}

		var cachedTopTrack spotify.FullTrack
		json.Unmarshal(redisData, &cachedTopTrack)
		if err != nil {
			fmt.Printf("Can't unmarshal value for artistID %s in Redis: "+err.Error(), string(artistID))
		}

		return cachedTopTrack
	}

	topTracks, err := spotifyClient.GetArtistsTopTracks(artistID, "US")
	if err != nil {
		fmt.Printf("Can't obtain artist %s top tracks: "+err.Error(), string(artistID))
		return topTracks[0]
	} else {
		serializedTopTrack, err := json.Marshal(topTracks[0])
		if err != nil {
			fmt.Printf("Can't marshal artist %s top tracks: "+err.Error(), string(artistID))
		}

		redisLayer.SetArtistTopTrack(string(artistID), serializedTopTrack)
		return topTracks[0]
	}
}

func GeneratePlayList(playlistName string, description string) spotify.ID {
	flag.Parse()

	currentUser, err := spotifyClient.CurrentUser()
	if err != nil {
		fmt.Printf("Error getting current Spotify user: " + err.Error())
	}

	displayName := currentUser.DisplayName

	playList, err := spotifyClient.CreatePlaylistForUser(displayName, playlistName, description, true)
	if err != nil {
		fmt.Printf("Error creating playlist for user %s: "+err.Error(), displayName)
	} else {
		fmt.Printf("Created playlist %s for user %s\n", playlistName, displayName)
	}

	return playList.ID
}

func AddTracksToPlaylist(playlistID spotify.ID, tracksToAdd []spotify.FullTrack) {
	for _, track := range tracksToAdd {
		_, err := spotifyClient.AddTracksToPlaylist(playlistID, track.ID)
		if err != nil {
			fmt.Printf("Error adding tracks to playlist: " + err.Error())
			return
		}
	}
}

func SearchAndFindSpotifyArtistID(artistName string) spotify.ID {
	artistNameAlreadyCached, err := redisLayer.Exists(artistName)
	if err != nil {
		fmt.Print("Couldn't access artist name %s from Redis cache: "+err.Error(), artistName)
	}

	if artistNameAlreadyCached {
		artistID, err := redisLayer.GetKeyString(artistName)

		if err != nil {
			fmt.Printf("Error getting value for artistID %s from Redis: "+err.Error(), artistID)
		}

		return spotify.ID(artistID)
	}

	searchResults, err := spotifyClient.Search(artistName, spotify.SearchTypeArtist)
	if err != nil {
		fmt.Printf("Error searching Spotify for artist %s"+err.Error(), artistName)
		return ""
	} else {
		if len(searchResults.Artists.Artists) != 0 {
			artistID := searchResults.Artists.Artists[0].ID
			redisLayer.SetKeyString(artistName, string(artistID))
			return artistID
		}
	}

	return ""
}
