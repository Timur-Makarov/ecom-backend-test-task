package main

import (
	"bytes"
	"ecom-backend-test-task/internal/pkg/app"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type App struct {
	repositories *app.Repositories
	handlers     *app.Handlers
	services     *app.Services
}

var a = App{}

func init() {
	hardCodedDSN := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	db, err := gorm.Open(postgres.Open(hardCodedDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to init db connection: %w", err))
	}

	a.repositories = app.GetRepositories(db)
	a.services = app.GetServices(a.repositories)
	a.handlers = app.GetHandlers(a.services)
}

func TestAddBanner(t *testing.T) {
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
	handler := http.HandlerFunc(a.handlers.BannerHandler.AddBanner)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetBannerCounterStats(t *testing.T) {
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
	handler := http.HandlerFunc(a.handlers.BannerHandler.GetCounterStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error [%v]: %v", status, rr.Body.String())
	}

	type GetCounterStatsDTO struct {
		BannerID      uint64 `json:"bannerId"`
		Count         uint64 `json:"count"`
		TimestampFrom uint64 `json:"timestampFrom"`
		TimestampTo   uint64 `json:"timestampTo"`
	}

	var response GetCounterStatsDTO

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

func TestUpdateBannerCounterStats(t *testing.T) {
	hardCodedBannerID := "1"

	req, err := http.NewRequest("POST", "/counter/"+hardCodedBannerID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("bannerID", hardCodedBannerID)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(a.handlers.BannerHandler.UpdateCounterStats)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error [%v]: %v", status, rr.Body.String())
	}
}
