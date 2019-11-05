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

//FindLocalEvents makes a request to the SeatGeek Events API using the postal code and range,
//and returns an array of SeatGeekEvents.
func FindLocalEvents(postalCode string, rangeMiles string) []SeatGeekEvent {
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

	//var seatGeekEvents2D = make([][]SeatGeekEvent, totalPages)
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
	SeatGeekLocalEventsURL := baseURL + "&per_page=100&page=" + strconv.FormatInt(int64(pageNumber), 10)

	resp, err := http.Get(SeatGeekLocalEventsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		seatGeekChan <- nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		seatGeekChan <- nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
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
			go GetSeatGeekArtistGenres(fmt.Sprintf("%d", int(performerData["id"].(float64))), channel)
		}
	}

	cases := make([][]reflect.SelectCase, len(genreChannels))

	remainingCases := 0

	for k, genreChannelArray := range genreChannels {
		cases[k] = make([]reflect.SelectCase, len(genreChannelArray))

		for m, genreChannel := range genreChannelArray {
			cases[k][m] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(genreChannel)}
			remainingCases++
		}
	}

	for remainingCases > 0 {
		for n, caseArray := range cases {
			for _, currentCase := range caseArray {
				_, value, ok := reflect.Select(caseArray)

				if !ok {
					//Channel has been closed; zero out channel to disable the case
					currentCase.Chan = reflect.ValueOf(nil)
					remainingCases--
					continue
				}
				seatGeekEvents[n].Genres = append(seatGeekEvents[n].Genres, value.Interface().([]string)...)
			}
		}
	}

	seatGeekChan <- seatGeekEvents
	close(seatGeekChan)
}

//GetSeatGeekArtistGenres returns an array of all genres pertinent to a performer.
func GetSeatGeekArtistGenres(performerID string, genreChannel chan<- []string) {
	SeatGeekPerformerURL := "https://api.seatgeek.com/2/performers/" + performerID + "?client_id=" + SEATGEEK_ID

	resp, err := http.Get(SeatGeekPerformerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		genreChannel <- nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		genreChannel <- nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		genreChannel <- nil
	}

	var genreArray []string

	if status, statusExists := responseData["status"].(string); statusExists {
		genreChannel <- append(genreArray, status)
	}

	if genresFromResponse, keyExists := responseData["genres"].([]interface{}); keyExists {
		for _, genre := range genresFromResponse {
			genreData := genre.(map[string]interface{})
			genreArray = append(genreArray, genreData["slug"].(string))
		}
	}

	genreChannel <- genreArray
	close(genreChannel)
}
