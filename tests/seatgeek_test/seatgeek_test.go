package seatgeek_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"otherside/api/seatgeekLayer"
	"testing"
	"time"
)

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

func TestSeatgeekEventsRequest(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"events": [{"title": "Winter Wonderland","type": "concert","datetime_local": "now","venue": {"name": "The Stage","display_location": "The Stage","url": "myURL"}, 		"performers": [{"short_name": "Billy Bob","genres": [{"slug": "rock"}]}]}]}`)
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

	if events[0].Title != "Winter Wonderland" {
		t.Errorf("Unexpected title: %s", events[0].Title)
	} else if events[0].EventType != "concert" {
		t.Errorf("Unexpected type: %s", events[0].EventType)
	} else if events[0].URL != "myURL" {
		t.Errorf("Unexpected URL: %s", events[0].URL)
	} else if events[0].Performers[0] != "Billy Bob" {
		t.Errorf("Unexpected Performer: %s", events[0].Performers[0])
	} else if events[0].Genres[0] != "rock" {
		t.Errorf("Unexpected genre: %s", events[0].Genres[0])
	} else if events[0].LocalShowtime != "now" {
		t.Errorf("Unexpected showtime: %s", events[0].LocalShowtime)
	} else if events[0].VenueName != "The Stage" {
		t.Errorf("Unexpected venue name: %s", events[0].VenueName)
	} else if events[0].VenueLocation != "The Stage" {
		t.Errorf("Unexpected venue location: %s", events[0].VenueLocation)
	}
}
