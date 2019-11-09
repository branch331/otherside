package seatgeekLayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"otherside/api/redisLayer"
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

//TimeToday saves the truncated value of the beginning and end times of the day
type TimeToday struct {
	BeginningOfDay time.Time
	EndOfDay       time.Time
}

var SEATGEEK_ID = os.Getenv("SEATGEEK_ID")
var requestExpirationTime time.Time
var timeToday TimeToday
var currentEventsTEST []SeatGeekEvent

//FindLocalEvents makes a request to the SeatGeek Events API using the postal code and range,
//and returns an array of SeatGeekEvents.
func FindLocalEvents(postalCode string, rangeMiles string) []SeatGeekEvent {

	t4 := time.Now()

	UTCTimeLocation, err := time.LoadLocation("UTC")
	if err != nil {
		fmt.Printf(err.Error())
	}

	redisLayer.Ping()
	redisLayer.SetSeatgeekEvents("78759")
	redisLayer.GetSeatgeekEvents("78759")

	if timeToday.EndOfDay.Sub(time.Now().In(UTCTimeLocation)) > 0 {
		//If postalcode s something in the cache...
		fmt.Println("NOT expired!! Here's your cached value")
		//Return seatGeekEvents here :) based on requested genres
		fmt.Println("[Time benchmark] Makin slow calls " + time.Since(t4).String())
		return currentEventsTEST
	}

	timeToday = GetTimeToday(UTCTimeLocation)

	BaseSeatGeekLocalEventsURL := "https://api.seatgeek.com/2/events?client_id=" +
		SEATGEEK_ID +
		"&geoip=" +
		postalCode +
		"&range=" +
		rangeMiles +
		"mi"

	totalSeatgeekEvents := FindTotalSeatgeekEvents(BaseSeatGeekLocalEventsURL)

	if totalSeatgeekEvents < 100 {
		var seatGeekEvents []SeatGeekEvent
		seatGeekChan := make(chan []SeatGeekEvent)
		go MakeSeatgeekEventsRequest(BaseSeatGeekLocalEventsURL, 1, timeToday, seatGeekChan)
		//****TO-DO select only one array here and return
		return seatGeekEvents
	}

	totalPages := int(math.Ceil(float64(totalSeatgeekEvents) / 100))

	var seatGeekEventChannels []chan []SeatGeekEvent

	for pageNumber := 1; pageNumber <= totalPages; pageNumber++ {
		seatGeekChan := make(chan []SeatGeekEvent)
		seatGeekEventChannels = append(seatGeekEventChannels, seatGeekChan)
		go MakeSeatgeekEventsRequest(BaseSeatGeekLocalEventsURL, pageNumber, timeToday, seatGeekChan)
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

	currentEventsTEST = append(currentEventsTEST, seatGeekEvents...)
	return seatGeekEvents
}

//GetTimeToday returns struct of the upper and lower datetime bounds of the current day in UTC
func GetTimeToday(loc *time.Location) TimeToday {
	var timeToday TimeToday

	currentTime := time.Now().In(loc)

	timeBeginningOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, loc)
	timeEndOfDay := timeBeginningOfDay.Add(24 * time.Hour)

	timeToday.BeginningOfDay = timeBeginningOfDay
	timeToday.EndOfDay = timeEndOfDay

	return timeToday
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
func MakeSeatgeekEventsRequest(baseURL string, pageNumber int, timeToday TimeToday, seatGeekChan chan<- []SeatGeekEvent) {

	SeatGeekLocalMusicEventsURL := baseURL +
		"&datetime_utc.gte=" +
		timeToday.BeginningOfDay.Format("2006-01-02") +
		"&datetime_utc.lte=" +
		timeToday.EndOfDay.Format("2006-01-02") +
		"&type=concert&type=music_festival&per_page=100&page=" +
		strconv.FormatInt(int64(pageNumber), 10)

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
