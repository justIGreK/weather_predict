package handler

import (
	"encoding/json"
	"net/http"
)

type Weather interface {
	GetWeatherSummary(city, date string) (map[string]string, error)
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	date := r.URL.Query().Get("date")

	if city == "" || date == "" {
		http.Error(w, "Missing 'city' or 'date' parameter", http.StatusBadRequest)
		return
	}

	summary, err := h.GetWeatherSummary(city, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
