package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"otherside/api/seatgeekLayer"
	"otherside/api/spotifyLayer"
	"reflect"
	"strings"
	"time"

	"otherside/api/redisLayer"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
)

type TopTrackResponse struct {
	Track        spotify.FullTrack
	ArtistExists bool
	err          error
}

type ArtistIDResponse struct {
	ID  spotify.ID
	err error
}

var applicationPort = "8081"

var cityPostcodeMap map[string]string
var availableGenres []string

func main() {
	cityPostcodeMap = generateCityPostcodeMap()
	availableGenres = generateGenres()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/cities", Cities)
	router.HandleFunc("/genres", Genres)
	router.HandleFunc("/authenticate", Authenticate)
	router.HandleFunc("/callback", Callback)
	router.HandleFunc("/localevents", LocalEvents)
	router.HandleFunc("/toptracks", TopTracks).Methods("POST")
	router.HandleFunc("/artistids", ArtistIDs).Methods("POST")
	router.HandleFunc("/buildplaylist", BuildPlaylist).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+applicationPort, router))
}

func generateCityPostcodeMap() map[string]string {
	postCodeMap := map[string]string{
		"Austin TX":        "78759",
		"Atlanta GA":       "30301",
		"Washington DC":    "20001",
		"Nashville TN":     "37011",
		"Las Vegas NV":     "88901",
		"New Haven CT":     "06501",
		"Buffalo NY":       "14201",
		"Troy NY":          "12180",
		"Kansas City MO":   "64030",
		"Tulsa OK":         "74008",
		"Denver CO":        "80014",
		"Omaha NE":         "68007",
		"San Diego CA":     "91945",
		"Boston MA":        "02101",
		"Indianapolis IN":  "46077",
		"Pittsburgh PA":    "15106",
		"St Louis MO":      "63101",
		"New Orleans LA":   "70032",
		"Detroit MI":       "48127",
		"Louisville KY":    "40018",
		"San Francisco CA": "94016",
		"Norfolk VA":       "23324",
		"Cincinatti OH":    "45203",
		"Birmingham AL":    "35005",
		"Charlotte NC":     "28105",
		"Des Moines IA":    "50047",
		"Philadelphia PA":  "19093",
		"Chicago IL":       "60007",
		"Houston TX":       "77001",
		"Dallas TX":        "75043",
	}

	return postCodeMap
}

func generateGenres() []string {
	genres := []string{
		"rock",
		"hard-rock",
		"indie",
		"hip-hop",
		"jazz",
		"pop",
		"soul",
		"rnb",
		"alternative",
		"classic-rock",
		"country",
		"folk",
		"punk",
		"electronic",
		"blues",
		"techno",
		"rap",
		"latin",
		"classical",
	}

	return genres
}

//GET
func Cities(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var cities []string
	for key := range cityPostcodeMap {
		cities = append(cities, key)
	}

	citiesJSON, err := json.Marshal(cities)
	if err != nil {
		fmt.Printf("Error Marshalling city keys: " + err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(citiesJSON)
	}
}

//GET
func Genres(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	genresJSON, err := json.Marshal(availableGenres)
	if err != nil {
		fmt.Printf("Error Marshalling city keys: " + err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(genresJSON)
	}
}

//GET
func LocalEvents(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	cities, ok := r.URL.Query()["cities"]
	if !ok || len(cities[0]) < 1 {
		fmt.Printf("cities parameter missing from localevents request.")
	}

	genres, ok := r.URL.Query()["genres"]
	if !ok || len(genres[0]) < 1 {
		fmt.Printf("genres parameter missing from localevents request.")
	}

	citiesArray := QueryStringToArray(cities[0])
	genreArray := QueryStringToArray(genres[0])

	var postCodeArray []string

	for _, val := range citiesArray {
		postCodeArray = append(postCodeArray, cityPostcodeMap[val])
	}

	localSeatGeekEvents := seatgeekLayer.FindLocalEvents(postCodeArray, genreArray)

	localSeatGeekEventsJSON, err := json.Marshal(localSeatGeekEvents)
	if err != nil {
		fmt.Printf("Error Marshalling localseatgeekevents data: " + err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(localSeatGeekEventsJSON)
	}
}

func QueryStringToArray(queryString string) []string {
	testsArray := strings.Split(strings.Trim(queryString, "[]"), ",")
	return testsArray
}

//POST
func ArtistIDs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	time.Sleep(200 * time.Millisecond)

	var localSeatGeekEvents []seatgeekLayer.SeatGeekEvent
	var artistIDs []spotify.ID

	err := json.NewDecoder(r.Body).Decode(&localSeatGeekEvents)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding local seatgeek events: " + err.Error()))
	}

	var artistChannels []chan ArtistIDResponse

	t4 := time.Now()
	for _, event := range localSeatGeekEvents {
		for _, performer := range event.Performers {
			artistChan := make(chan ArtistIDResponse)
			artistChannels = append(artistChannels, artistChan)
			go GetArtistID(performer, artistChan)
		}
	}
	cases := make([]reflect.SelectCase, len(artistChannels))

	for i, artistChan := range artistChannels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(artistChan)}
	}
	remainingCases := len(cases)

	for remainingCases > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			//Channel has been closed; zero out channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remainingCases--
			continue
		}
		response := value.Interface().(ArtistIDResponse)
		if response.err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error getting artist ID for spotify artist: " + err.Error() + "\n"))
		} else {
			artistIDs = append(artistIDs, response.ID)
		}
	}
	fmt.Println("[Time benchmark] Artist IDs " + time.Since(t4).String())
	artistIDJSON, err := json.Marshal(artistIDs)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshaling spotify artist IDs")))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(artistIDJSON)
	}
}

