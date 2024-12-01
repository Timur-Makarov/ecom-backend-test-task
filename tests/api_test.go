package main

import (
	"bytes"
	"ecom-backend-test-task/config"
	"ecom-backend-test-task/internal/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func setup(t *testing.T) *config.APIComponents {
	hardCodedDSN := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := config.OpenDB(hardCodedDSN)

	if err != nil {
		t.Fatal(err)
	}

	apiComponents := config.SetupAPIComponents(db)
	return apiComponents
}

func TestAddBanner(t *testing.T) {
	apiComponents := setup(t)

	type Request struct {
		Name string `json:"name"`
	}

	requestData := Request{"New-Banner"}
	jsonBody, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/banner/add", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiComponents.Handlers.AddBanner)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetBannerCounterStats(t *testing.T) {
	apiComponents := setup(t)

	hardCodedBannerID := "1"

	requestURL := url.URL{
		Path: "/stats/" + hardCodedBannerID,
	}

	ts := time.Now()
	timestampFrom := ts.Truncate(time.Minute).Unix()
	timestampTo := ts.Truncate(time.Minute).Add(time.Minute).Unix() - 1

	tsFrom := strconv.FormatInt(timestampFrom, 10)
	tsTo := strconv.FormatInt(timestampTo, 10)

	queryParams := requestURL.Query()
	queryParams.Add("tsFrom", tsFrom)
	queryParams.Add("tsTo", tsTo)

	requestURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetPathValue("bannerID", hardCodedBannerID)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiComponents.Handlers.GetBannerCounterStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error [%v]: %v", status, rr.Body.String())
	}

	var response services.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

func TestUpdateBannerCounterStats(t *testing.T) {
	apiComponents := setup(t)

	hardCodedBannerID := "1"

	req, err := http.NewRequest("POST", "/counter/"+hardCodedBannerID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("bannerID", hardCodedBannerID)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiComponents.Handlers.UpdateBannerCounterStats)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error [%v]: %v", status, rr.Body.String())
	}
}
