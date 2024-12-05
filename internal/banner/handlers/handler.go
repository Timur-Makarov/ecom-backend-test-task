package handlers

import (
	"ecom-backend-test-task/internal/banner/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type BannerHandler struct {
	Service services.BannerService
}

func (h BannerHandler) AddBanner(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		Name string `json:"name"`
	}

	var reqBody RequestBody

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err := h.Service.AddBanner(reqBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h BannerHandler) UpdateCounterStats(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bannerIdString := req.PathValue("bannerID")
	bannerIdInt, err := strconv.Atoi(bannerIdString)
	if err != nil || bannerIdInt < 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	h.Service.UpdateBannerCounterStats(uint64(bannerIdInt))
}

func (h BannerHandler) GetCounterStats(w http.ResponseWriter, req *http.Request) {
	bannerIdString := req.PathValue("bannerID")
	bannerIdInt, err := strconv.Atoi(bannerIdString)
	if err != nil || bannerIdInt < 0 {
		http.Error(w, "Bad Request: invalid bannerID", http.StatusBadRequest)
		return
	}

	tsFromString := req.URL.Query().Get("tsFrom")
	tsToString := req.URL.Query().Get("tsTo")

	tsFromInt, err := strconv.Atoi(tsFromString)
	if err != nil || tsFromInt < 0 {
		http.Error(w, "Bad request: invalid tsFrom", http.StatusBadRequest)
		return
	}

	tsToInt, err := strconv.Atoi(tsToString)
	if err != nil || tsToInt < 0 {
		http.Error(w, "Bad request: invalid tsFrom", http.StatusBadRequest)
		return
	}

	res, err := h.Service.GetCounterStats(uint64(bannerIdInt), uint64(tsFromInt), uint64(tsToInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
