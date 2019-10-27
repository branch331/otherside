package seatgeekLayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type SeatGeekEvent struct {
	title         string
	eventType     string
	url           string
	performers    []string
	genres        []string //possibly nil
	localShowtime string
	venueName     string
	venueLocation string
}

var SEATGEEK_ID = os.Getenv("SEATGEEK_ID")

func findLocalEvents(postalCode string, rangeMiles string) []SeatGeekEvent {
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

		seatGeekEvents[i].title = eventData["title"].(string)
		seatGeekEvents[i].eventType = eventData["type"].(string)
		seatGeekEvents[i].localShowtime = eventData["datetime_local"].(string)

		venueData := eventData["venue"].(map[string]interface{})

		seatGeekEvents[i].venueName = venueData["name"].(string)
		seatGeekEvents[i].venueLocation = venueData["display_location"].(string)
		seatGeekEvents[i].url = venueData["url"].(string)

		performersArray := eventData["performers"].([]interface{})

		for _, performer := range performersArray {
			performerData := performer.(map[string]interface{})
			seatGeekEvents[i].performers = append(seatGeekEvents[i].performers, performerData["short_name"].(string))
			seatGeekEvents[i].genres = GetSeatGeekArtistGenres(fmt.Sprintf("%d", int(performerData["id"].(float64))))
		}
	}

	return seatGeekEvents
}

func GetSeatGeekArtistGenres(performerId string) []string {
	SeatGeekPerformerURL := "https://api.seatgeek.com/2/performers/" + performerId + "?client_id=" + SEATGEEK_ID
	fmt.Printf(SeatGeekPerformerURL)

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

	if genresFromResponse, keyExists := responseData["genres"].([]interface{}); keyExists {
		for _, genre := range genresFromResponse {
			genreData := genre.(map[string]interface{})
			genreArray = append(genreArray, genreData["slug"].(string))
		}
	}

	return genreArray
}
