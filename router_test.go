package gohttproute

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const testError = "expected to be %v, actual %v"

func TestSetupRoutes(t *testing.T) {
	// Test init function sets up all the routes
	expected := 1
	if len(routes) != expected {
		t.Errorf(testError, expected, len(routes))
	}
}

func TestAddRoute(t *testing.T) {
	expected := 1

	testRoutes := Routes{}

	testRoutes.AddRoute("test", "GET", "/test", nil)
	if len(testRoutes) != expected {
		t.Errorf(testError, expected, len(routes))
	}
}

func TestGetRoute(t *testing.T) {
	r := httptest.NewRequest("GET", "/view", nil)
	route, err := routes.GetRoute(r)
	if err != nil {
		t.Errorf(testError, nil, err)
	}

	if route.name != "view" {
		t.Errorf(testError, "view", route.name)
	}
}

func TestGetInvalidRoute(t *testing.T) {
	r := httptest.NewRequest("GET", "/blah", nil)
	_, err := routes.GetRoute(r)
	if err == nil {
		t.Errorf(testError, "", err)
	}
}

func TestServeNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/blah", nil)
	Serve(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf(testError, http.StatusNotFound, w.Code)
	}
}

func TestServeValidRoute(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/view", nil)
	Serve(w, r)
	if w.Code != http.StatusOK {
		t.Errorf(testError, http.StatusOK, w.Code)
	}
}

func TestServeMethodNotAllowed(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PWN", "/view", nil)
	Serve(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf(testError, http.StatusMethodNotAllowed, w.Code)
	}
}
