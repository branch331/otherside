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

type SpotifyArtistImage struct {
	Id       spotify.ID
	ImageURL string
}

func ObtainAuthenticationURL(state string) string {
	spotifyAuth.SetAuthInfo(SPOTIFY_ID, SPOTIFY_SECRET)

	return spotifyAuth.AuthURL(state)
}

func SetNewSpotifyClient(w http.ResponseWriter, r *http.Request, state string) error {
	token, err := spotifyAuth.Token(state, r)

	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return err
	}

	spotifyClient = spotifyAuth.NewClient(token)

	return nil
}

func GetTopSpotifyArtistTrack(artistID spotify.ID) (spotify.FullTrack, error) {
	var cachedTopTrack spotify.FullTrack

	topTrackAlreadyCached, err := redisLayer.Exists(string(artistID))
	if err != nil {
		fmt.Printf("Couldn't check if artistID %s exists in Redis: "+err.Error(), string(artistID))
		return cachedTopTrack, err
	}

	if topTrackAlreadyCached {
		redisData, err := redisLayer.GetArtistTopTrack(string(artistID))
		if err != nil {
			fmt.Printf("Couldn't get value for artistID key %s in Redis: "+err.Error(), string(artistID))
			return cachedTopTrack, err
		}

		json.Unmarshal(redisData, &cachedTopTrack)
		if err != nil {
			fmt.Printf("Can't unmarshal value for artistID %s in Redis: "+err.Error(), string(artistID))
			return cachedTopTrack, err
		}

		return cachedTopTrack, nil
	}

	topTracks, err := spotifyClient.GetArtistsTopTracks(artistID, "US")
	if err != nil {
		fmt.Printf("Can't obtain artist %s top tracks: "+err.Error(), string(artistID))
		return cachedTopTrack, err
	} else {
		if len(topTracks) > 0 {
			serializedTopTrack, err := json.Marshal(topTracks[0])
			if err != nil {
				fmt.Printf("Can't marshal artist %s top tracks: "+err.Error(), string(artistID))
			}

			redisLayer.SetArtistTopTrack(string(artistID), serializedTopTrack)
			return topTracks[0], nil
		} else {
			return cachedTopTrack, nil
		}
	}
}

func GeneratePlayList(playlistName string, description string) (spotify.ID, error) {
	flag.Parse()

	currentUser, err := spotifyClient.CurrentUser()
	if err != nil {
		fmt.Printf("Error getting current Spotify user: " + err.Error())
		return "", err
	}

	displayName := currentUser.DisplayName

	playList, err := spotifyClient.CreatePlaylistForUser(displayName, playlistName, description, true)
	if err != nil {
		fmt.Printf("Error creating playlist for user %s: "+err.Error(), displayName)
		return "", err
	} else {
		fmt.Printf("Created playlist %s for user %s\n", playlistName, displayName)
	}

	return playList.ID, nil
}

func AddTracksToPlaylist(playlistID spotify.ID, tracksToAdd []spotify.FullTrack) error {
	for _, track := range tracksToAdd {
		_, err := spotifyClient.AddTracksToPlaylist(playlistID, track.ID)
		if err != nil {
			fmt.Printf("Error adding tracks to playlist: " + err.Error())
			return err
		}
	}

	return nil
}

func SearchAndFindSpotifyArtistID(artistName string) (SpotifyArtistImage, error) {
	var spotifyArtistImage SpotifyArtistImage

	artistNameAlreadyCached, err := redisLayer.Exists(artistName)
	if err != nil {
		fmt.Print("Couldn't access artist name %s from Redis cache: "+err.Error(), artistName)
		return spotifyArtistImage, err
	}

	if artistNameAlreadyCached {
		redisData, err := redisLayer.GetKeyBytes(artistName)

		if err != nil {
			fmt.Printf("Error getting value for artist %s from Redis: "+err.Error(), artistName)
			return spotifyArtistImage, err
		}

		json.Unmarshal(redisData, &spotifyArtistImage)
		if err != nil {
			fmt.Printf("Error unmarshalling value for artist %s from Redis: "+err.Error(), artistName)
		}

		return spotifyArtistImage, nil
	}

	searchResults, err := spotifyClient.Search(artistName, spotify.SearchTypeArtist)
	if err != nil {
		fmt.Printf("Error searching Spotify for artist %s"+err.Error(), artistName)
		return spotifyArtistImage, err
	} else {
		if len(searchResults.Artists.Artists) != 0 {
			artistID := searchResults.Artists.Artists[0].ID
			spotifyArtistImage.Id = artistID
			spotifyArtistImage.ImageURL = searchResults.Artists.Artists[0].Images[0].URL

			spotifyArtistSerialized, err := json.Marshal(spotifyArtistImage)
			if err != nil {
				fmt.Printf("Error marshalling spotify artist: " + err.Error())
			}

			redisLayer.SetKeyBytes(artistName, spotifyArtistSerialized)
			return spotifyArtistImage, nil
		}
	}

	return spotifyArtistImage, nil
}
