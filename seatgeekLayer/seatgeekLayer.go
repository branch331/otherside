package seatgeekLayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

//SeatGeekEvent is a struct to handle pertinent SeatGeek response data.
type SeatGeekEvent struct {
	Title         string
	EventType     string
	URL           string
	Performers    []string
	Genres        []string //possibly nil
	LocalShowtime string
	VenueName     string
	VenueLocation string
}

var SEATGEEK_ID = os.Getenv("SEATGEEK_ID")
var requestExpirationTime time.Time

//FindLocalEvents makes a request to the SeatGeek Events API using the postal code and range,
//and returns an array of SeatGeekEvents.
func FindLocalEvents(postalCode string, rangeMiles string) []SeatGeekEvent {

	if requestExpirationTime.Sub(time.Now()) > 0 {
		//If postalcode, rangeMiles, genres? is something in the cache...
		fmt.Println("NOT expired!! Here's your cached value")
		//Return seatGeekEvents here :)
	}

	requestExpirationTime = time.Now().Add(time.Hour * 24)

	BaseSeatGeekLocalEventsURL := "https://api.seatgeek.com/2/events?client_id=" +
		SEATGEEK_ID +
		"&geoip=" +
		postalCode +
		"&range=" +
		rangeMiles +
		"mi"

	t4 := time.Now()

	totalSeatgeekEvents := FindTotalSeatgeekEvents(BaseSeatGeekLocalEventsURL)

	if totalSeatgeekEvents < 100 {
		var seatGeekEvents []SeatGeekEvent
		seatGeekChan := make(chan []SeatGeekEvent)
		go MakeSeatgeekEventsRequest(BaseSeatGeekLocalEventsURL, 1, seatGeekChan)
		//****TO-DO select only one array here and return
		return seatGeekEvents
	}

	totalPages := int(math.Ceil(float64(totalSeatgeekEvents) / 100))

	var seatGeekEventChannels []chan []SeatGeekEvent

	for pageNumber := 1; pageNumber <= totalPages; pageNumber++ {
		seatGeekChan := make(chan []SeatGeekEvent)
		seatGeekEventChannels = append(seatGeekEventChannels, seatGeekChan)
		go MakeSeatgeekEventsRequest(BaseSeatGeekLocalEventsURL, pageNumber, seatGeekChan)
	}

	cases := make([]reflect.SelectCase, len(seatGeekEventChannels))

	for i, seatGeekEventChan := range seatGeekEventChannels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(seatGeekEventChan)}
	}

	var seatGeekEvents []SeatGeekEvent
	remainingCases := len(cases)

	for remainingCases > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			//Channel has been closed; zero out channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remainingCases--
			continue
		}
		seatGeekEvents = append(seatGeekEvents, value.Interface().([]SeatGeekEvent)...)
	}

	fmt.Println("[Time benchmark] Makin slow calls " + time.Since(t4).String())

	return seatGeekEvents
}

//FindTotalSeatgeekEvents returns the total amount of Seatgeek events in the area
func FindTotalSeatgeekEvents(baseURL string) int {
	SingleEventSeatgeekURL := baseURL + "&per_page=1"

	resp, err := http.Get(SingleEventSeatgeekURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 0
	} else {
		fmt.Printf("Obtained local events data from SeatGeek.\n")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 0
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	metaData := responseData["meta"].(interface{})
	return int(metaData.(map[string]interface{})["total"].(float64))
}

//MakeSeatgeekEventsRequest performs an HTTP request to obtain a single page of event information for an area
func MakeSeatgeekEventsRequest(baseURL string, pageNumber int, seatGeekChan chan<- []SeatGeekEvent) {
	SeatGeekLocalMusicEventsURL := baseURL + "&type=concert&type=music_festival&datetime_utc.gte=2019-11-05&datetime_utc.lte=2019-11-12&per_page=100&page=" + strconv.FormatInt(int64(pageNumber), 10)

	resp, err := http.Get(SeatGeekLocalMusicEventsURL)
	if err != nil {
		fmt.Println("[MakeSeatgeekEventsRequest] initial GET")
		fmt.Fprintf(os.Stderr, err.Error())
		seatGeekChan <- nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[MakeSeatgeekEventsRequest] ioutil ReadAll")
		fmt.Fprintf(os.Stderr, err.Error())
		seatGeekChan <- nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("[MakeSeatgeekEventsRequest] unmarshal")
		fmt.Fprintf(os.Stderr, err.Error())
	}

	eventsFromResponse := responseData["events"].([]interface{})
	seatGeekEvents := make([]SeatGeekEvent, len(eventsFromResponse))
	genreChannels := make([][]chan []string, len(eventsFromResponse))

	for i, event := range eventsFromResponse {
		eventData := event.(map[string]interface{})

		seatGeekEvents[i].Title = eventData["title"].(string)
		seatGeekEvents[i].EventType = eventData["type"].(string)
		seatGeekEvents[i].LocalShowtime = eventData["datetime_local"].(string)

		venueData := eventData["venue"].(map[string]interface{})

		seatGeekEvents[i].VenueName = venueData["name"].(string)
		seatGeekEvents[i].VenueLocation = venueData["display_location"].(string)
		seatGeekEvents[i].URL = venueData["url"].(string)

		performersArray := eventData["performers"].([]interface{})

		for _, performer := range performersArray {
			performerData := performer.(map[string]interface{})
			channel := make(chan []string)
			genreChannels[i] = append(genreChannels[i], channel)

			seatGeekEvents[i].Performers = append(seatGeekEvents[i].Performers, performerData["short_name"].(string))

			if performerData["genres"] != nil {
				genreArray := performerData["genres"].([]interface{})

				for _, genre := range genreArray {
					genreData := genre.(map[string]interface{})
					seatGeekEvents[i].Genres = append(seatGeekEvents[i].Genres, genreData["slug"].(string))
				}
			}
		}
	}

	seatGeekChan <- seatGeekEvents
	close(seatGeekChan)
}
