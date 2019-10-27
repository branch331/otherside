package seatgeek_test

import (
	"otherside/api/seatgeekLayer"
	"testing"
)

func TestGenreLookupBeastieBoys(t *testing.T) {
	genres := seatgeekLayer.GetSeatGeekArtistGenres("266") //Beastie Boys

	if len(genres) != 4 {
		t.Errorf("Expected 4 genres for Beastie Boys. Instead, got %d.", len(genres))
	}

	if genres[2] != "hip-hop" {
		t.Errorf("Expected hip-hop genre. Instead, got %s.", genres[2])
	}
}

func TestGenreLookupNonExistingArtist(t *testing.T) {
	genres := seatgeekLayer.GetSeatGeekArtistGenres("44444") //Non-existing performer ID

	if genres[0] != "error" {
		t.Errorf("Expected error; did not receive error response.")
	}
}
