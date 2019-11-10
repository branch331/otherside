package seatgeek_test

import (
	"otherside/api/seatgeekLayer"
	"testing"
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
