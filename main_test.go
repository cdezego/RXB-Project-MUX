package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// The primary TestGetFilms Test method will test the full GetFilms response and compare vs a saved test response
func TestGetFilms(t *testing.T) {
	// Switch to test mode
	gin.SetMode(gin.TestMode)

	// Setup the router
	r := gin.Default()
	r.GET("/", getFilms)

	// Create the mock request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	// Check to see if the response was what you expected
	if w.Code == http.StatusOK {
		content, err := ioutil.ReadFile("test_response.json")
		if err == nil {
			if string(content) == w.Body.String() {
				t.Logf("Test Case PASSED JSON VALIDATION")
			} else {
				t.Fatalf("Unable to validate the JSON response!")
			}
		}

	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}

func TestGetByFilm_ID(t *testing.T) {

}

func TestGetByTitle(t *testing.T) {

}

func TestGetByRating(t *testing.T) {

}

func TestGetByCategory(t *testing.T) {

}
