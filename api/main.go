package main

import (
	"flag"
	"fmt"
	"html"
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

//var SEATGEEK_ID = os.Getenv("SEATGEEK_ID")
var spotifyAuth = spotify.NewAuthenticator(redirectURL, spotify.ScopePlaylistModifyPublic)
var spotifyClient spotify.Client

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/authenticate", Authenticate)
	router.HandleFunc("/callback", Callback)
	log.Fatal(http.ListenAndServe(":"+applicationPort, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	GetAndPrintProfileData(spotifyClient)

	GeneratePlayList(spotifyClient, "testing the waters", "a playlist to test the waters. duh!")

	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
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

	_, err := client.CreatePlaylistForUser(*userID, playlistName, description, true)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	} else {
		fmt.Printf("Created playlist %s for user %s", playlistName, *userID)
	}
}