//POST
func TopTracks(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	time.Sleep(200 * time.Millisecond)

	var artistIDs []spotify.ID
	var topTracks []spotify.FullTrack

	err := json.NewDecoder(r.Body).Decode(&artistIDs)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding Spotify Artist IDs: " + err.Error() + "\n"))
	}

	var topTrackChannels []chan TopTrackResponse

	t4 := time.Now()
	for _, ID := range artistIDs {
		topTrackChan := make(chan TopTrackResponse)
		topTrackChannels = append(topTrackChannels, topTrackChan)
		go GetArtistTopTrack(ID, topTrackChan)
	}
	cases := make([]reflect.SelectCase, len(topTrackChannels))

	for i, topTrackChan := range topTrackChannels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(topTrackChan)}
	}
	remainingCases := len(cases)

	for remainingCases > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			//Channel has been closed; zero out channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remainingCases--
			continue
		}
		response := value.Interface().(TopTrackResponse)
		if response.err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error getting top track for spotify artist: " + err.Error()))
		} else {
			if response.ArtistExists {
				topTracks = append(topTracks, response.Track)
			}
		}
	}
	fmt.Println("[Time benchmark] Top tracks " + time.Since(t4).String())
	topTracksJSON, err := json.Marshal(topTracks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshaling spotify top tracks")))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(topTracksJSON)
	}
}

func GetArtistID(performer string, artistChan chan<- ArtistIDResponse) {
	var response ArtistIDResponse

	artistID, err := spotifyLayer.SearchAndFindSpotifyArtistID(performer)
	if err != nil {
		response.err = err
		artistChan <- response
		close(artistChan)
	} else {
		response.ID = artistID
		response.err = err
		artistChan <- response
		close(artistChan)
	}
}

func GetArtistTopTrack(artistID spotify.ID, topTrackChan chan<- TopTrackResponse) {
	var response TopTrackResponse

	if artistID != "" {
		topArtistTrack, err := spotifyLayer.GetTopSpotifyArtistTrack(artistID)
		response.Track = topArtistTrack
		response.ArtistExists = true
		response.err = err
		topTrackChan <- response
		close(topTrackChan)
	} else {
		response.ArtistExists = false
		response.err = nil
		topTrackChan <- response
		close(topTrackChan)
	}
}

//POST
func BuildPlaylist(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	playlistName, ok := r.URL.Query()["name"]
	if !ok || len(playlistName[0]) < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Playlist 'name' parameter missing from buildplaylist request."))
	}

	playlistDesc, ok := r.URL.Query()["desc"]
	if !ok || len(playlistDesc[0]) < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Playlist 'desc' parameter missing from buildplaylist request."))
	}

	var topTracks []spotify.FullTrack

	err := json.NewDecoder(r.Body).Decode(&topTracks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding Spotify Top Tracks: " + err.Error()))
	}

	playlistID, err := spotifyLayer.GeneratePlayList(playlistName[0], playlistDesc[0])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating playlist: " + err.Error()))
	}

	err = spotifyLayer.AddTracksToPlaylist(playlistID, topTracks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error adding tracks to playlist: " + err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Playlist generated."))
	}
}

//POST
func ConfigureCallbackURL(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	redirectURL, ok := r.URL.Query()["redirectURL"]
	if !ok || len(redirectURL[0]) < 1 {
		fmt.Printf("redirectURL parameter missing from configuration request.")
	}
}

//GET
//Callback is called from the Spotify authentication flow, and redirects to <Host>/#/callback
func Callback(w http.ResponseWriter, r *http.Request) {
	state, ok := r.URL.Query()["state"]
	if !ok || len(state[0]) < 1 {
		fmt.Printf("State parameter missing from callback.")
	}

	spotifyLayer.SetNewSpotifyClient(w, r, state[0])

	clientOrigin, err := redisLayer.GetKeyString(state[0])
	if err != nil {
		fmt.Printf("Error getting state/origin Redis key: " + err.Error())
	}

	redirectURL := clientOrigin + "/#/callback?state=" + state[0]

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

//GET
func Authenticate(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	state, ok := r.URL.Query()["state"]
	if !ok || len(state[0]) < 1 {
		fmt.Printf("State parameter missing from authenticate request.")
	}

	err := redisLayer.SetKeyString(state[0], r.Header.Get("Origin"))
	if err != nil {
		fmt.Printf("Error setting state/origin Redis key: " + err.Error())
	}

	authenticationUrl := spotifyLayer.ObtainAuthenticationURL(state[0])
	fmt.Fprint(w, authenticationUrl)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
