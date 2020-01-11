package seatgeekLayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"otherside/api/redisLayer"
	"reflect"
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

//FindLocalEvents makes a request to the SeatGeek Events API using the postal code and range,
//and returns an array of SeatGeekEvents.
func FindLocalEvents(postalCodes []string, genres []string, timeToday TimeToday) []SeatGeekEvent {

	t4 := time.Now()

	var seatGeekEventChannels []chan []SeatGeekEvent
	var seatGeekEvents []SeatGeekEvent

	for _, postCode := range postalCodes {

		postCodeAlreadyCached, err := redisLayer.Exists(postCode)
		if err != nil {
			fmt.Printf("Error checking if postcode key %s exists in Redis: "+err.Error(), postCode)
		}

		if postCodeAlreadyCached {
			redisData, err := redisLayer.GetKeyBytes(postCode)

			if err != nil {
				fmt.Printf("Error getting value for postcode key %s in Redis: "+err.Error(), postCode)
			}

			var cachedSeatgeekEvents []SeatGeekEvent
			json.Unmarshal(redisData, &cachedSeatgeekEvents)
			if err != nil {
				fmt.Printf("Error unmarshalling value for postcode key %s from Redis: "+err.Error(), postCode)
			}

			seatGeekEvents = append(seatGeekEvents, cachedSeatgeekEvents...)
		} else {
			BaseSeatGeekLocalEventsURL := "https://api.seatgeek.com/2/events?client_id=" +
				SEATGEEK_ID +
				"&range=" +
				"50mi"

			seatGeekChan := make(chan []SeatGeekEvent)
			seatGeekEventChannels = append(seatGeekEventChannels, seatGeekChan)
			go MakeSeatgeekEventsRequest(BaseSeatGeekLocalEventsURL, postCode, timeToday, seatGeekChan)
		}
	}

	cases := make([]reflect.SelectCase, len(seatGeekEventChannels))

	for i, seatGeekEventChan := range seatGeekEventChannels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(seatGeekEventChan)}
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

		seatGeekEvents = append(seatGeekEvents, value.Interface().([]SeatGeekEvent)...)
	}

	fmt.Println("[Time benchmark] Makin slow calls " + time.Since(t4).String())
	return FilterByGenres(seatGeekEvents, genres)
}

//FilterByGenres returns an array of SeatGeekEvent items containing only the genres listed
func FilterByGenres(events []SeatGeekEvent, genres []string) []SeatGeekEvent {
	var filteredEvents []SeatGeekEvent

	for _, event := range events {
		for _, genre := range genres {
			if stringInSlice(genre, event.Genres) {
				filteredEvents = append(filteredEvents, event)
				break
			}
		}
	}

	return filteredEvents
}

func stringInSlice(stringToFind string, list []string) bool {
	for _, val := range list {
		if stringToFind == val {
			return true
		}
	}
	return false
}

//MakeSeatgeekEventsRequest performs an HTTP request to obtain a single page of event information for an area
func MakeSeatgeekEventsRequest(baseURL string, postCode string, timeToday TimeToday, seatGeekChan chan<- []SeatGeekEvent) {

	SeatGeekLocalMusicEventsURL := baseURL +
		"&geoip=" +
		postCode +
		"&datetime_utc.gte=" +
		timeToday.BeginningOfDay.Format("2006-01-02") +
		"&datetime_utc.lte=" +
		timeToday.EndOfDay.AddDate(0, 0, 7).Format("2006-01-02") +
		"&type=concert&type=music_festival&per_page=100&page=1"

	resp, err := http.Get(SeatGeekLocalMusicEventsURL)
	if err != nil {
		fmt.Printf("Error making request to SeatGeek API %s: "+err.Error(), SeatGeekLocalMusicEventsURL)
		seatGeekChan <- nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error with ReadAll for local events response: " + err.Error())
		seatGeekChan <- nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("Error with SeatGeek response unmarshalling: " + err.Error())
	}

	if responseData["events"] == nil {
		seatGeekChan <- nil
		redisLayer.SetKeyBytes(postCode, nil)
		close(seatGeekChan)
	} else {
		eventsFromResponse := responseData["events"].([]interface{})
		seatGeekEvents := make([]SeatGeekEvent, len(eventsFromResponse))

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

		seatGeekEventsSerialized, err := json.Marshal(seatGeekEvents)
		if err != nil {
			fmt.Printf("Error marshalling seatgeek events data: " + err.Error())
		}

		redisLayer.SetKeyBytes(postCode, seatGeekEventsSerialized)
		
		seatGeekChan <- seatGeekEvents
		close(seatGeekChan)
	}
}

func GetTimeToday(loc *time.Location) TimeToday {
	var timeToday TimeToday

	currentTime := time.Now().In(loc)

	timeBeginningOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, loc)
	timeEndOfDay := timeBeginningOfDay.Add(24 * time.Hour)

	timeToday.BeginningOfDay = timeBeginningOfDay
	timeToday.EndOfDay = timeEndOfDay

	return timeToday
}
