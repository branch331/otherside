package seatgeek_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"otherside/api/seatgeekLayer"
	"testing"
	"time"
)

type seatGeekTestData struct {
	Title              string
	Type               string
	Datetime           string
	VenueName          string
	VenueLocation      string
	VenueURL           string
	PerformerName      string
	PerformerGenreSlug string
}

var seatGeekData seatGeekTestData

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	seatGeekData.Title = "Winter Wonderland"
	seatGeekData.Type = "Concert"
	seatGeekData.Datetime = "Now"
	seatGeekData.VenueName = "The Stage"
	seatGeekData.VenueLocation = "Also The Stage"
	seatGeekData.VenueURL = "URL"
	seatGeekData.PerformerName = "Billy Ray Dyrus"
	seatGeekData.PerformerGenreSlug = "Rock n' roll"
}

func TestFilterByGenres(t *testing.T) {
	events := []seatgeekLayer.SeatGeekEvent{
		seatgeekLayer.SeatGeekEvent{
			Title:  "A",
			Genres: []string{"rock", "hip-hop"},
		},
		seatgeekLayer.SeatGeekEvent{
			Title:  "B",
			Genres: []string{"electronic", "hip-hop"},
		},
		seatgeekLayer.SeatGeekEvent{
			Title:  "C",
			Genres: []string{"country", "rock"},
		},
		seatgeekLayer.SeatGeekEvent{
			Title:  "D",
			Genres: []string{"indie"},
		},
	}

	genreFilter := []string{"rock"}
	filteredEvents := seatgeekLayer.FilterByGenres(events, genreFilter)

	if filteredEvents[0].Title != "A" || filteredEvents[1].Title != "C" {
		t.Errorf("Expected A and C for event titles. Got: %s and %s", filteredEvents[0].Title, filteredEvents[1].Title)
	}

	genreFilter = []string{"electronic", "indie"}
	filteredEvents = seatgeekLayer.FilterByGenres(events, genreFilter)

	if filteredEvents[0].Title != "B" || filteredEvents[1].Title != "D" {
		t.Errorf("Expected B and D for event titles. Got: %s and %s", filteredEvents[0].Title, filteredEvents[1].Title)
	}
}

func TestSeatgeekSingleEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, fmt.Sprintf(
			`{"events": [
				{
					"title": "%s",
					"type": "%s",
					"datetime_local": "%s",
					"venue": {"name": "%s",
					"display_location": "%s",
					"url": "%s"},
					"performers": 
					[
						{
							"short_name": "%s",
							"genres": [{"slug": "%s"}]
						}	
					]
				}		
			]}`,
			seatGeekData.Title,
			seatGeekData.Type,
			seatGeekData.Datetime,
			seatGeekData.VenueName,
			seatGeekData.VenueLocation,
			seatGeekData.VenueURL,
			seatGeekData.PerformerName,
			seatGeekData.PerformerGenreSlug))
	}))
	defer ts.Close()

	seatgeekURL := ts.URL
	eventsChan := make(chan []seatgeekLayer.SeatGeekEvent)

	var timeToday seatgeekLayer.TimeToday
	UTCTimeLocation, err := time.LoadLocation("UTC")
	if err != nil {
		fmt.Printf("Error creating time LoadLocation: " + err.Error())
	}

	timeToday = seatgeekLayer.GetTimeToday(UTCTimeLocation)

	go seatgeekLayer.MakeSeatgeekEventsRequest(seatgeekURL+"/test", "", timeToday, eventsChan)

	events := <-eventsChan

	if events[0].Title != seatGeekData.Title {
		t.Errorf("Unexpected title: %s", events[0].Title)
	} else if events[0].EventType != seatGeekData.Type {
		t.Errorf("Unexpected type: %s", events[0].EventType)
	} else if events[0].URL != seatGeekData.VenueURL {
		t.Errorf("Unexpected URL: %s", events[0].URL)
	} else if events[0].Performers[0] != seatGeekData.PerformerName {
		t.Errorf("Unexpected Performer: %s", events[0].Performers[0])
	} else if events[0].Genres[0] != seatGeekData.PerformerGenreSlug {
		t.Errorf("Unexpected genre: %s", events[0].Genres[0])
	} else if events[0].LocalShowtime != seatGeekData.Datetime {
		t.Errorf("Unexpected showtime: %s", events[0].LocalShowtime)
	} else if events[0].VenueName != seatGeekData.VenueName {
		t.Errorf("Unexpected venue name: %s", events[0].VenueName)
	} else if events[0].VenueLocation != seatGeekData.VenueLocation {
		t.Errorf("Unexpected venue location: %s", events[0].VenueLocation)
	}
}

func TestSeatgeekMultiEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, fmt.Sprintf(
			`{"events": [
				{
					"title": "%s",
					"type": "%s",
					"datetime_local": "%s",
					"venue": {"name": "%s",
					"display_location": "%s",
					"url": "%s"},
					"performers": 
					[
						{
							"short_name": "%s",
							"genres": [{"slug": "%s"}]
						}	
					]
				},
				{
					"title": "%s",
					"type": "%s",
					"datetime_local": "%s",
					"venue": {"name": "%s",
					"display_location": "%s",
					"url": "%s"},
					"performers": 
					[
						{
							"short_name": "%s",
							"genres": [{"slug": "%s"}]
						}	
					]
				}	
			]}`,
			seatGeekData.Title,
			seatGeekData.Type,
			seatGeekData.Datetime,
			seatGeekData.VenueName,
			seatGeekData.VenueLocation,
			seatGeekData.VenueURL,
			seatGeekData.PerformerName,
			seatGeekData.PerformerGenreSlug,
			seatGeekData.Title,
			seatGeekData.Type,
			seatGeekData.Datetime,
			seatGeekData.VenueName,
			seatGeekData.VenueLocation,
			seatGeekData.VenueURL,
			seatGeekData.PerformerName,
			seatGeekData.PerformerGenreSlug))
	}))
	defer ts.Close()

	seatgeekURL := ts.URL
	eventsChan := make(chan []seatgeekLayer.SeatGeekEvent)

	var timeToday seatgeekLayer.TimeToday
	UTCTimeLocation, err := time.LoadLocation("UTC")
	if err != nil {
		fmt.Printf("Error creating time LoadLocation: " + err.Error())
	}

	timeToday = seatgeekLayer.GetTimeToday(UTCTimeLocation)

	go seatgeekLayer.MakeSeatgeekEventsRequest(seatgeekURL+"/test", "", timeToday, eventsChan)

	events := <-eventsChan

	for _, event := range events {
		if event.Title != seatGeekData.Title {
			t.Errorf("Unexpected title: %s", events[0].Title)
		} else if event.EventType != seatGeekData.Type {
			t.Errorf("Unexpected type: %s", events[0].EventType)
		} else if event.URL != seatGeekData.VenueURL {
			t.Errorf("Unexpected URL: %s", events[0].URL)
		} else if event.Performers[0] != seatGeekData.PerformerName {
			t.Errorf("Unexpected Performer: %s", events[0].Performers[0])
		} else if event.Genres[0] != seatGeekData.PerformerGenreSlug {
			t.Errorf("Unexpected genre: %s", events[0].Genres[0])
		} else if event.LocalShowtime != seatGeekData.Datetime {
			t.Errorf("Unexpected showtime: %s", events[0].LocalShowtime)
		} else if event.VenueName != seatGeekData.VenueName {
			t.Errorf("Unexpected venue name: %s", events[0].VenueName)
		} else if event.VenueLocation != seatGeekData.VenueLocation {
			t.Errorf("Unexpected venue location: %s", events[0].VenueLocation)
		}
	}
}
