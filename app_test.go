package main

import (
	"strings"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthCheckEndpoint(t *testing.T) {
    // Create a request to pass to our handler. No query parameters, so we'll
    // pass 'nil' as the third parameter.
    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatal(err)
    }

    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HealthCheckEndpoint)

    // handlers satisfy http.Handler, so call ServeHTTP method
    // directly and pass in Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check the response body is what we expect.
    expected := `{"alive":true}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestCreateTodoEndPoint(t *testing.T) {
	bodyReader := strings.NewReader(`{"description":"buy beer", "completed":false}`)
	req, err := http.NewRequest("POST", "/todos", bodyReader)
    if err != nil {
        t.Fatal(err)
    }

    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(CreateTodoEndPoint)

    // handlers satisfy http.Handler, so call ServeHTTP method
    // directly and pass in Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
    }
}