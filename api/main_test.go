package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestPostCodeMap(t *testing.T) {
	postcodeMap := generateCityPostcodeMap()
	if postcodeMap["Austin TX"] != "78759" {
		t.Errorf("Expected 78759; got %s", postcodeMap["Austin TX"])
	}
	if postcodeMap["Troy NY"] != "12180" {
		t.Errorf("Expected 78759; got %s", postcodeMap["Austin TX"])
	}
}

func TestLocalEvents(t *testing.T) {
	req, err := http.NewRequest("GET", "/localevents?cities=[Austin TX]&genres=[rock]", nil)
	if err != nil {
		fmt.Printf("Error obtaining localevents: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LocalEvents)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("Unexpected status code: %v", status)
	}
}

func TestLocalEventsMissingCities(t *testing.T) {
	req, err := http.NewRequest("GET", "/localevents", nil)
	if err != nil {
		fmt.Printf("Error obtaining localevents: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LocalEvents)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusInternalServerError {
		t.Errorf("Unexpected status code: %v", status)
	}

	body := rr.Body.String()
	if body != "Cities array parameter missing from request." {
		t.Errorf("Unexpected response: %s", body)
	}
}

func TestLocalEventsMissingGenres(t *testing.T) {
	req, err := http.NewRequest("GET", "/localevents?cities=[Austin TX,Boston MA]", nil)
	if err != nil {
		fmt.Printf("Error obtaining localevents: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LocalEvents)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusInternalServerError {
		t.Errorf("Unexpected status code: %v", status)
	}

	body := rr.Body.String()
	if body != "Genres array parameter missing from request." {
		t.Errorf("Unexpected response: %s", body)
	}
}

func TestCallbackMissingParam(t *testing.T) {
	req, err := http.NewRequest("GET", "/callback", nil)
	if err != nil {
		fmt.Printf("Error obtaining callback: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Callback)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusInternalServerError {
		t.Errorf("Unexpected status code: %v", status)
	}

	body := rr.Body.String()
	if body != "State parameter missing from request." {
		t.Errorf("Unexpected response: %s", body)
	}
}

func TestBuildPlaylistMissingName(t *testing.T) {
	req, err := http.NewRequest("GET", "/buildplaylist", nil)
	if err != nil {
		fmt.Printf("Error obtaining buildplaylist: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BuildPlaylist)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusInternalServerError {
		t.Errorf("Unexpected status code: %v", status)
	}

	body := rr.Body.String()
	if body != "name parameter missing from request." {
		t.Errorf("Unexpected response: %s", body)
	}
}

func TestBuildPlaylistMissingDesc(t *testing.T) {
	req, err := http.NewRequest("GET", "/buildplaylist?name=wawa", nil)
	if err != nil {
		fmt.Printf("Error obtaining buildplaylist: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BuildPlaylist)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusInternalServerError {
		t.Errorf("Unexpected status code: %v", status)
	}

	body := rr.Body.String()
	if body != "desc parameter missing from request." {
		t.Errorf("Unexpected response: %s", body)
	}
}

func TestQueryStringToArray(t *testing.T) {
	array := QueryStringToArray("[Blarg,Glarb,Bebarg]")

	if len(array) != 3 {
		t.Errorf("Expected array of size 3; got %v", len(array))
	}

	if array[0] != "Blarg" {
		t.Errorf("Unexpected array value: %s", array[0])
	}
}
