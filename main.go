package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var userID = flag.String("user", "themobil9", "the Spotify user ID to look up")

//var SPOTIFY_ID = os.Getenv("SPOTIFY_ID")
//var SPOTIFY_SECRET = os.Getenv("SPOTIFY_SECRET")
//var SEATGEEK_ID = os.Getenv("SEATGEEK_ID")

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8081", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	client := SetupSpotifyClient()
	GetAndPrintProfileData(client)
	GeneratePlayList(client, "testing the waters", "a playlist to test the waters. duh!")
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func SetupSpotifyClient() spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	return spotify.Authenticator{}.NewClient(token)
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

	_, err := client.CreatePlaylistForUser(*userID, playlistName, description, true)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	} else {
		fmt.Printf("Created playlist %s for user %s", playlistName, *userID)
	}
}
