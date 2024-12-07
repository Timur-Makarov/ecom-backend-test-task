package main

import (
	"bytes"
	"ecom-backend-test-task/internal/pkg/app"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

var a *app.App

func init() {
	var err error
	a, err = app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateBanner(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}

	requestData := Request{"New-Banner"}
	jsonBody, err := json.Marshal(requestData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/banners", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err := a.Http.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetCounterStatistics(t *testing.T) {
	hardCodedBannerID := "1"

	requestURL := url.URL{
		Path: "/banners/" + hardCodedBannerID + "/stats/",
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
	assert.NoError(t, err)

	res, err := a.Http.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	type GetCounterStatsDTO struct {
		BannerID      uint64 `json:"bannerId"`
		Count         uint64 `json:"count"`
		TimestampFrom uint64 `json:"timestampFrom"`
		TimestampTo   uint64 `json:"timestampTo"`
	}

	var response GetCounterStatsDTO

	assert.NoError(t, json.NewDecoder(res.Body).Decode(&response))
}

func TestUpdateBannerCounterStatistics(t *testing.T) {
	hardCodedBannerID := "1"

	req, err := http.NewRequest("PUT", "/banners/"+hardCodedBannerID+"/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := a.Http.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
