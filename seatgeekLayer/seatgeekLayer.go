package seatgeekLayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type SeatGeekEvent struct {
	Title         string
	EventType     string
	Url           string
	Performers    []string
	Genres        []string //possibly nil
	LocalShowtime string
	VenueName     string
	VenueLocation string
}

var SEATGEEK_ID = os.Getenv("SEATGEEK_ID")

func FindLocalEvents(postalCode string, rangeMiles string) []SeatGeekEvent {
	SeatGeekLocalEventsURL := "https://api.seatgeek.com/2/events?client_id=" +
		SEATGEEK_ID +
		"&geoip=" +
		postalCode +
		"&range=" +
		rangeMiles +
		"mi"

	resp, err := http.Get(SeatGeekLocalEventsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	} else {
		fmt.Printf("Obtained local events data from SeatGeek.\n")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

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
		seatGeekEvents[i].Url = venueData["url"].(string)

		performersArray := eventData["performers"].([]interface{})

		for _, performer := range performersArray {
			performerData := performer.(map[string]interface{})
			seatGeekEvents[i].Performers = append(seatGeekEvents[i].Performers, performerData["short_name"].(string))
			seatGeekEvents[i].Genres = GetSeatGeekArtistGenres(fmt.Sprintf("%d", int(performerData["id"].(float64))))
		}
	}

	return seatGeekEvents
}

func GetSeatGeekArtistGenres(performerId string) []string {
	SeatGeekPerformerURL := "https://api.seatgeek.com/2/performers/" + performerId + "?client_id=" + SEATGEEK_ID

	resp, err := http.Get(SeatGeekPerformerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	var genreArray []string

	if status, statusExists := responseData["status"].(string); statusExists {
		return append(genreArray, status)
	}

	if genresFromResponse, keyExists := responseData["genres"].([]interface{}); keyExists {
		for _, genre := range genresFromResponse {
			genreData := genre.(map[string]interface{})
			genreArray = append(genreArray, genreData["slug"].(string))
		}
	}

	return genreArray
}
